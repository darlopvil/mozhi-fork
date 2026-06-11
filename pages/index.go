package pages

import (
	"os"
	"slices"
	"strings"

	"codeberg.org/aryak/libmozhi"
	"codeberg.org/aryak/mozhi/utils"
	"github.com/gofiber/fiber/v2"
)

func langListMerge(engines []string) ([]libmozhi.List, []libmozhi.List) {
	sl := []libmozhi.List{}
	tl := []libmozhi.List{}
	for _, key := range engines {
		temp, _ := libmozhi.LangList(key, "sl")
		temp2, _ := libmozhi.LangList(key, "tl")
		sl = append(sl, temp...)
		tl = append(tl, temp2...)
	}
	return utils.DeDuplicateLists(sl), utils.DeDuplicateLists(tl)
}

func HandleIndex(c *fiber.Ctx) error {
	engines := utils.EngineList()
	type enginesStruct struct {
		Engines []string `query:"engines"`
	}
	enginesSome := new(enginesStruct)
	c.QueryParser(enginesSome)
	var enginesAsArray []string
	for engine := range engines {
		enginesAsArray = append(enginesAsArray, engine)
	}
	defaultEngine := os.Getenv("MOZHI_DEFAULT_ENGINE")
	if defaultEngine == "" || !slices.Contains(enginesAsArray, defaultEngine) {
		if slices.Contains(enginesAsArray, "google") {
			defaultEngine = "google"
		} else if len(enginesAsArray) > 0 {
			defaultEngine = enginesAsArray[0]
		}
	}
	engineCookie := c.Cookies("engine")
	fromCookie := c.Cookies("from")
	toCookie := c.Cookies("to")
	var engine = utils.GetQueryOrFormValue(c, "engine")
	if engine == "" || !slices.Contains(enginesAsArray, engine) {
		if engineCookie != "" {
			engine = engineCookie
		} else {
			engine = defaultEngine
		}
	}

	var sourceLanguages []libmozhi.List
	var targetLanguages []libmozhi.List
	if engine == "all" {
		sourceLanguages, targetLanguages = langListMerge(enginesAsArray)
	} else if engine == "some" {
		if c.Cookies("engines") == "" && enginesSome.Engines == nil {
			enginesSome.Engines = append(enginesSome.Engines, "google")
		} else if enginesSome.Engines == nil && c.Cookies("engines") != "" {
			enginesSome.Engines = strings.Split(c.Cookies("engines"), ",")
		}
		for i, engine := range enginesSome.Engines {
			if !slices.Contains(enginesAsArray, engine) || engine == "some" {
				// Delete array from slice if its not in engines list
				enginesSome.Engines = append(enginesSome.Engines[:i], enginesSome.Engines[i+1:]...)
			} else if engine == "all" {
				enginesSome.Engines = enginesAsArray
			}
		}
		if enginesSome.Engines == nil {
			enginesSome.Engines = append(enginesSome.Engines, "google")
		}
		sourceLanguages, targetLanguages = langListMerge(enginesSome.Engines)
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
	var transmany []libmozhi.LangOut
	var tlerr error
	var ttsFrom string
	var ttsTo string
	if engine != "" && originalText != "" && from != "" && to != "" {
		if engine == "all" {
			transmany = libmozhi.TranslateAll(to, from, originalText)
			translationExists = true
		} else if engine == "some" {
			// The error doesn't really matter since it just checks if the engines are valid, which is already checked in the code at the beginning of the func
			transmany, _ = libmozhi.TranslateSome(enginesSome.Engines, to, from, originalText)
			translationExists = true
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
	if defaultLang == "" || (preferAutoDetect == "true" && len(sourceLanguages) > 0 && sourceLanguages[0].Id == "auto") {
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
	if enginesSome.Engines != nil {
		cookie := new(fiber.Cookie)
		cookie.Name = "engines"
		cookie.Value = strings.Join(enginesSome.Engines, ",")
		c.Cookie(cookie)
	}
	return c.Render("index", fiber.Map{
		"Engine":            engine,
		"enginesNames":      engines,
		"SomeEngines":       enginesSome.Engines,
		"SomeEnginesStr":    strings.Join(enginesSome.Engines, ","),
		"SourceLanguages":   sourceLanguages,
		"TargetLanguages":   targetLanguages,
		"OriginalText":      originalText,
		"Translation":       translation,
		"TranslationExists": translationExists,
		"TranslateMany":     transmany,
		"From":              from,
		"To":                to,
		"TtsFrom":           ttsFrom,
		"TtsTo":             ttsTo,
		"defaultLang":       defaultLang,
		"defaultLangTarget": defaultLangTarget,
	})
}
