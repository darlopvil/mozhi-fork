package pages

import (
	"codeberg.org/aryak/libmozhi"
	"strings"
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
	lang := utils.Sanitize(c.Query("lang"), "alpha")
	text := c.Query("text")
	if engine == "" || text == "" || lang == "" {
		return fiber.NewError(fiber.StatusBadRequest, "engine, lang, text are required query strings.")
	}
	data, err := libmozhi.TTS(engine, lang, text)
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, err.Error())
	}
	c.Set("Content-Type", "audio/mpeg")
	return c.Status(fiber.StatusOK).Send(data)
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
	from := utils.Sanitize(utils.GetQueryOrFormValue(c, "from"), "alpha")
	to := utils.Sanitize(utils.GetQueryOrFormValue(c, "to"), "alpha")
	text := utils.GetQueryOrFormValue(c, "text")
	if engine == "" || from == "" || to == "" || text == "" {
		return fiber.NewError(fiber.StatusBadRequest, "from, to, engine, text are required query strings.")
	}
	if engine == "all" {
		dataarr := libmozhi.TranslateAll(to, from, text)
		return c.JSON(dataarr)
	} else if strings.Contains(engine, ",") == true {
		engineArr := strings.Split(engine, ",")
		dataarr, err := libmozhi.TranslateSome(engineArr, to, from, text)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(dataarr)
	} else {
		data, err := libmozhi.Translate(engine, to, from, text)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		return c.JSON(data)
	}
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
