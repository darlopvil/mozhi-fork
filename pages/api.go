package pages

import (
	"bufio"
	"encoding/base64"
	"io"
	"strings"

	"codeberg.org/aryak/libmozhi"
	"codeberg.org/aryak/mozhi/utils"
	"github.com/gofiber/fiber/v2"
)

// HandleSourceLanguages godoc
//
//	@Summary	Show list of available source languages for engine
//	@Param		engine	query		string	true	"Engine name"
//	@Success	200		{object}	libmozhi.List
//	@Router		/api/source_languages [get]
func HandleSourceLanguages(c *fiber.Ctx) error {
	engine := utils.Sanitize(c.Query("engine"), "alpha")
	if engine == "" {
		return fiber.NewError(fiber.StatusBadRequest, "engine is a required query string.")
	}
	data, err := libmozhi.LangList(engine, "sl")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(data)
}

// HandleTargetLanguages godoc
//
//	@Summary	Show list of available target languages for engine
//	@Param		engine	query		string	true	"Engine name"
//	@Success	200		{object}	libmozhi.List
//	@Router		/api/target_languages [get]
func HandleTargetLanguages(c *fiber.Ctx) error {
	engine := utils.Sanitize(c.Query("engine"), "alpha")
	if engine == "" {
		return fiber.NewError(fiber.StatusBadRequest, "engine is a required query string.")
	}
	data, err := libmozhi.LangList(engine, "tl")
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(data)
}

// HandleTTS godoc
//
//	@Summary	Get Text-To-Speech for specified language using specified engine
//	@Param		engine	query	string	true	"Engine name"
//	@Param		lang	query	string	true	"Language being TTS'd"
//	@Param		text	query	string	true	"Text being TTS'd"
//	@Router		/api/tts [get]
func HandleTTS(c *fiber.Ctx) error {
	engine := utils.Sanitize(c.Query("engine"), "alpha")
	lang := c.Query("lang")
	text := c.Query("text")
	if engine == "" || text == "" || lang == "" {
		return fiber.NewError(fiber.StatusBadRequest, "engine, lang, text are required query strings.")
	}
	if !utils.EngineEnabled(engine) {
		return fiber.NewError(fiber.StatusBadRequest, "engine "+engine+" is not enabled")
	}
	if utils.TextTooLong(text) {
		return utils.TextTooLongError()
	}
	// Long texts are synthesized in chunks (each engine caps the text per
	// request); every chunk is sent as soon as it is ready, so playback starts
	// after the first one. Audio types are never compressed, so it streams.
	contentType, write, err := libmozhi.TTSStream(engine, lang, text)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	c.Set("Content-Type", contentType)
	c.Context().SetBodyStreamWriter(func(w *bufio.Writer) {
		_ = write(w)
	})
	return nil
}

// HandleTranslate godoc
//
//	@Summary		Translate text
//	@Description	When engine is set to all, it will return an array of libmozhi.LangOut.
//	@Param			engine	query		string	true	"Engine name"
//	@Param			from	query		string	true	"Source language"
//	@Param			to		query		string	true	"Target language"
//	@Param			text	query		string	true	"Text being translated"
//	@Success		200		{object}	libmozhi.LangOut
//	@Router			/api/translate [get]
func HandleTranslate(c *fiber.Ctx) error {
	engine := utils.Sanitize(utils.GetQueryOrFormValue(c, "engine"), "alpha")
	from := utils.GetQueryOrFormValue(c, "from")
	to := utils.GetQueryOrFormValue(c, "to")
	text := utils.GetQueryOrFormValue(c, "text")
	if engine == "" || from == "" || to == "" || text == "" {
		return fiber.NewError(fiber.StatusBadRequest, "from, to, engine, text are required query strings.")
	}
	if utils.TextTooLong(text) {
		return utils.TextTooLongError()
	}
	if engine == "all" {
		// Only the engines enabled on this instance.
		dataarr, _ := libmozhi.TranslateSome(utils.EnabledEngines(), to, from, text)
		return c.JSON(dataarr)
	} else if strings.Contains(engine, ",") == true {
		engineArr := strings.Split(engine, ",")
		for _, e := range engineArr {
			if !utils.EngineEnabled(e) {
				return fiber.NewError(fiber.StatusBadRequest, "engine "+e+" is not enabled")
			}
		}
		dataarr, err := libmozhi.TranslateSome(engineArr, to, from, text)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(dataarr)
	} else {
		if !utils.EngineEnabled(engine) {
			return fiber.NewError(fiber.StatusBadRequest, "engine "+engine+" is not enabled")
		}
		data, err := libmozhi.Translate(engine, to, from, text)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(data)
	}
}

// HandleImg godoc
//
//	@Summary		Translate an image
//	@Description	When engine is set to all, it will return an array of libmozhi.LangOut.
//	@Param			engine	formData		string	true	"Engine name"
//	@Param			from	formData		string	true	"Source language"
//	@Param			to		formData		string	true	"Target language"
//	@Param			image	formData		file	true	"PNG image"
//	@Success		200		{object}	libmozhi.ImgOut
//	@Router			/api/image [post]
func HandleImg(c *fiber.Ctx) error {
	engine := utils.Sanitize(c.FormValue("engine"), "alpha")
	from := c.FormValue("from")
	to := c.FormValue("to")
	file, err := c.FormFile("image")
	if err != nil {
		return err
	}
	if file.Header.Get("Content-Type") != "image/png" {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid image, only png images are accepted")
	}
	fileData, err := file.Open()
	if err != nil {
		return err
	}
	image, err := io.ReadAll(fileData)
	if err != nil {
		return err
	}

	if engine == "" || from == "" || to == "" {
		return fiber.NewError(fiber.StatusBadRequest, "from, to, engine are required arguments.")
	}
	imageEncoded := base64.StdEncoding.EncodeToString(image)
	data, err := libmozhi.ImageGoogle(to, from, imageEncoded)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	return c.JSON(data)
}

// HandleEngines godoc
//
//	@Summary		List engines
//	@Description	Lists available Engines.
//	@Success		200		{object}	map[string]string
//	@Router			/api/engines [get]
func HandleEngines(c *fiber.Ctx) error {
	engines := utils.EngineList()
	serializedData := make(map[string]interface{}, len(engines))
	for engineId, engineName := range engines {
		serializedData[engineId] = engineName
	}
	return c.JSON(serializedData)
}
