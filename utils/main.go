package utils

import (
	"codeberg.org/aryak/libmozhi"
	"github.com/gofiber/fiber/v2"
	"os"
	"regexp"
)

func GetQueryOrFormValue(c *fiber.Ctx, key string) string {
	if c.Method() == "POST" {
		return c.FormValue(key)
	} else {
		return c.Query(key)
	}
}

func EnvTrueNoExist(env string) bool {
	if _, ok := os.LookupEnv(env); ok == false || os.Getenv(env) == "true" {
		return true
	}
	return false
}

func Sanitize(str string, strip string) string {
	nonAlphanumericRegex := regexp.MustCompile(`[^a-zA-Z]+`)
	nonAlphaRegex := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	if strip == "alpha" {
		return nonAlphaRegex.ReplaceAllString(str, "")
	} else if strip == "alphanumeric" {
		return nonAlphanumericRegex.ReplaceAllString(str, "")
	}
	return ""
}

func EngineList() map[string]string {
	engines := map[string]string{"all": "All Engines", "google": "Google", "deepl": "DeepL", "duckduckgo": "DuckDuckGo", "libre": "LibreTranslate", "mymemory": "MyMemory", "reverso": "Reverso", "watson": "Watson", "yandex": "Yandex"}
	if EnvTrueNoExist("MOZHI_GOOGLE_ENABLED") == false {
		delete(engines, "google")
	} else if EnvTrueNoExist("MOZHI_DEEPL_ENABLED") == false {
		delete(engines, "deepl")
	} else if EnvTrueNoExist("MOZHI_DUCKDUCKGO_ENABLED") == false {
		delete(engines, "duckduckgo")
	} else if EnvTrueNoExist("MOZHI_LIBRETRANSLATE_ENABLED") == false || EnvTrueNoExist("MOZHI_LIBRETRANSLATE_URL") {
		delete(engines, "libre")
	} else if EnvTrueNoExist("MOZHI_MYMEMORY_ENABLED") == false {
		delete(engines, "mymemory")
	} else if EnvTrueNoExist("MOZHI_REVERSO_ENABLED") == false {
		delete(engines, "reverso")
	} else if EnvTrueNoExist("MOZHI_WATSON_ENABLED") == false {
		delete(engines, "watson")
	} else if EnvTrueNoExist("MOZHI_YANDEX_ENABLED") == false {
		delete(engines, "yandex")
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
