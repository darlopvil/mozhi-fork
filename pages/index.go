package pages

import (
	"net/url"
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

	// Orden del dropdown: all/some arriba fijos, resto alfabético por nombre visible.
	enginesOrder := []string{}
	for key := range engines {
		if key == "all" || key == "some" {
			continue
		}
		enginesOrder = append(enginesOrder, key)
	}
	slices.SortFunc(enginesOrder, func(a, b string) int {
		return strings.Compare(engines[a], engines[b])
	})
	enginesPrefix := []string{}
	if _, ok := engines["all"]; ok {
		enginesPrefix = append(enginesPrefix, "all")
	}
	if _, ok := engines["some"]; ok {
		enginesPrefix = append(enginesPrefix, "some")
	}
	enginesOrder = append(enginesPrefix, enginesOrder...)
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
		if engineCookie != "" && slices.Contains(enginesAsArray, engineCookie) {
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
		filtered := enginesSome.Engines[:0]
		for _, engine := range enginesSome.Engines {
			if engine == "all" {
				filtered = enginesAsArray
				break
			}
			if slices.Contains(enginesAsArray, engine) && engine != "some" {
				filtered = append(filtered, engine)
			}
		}
		enginesSome.Engines = filtered
		if enginesSome.Engines == nil {
			enginesSome.Engines = append(enginesSome.Engines, "google")
		}
		sourceLanguages, targetLanguages = langListMerge(enginesSome.Engines)
	} else {
		sourceLanguages, _ = libmozhi.LangList(engine, "sl")
		targetLanguages, _ = libmozhi.LangList(engine, "tl")
	}

	originalText := utils.GetQueryOrFormValue(c, "text")
	if utils.TextTooLong(originalText) {
		return utils.TextTooLongError()
	}
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
	// Each engine may name a language its own way (es / es-ES, zh-CN /
	// zh-Hans): match the choice to this engine's list, and fall back to the
	// configured defaults, never to the first option of the list.
	defaultLang, defaultLangTarget := defaultLangs(sourceLanguages, targetLanguages)
	if from = libmozhi.ResolveLang(sourceLanguages, from); from == "" {
		from = defaultLang
	}
	if to = libmozhi.ResolveLang(targetLanguages, to); to == "" {
		to = defaultLangTarget
	}

	var translation libmozhi.LangOut
	var translationExists bool
	var transmany []libmozhi.LangOut
	var tlerr error
	var ttsFrom string
	var ttsTo string
	if engine != "" && originalText != "" && from != "" && to != "" {
		if engine == "all" {
			// Only the engines enabled on this instance.
			transmany, _ = libmozhi.TranslateSome(utils.EnabledEngines(), to, from, originalText)
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
		if engine == "google" || engine == "reverso" || engine == "textra" || engine == "groq" {
			if from == "auto" && translation.AutoDetect != "" {
				ttsFrom = ttsURL(translation.AutoDetect, engine, originalText)
			} else {
				ttsFrom = ttsURL(from, engine, originalText)
			}
			ttsTo = ttsURL(to, engine, translation.OutputText)
		}
	} else {
		translationExists = false
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
		"enginesOrder":      enginesOrder,
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

// ttsURL builds the /api/tts link for the audio players. Every value is
// query-escaped: a text containing &, #, + or % used to be cut or corrupted.
func ttsURL(lang, engine, text string) string {
	return "/api/tts?lang=" + url.QueryEscape(lang) + "&engine=" + url.QueryEscape(engine) + "&text=" + url.QueryEscape(text)
}

// defaultLangs returns the configured default source and target languages
// (MOZHI_DEFAULT_SOURCE_LANG, MOZHI_DEFAULT_PREFER_AUTODETECT,
// MOZHI_DEFAULT_TARGET_LANG) matched to the codes of these lists.
func defaultLangs(source, target []libmozhi.List) (string, string) {
	from := os.Getenv("MOZHI_DEFAULT_SOURCE_LANG")
	if from == "" || (os.Getenv("MOZHI_DEFAULT_PREFER_AUTODETECT") == "true" && len(source) > 0 && source[0].Id == "auto") {
		from = "auto"
	}
	to := os.Getenv("MOZHI_DEFAULT_TARGET_LANG")
	if to == "" {
		to = "en"
	}
	return libmozhi.ResolveLang(source, from), libmozhi.ResolveLang(target, to)
}
