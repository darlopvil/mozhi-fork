package utils

import (
	"os"
	"regexp"

	"codeberg.org/aryak/libmozhi"
	"github.com/gofiber/fiber/v2"
)

var nonAlphanumericRegex = regexp.MustCompile(`[^a-zA-Z,]+`)
var nonAlphaRegex = regexp.MustCompile(`[^a-zA-Z0-9,]+`)

func GetQueryOrFormValue(c *fiber.Ctx, key string) string {
	if c.Method() == "POST" {
		return c.FormValue(key)
	} else {
		return c.Query(key)
	}
}

func EnvTrueNoExist(env string) bool {
	_, envFound := os.LookupEnv(env)

	return !envFound || os.Getenv(env) == "true"
}

func Sanitize(str string, strip string) string {
	if strip == "alpha" {
		return nonAlphaRegex.ReplaceAllString(str, "")
	} else if strip == "alphanumeric" {
		return nonAlphanumericRegex.ReplaceAllString(str, "")
	}
	return ""
}

func EngineList() map[string]string {
	engines := map[string]string{"all": "All Engines", "some": "Some Engines", "google": "Google", "deepl": "DeepL", "duckduckgo": "DuckDuckGo", "gemini": "Gemini", "mymemory": "MyMemory", "textra": "TexTra (NICT)", "yandex": "Yandex", "groq": "Groq (LLM)","gptoss":   "OpenRouter (GPT-OSS)",
"gemma":    "OpenRouter (Gemma)",
"nemotron": "OpenRouter (Nemotron)",}
	if EnvTrueNoExist("MOZHI_GOOGLE_ENABLED") == false {
		delete(engines, "google")
	}
	if EnvTrueNoExist("MOZHI_DEEPL_ENABLED") == false {
		delete(engines, "deepl")
	}
	if EnvTrueNoExist("MOZHI_DUCKDUCKGO_ENABLED") == false {
		delete(engines, "duckduckgo")
	}
	if EnvTrueNoExist("MOZHI_MYMEMORY_ENABLED") == false {
		delete(engines, "mymemory")
	}
	if EnvTrueNoExist("MOZHI_YANDEX_ENABLED") == false {
		delete(engines, "yandex")
	}
	if EnvTrueNoExist("MOZHI_GEMINI_ENABLED") == false {
		delete(engines, "gemini")
	}
	if EnvTrueNoExist("MOZHI_TEXTRA_ENABLED") == false {
		delete(engines, "textra")
	}
	return engines
}

// DeduplicateLists deduplicates a slice of List based on the Id field
func DeDuplicateLists(input []libmozhi.List) []libmozhi.List {
	// Create a map to store unique Ids
	uniqueIds := make(map[string]struct{})
	result := []libmozhi.List{}

	// Iterate over the input slice
	for _, item := range input {
		// Check if the Id is unique
		if _, found := uniqueIds[item.Id]; !found {
			// Add the Id to the map and append the List to the result slice
			uniqueIds[item.Id] = struct{}{}
			result = append(result, item)
		}
	}

	return result
}
