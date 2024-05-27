package pages

import (
	"os"
	"slices"

	"codeberg.org/aryak/libmozhi"
	"codeberg.org/aryak/mozhi/utils"
	"github.com/gofiber/fiber/v2"
)

func langListMerge(engines map[string]string) ([]libmozhi.List, []libmozhi.List) {
	sl := []libmozhi.List{}
	tl := []libmozhi.List{}
	for key := range engines {
		temp, _ := libmozhi.LangList(key, "sl")
		temp2, _ := libmozhi.LangList(key, "tl")
		sl = append(sl, temp...)
		tl = append(tl, temp2...)
	}
	return utils.DeDuplicateLists(sl), utils.DeDuplicateLists(tl)
}

func HandleIndex(c *fiber.Ctx) error {
	engines := utils.EngineList()
	var enginesAsArray []string
	for engine := range engines {
		enginesAsArray = append(enginesAsArray, engine)
	}
	engineCookie := c.Cookies("engine")
	fromCookie := c.Cookies("from")
	toCookie := c.Cookies("to")
	var engine = utils.GetQueryOrFormValue(c, "engine")
	if engine == "" || !slices.Contains(enginesAsArray, engine) {
		if engineCookie != "" {
			engine = engineCookie
		} else {
			engine = "google"
		}
	}

	var sourceLanguages []libmozhi.List
	var targetLanguages []libmozhi.List
	if engine == "all" {
		sourceLanguages, targetLanguages = langListMerge(engines)
	} else {
		sourceLanguages, _ = libmozhi.LangList(engine, "sl")
		targetLanguages, _ = libmozhi.LangList(engine, "tl")
	}

	originalText := utils.GetQueryOrFormValue(c, "text")
	to := utils.GetQueryOrFormValue(c, "to")
	from := utils.GetQueryOrFormValue(c, "from")
	if from == "" {
		from = utils.GetQueryOrFormValue(c, "sl")
		if from == "" && fromCookie != "" {
			from = fromCookie
		}
	}
	if to == "" {
		to = utils.GetQueryOrFormValue(c, "tl")
		if to == "" && toCookie != "" {
			to = toCookie
		}
	}

	var translation libmozhi.LangOut
	var translationExists bool
	var transall []libmozhi.LangOut
	var tlerr error
	var ttsFrom string
	var ttsTo string
	if engine != "" && originalText != "" && from != "" && to != "" {
		if engine == "all" {
			transall = libmozhi.TranslateAll(to, from, originalText)
		} else {
			translation, tlerr = libmozhi.Translate(engine, to, from, originalText)
			if tlerr != nil {
				return fiber.NewError(fiber.StatusInternalServerError, tlerr.Error())
			}
			translationExists = true
		}
		if engine == "google" || engine == "reverso" {
			if from == "auto" && translation.AutoDetect != "" {
				ttsFrom = "/api/tts?lang=" + translation.AutoDetect + "&engine=" + engine + "&text=" + originalText
			} else {
				ttsFrom = "/api/tts?lang=" + from + "&engine=" + engine + "&text=" + originalText
			}
			ttsTo = "/api/tts?lang=" + to + "&engine=" + engine + "&text=" + translation.OutputText
		}
	} else {
		translationExists = false
	}

	defaultLang := os.Getenv("MOZHI_DEFAULT_SOURCE_LANG")
	preferAutoDetect := os.Getenv("MOZHI_DEFAULT_PREFER_AUTODETECT")
	defaultLangTarget := os.Getenv("MOZHI_DEFAULT_TARGET_LANG")
	if defaultLang == "" || preferAutoDetect == "true" {
		defaultLang = "auto"
	}
	if defaultLangTarget == "" {
		defaultLangTarget = "en"
	}
	cookie := new(fiber.Cookie)
	cookie.Name = "engine"
	cookie.Value = engine
	c.Cookie(cookie)
	if from != "" {
		cookie := new(fiber.Cookie)
		cookie.Name = "from"
		cookie.Value = from
		c.Cookie(cookie)
	}
	if to != "" {
		cookie := new(fiber.Cookie)
		cookie.Name = "to"
		cookie.Value = to
		c.Cookie(cookie)
	}
	return c.Render("index", fiber.Map{
		"Engine":            engine,
		"enginesNames":      engines,
		"SourceLanguages":   sourceLanguages,
		"TargetLanguages":   targetLanguages,
		"OriginalText":      originalText,
		"Translation":       translation,
		"TranslationExists": translationExists,
		"TranslateAll":      transall,
		"From":              from,
		"To":                to,
		"TtsFrom":           ttsFrom,
		"TtsTo":             ttsTo,
		"defaultLang":       defaultLang,
		"defaultLangTarget": defaultLangTarget,
	})
}
