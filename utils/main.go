package utils

import (
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode/utf8"

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
	engines := map[string]string{"all": "All Engines", "some": "Some Engines", "google": "Google", "deepl": "DeepL", "duckduckgo": "DuckDuckGo", "gemini": "Gemini", "mymemory": "MyMemory", "textra": "TexTra (NICT)", "yandex": "Yandex", "groq": "Groq (LLM)", "qwen": "Qwen (Groq)", "cohere": "Cohere", "mistral": "Mistral"}
	// Every engine can be turned off with MOZHI_<ID>_ENABLED=false (on by default).
	for id := range engines {
		if id == "all" || id == "some" {
			continue
		}
		if !EnvTrueNoExist("MOZHI_" + strings.ToUpper(id) + "_ENABLED") {
			delete(engines, id)
		}
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

// EnabledEngines returns the ids of the enabled translation engines, without
// the "all" and "some" pseudo-engines. All mode translates with exactly these.
func EnabledEngines() []string {
	engines := EngineList()
	ids := []string{}
	for id := range engines {
		if id != "all" && id != "some" {
			ids = append(ids, id)
		}
	}
	// Stable order, by display name like the dropdown: map order is random
	// and All mode returns its results in this order.
	sort.Slice(ids, func(i, j int) bool { return engines[ids[i]] < engines[ids[j]] })
	return ids
}

// EngineEnabled reports whether id is a translation engine enabled on this
// instance. Engines that are disabled, or not listed at all, are refused by
// the web UI and the API alike.
func EngineEnabled(id string) bool {
	if id == "all" || id == "some" {
		return false
	}
	_, ok := EngineList()[id]
	return ok
}

// MaxTextLength is the longest text (in characters) accepted for translation
// or TTS: MOZHI_MAX_TEXT_LENGTH, default 5000 (the limit of Google's own web
// translator), 0 for no limit. In All mode one request fans out to every
// enabled engine, so an unbounded text multiplies quota use across all of them.
func MaxTextLength() int {
	if n, err := strconv.Atoi(os.Getenv("MOZHI_MAX_TEXT_LENGTH")); err == nil && n >= 0 {
		return n
	}
	return 5000
}

// TextTooLong reports whether text exceeds MaxTextLength.
func TextTooLong(text string) bool {
	max := MaxTextLength()
	return max > 0 && utf8.RuneCountInString(text) > max
}

// TextTooLongError is the error returned for a text over the limit.
func TextTooLongError() error {
	return fiber.NewError(fiber.StatusRequestEntityTooLarge, "text is longer than "+strconv.Itoa(MaxTextLength())+" characters")
}

// ReadBufferSize is Fiber's per-connection read buffer, which caps the
// request line plus headers (default 4096 bytes: GET requests over ~4 KB got
// 431). It is sized for a MaxTextLength text in the URL of a TTS link: up to
// 12 bytes per character once percent-encoded, plus 8 KB for headers.
func ReadBufferSize() int {
	if max := MaxTextLength(); max > 0 {
		return 12*max + 8192
	}
	return 1 << 20
}

// PinnedLangs returns MOZHI_PINNED_LANGS (comma-separated language codes,
// e.g. "ja,es,en"): languages listed first in the language selectors.
func PinnedLangs() []string {
	var codes []string
	for _, code := range strings.Split(os.Getenv("MOZHI_PINNED_LANGS"), ",") {
		if code = strings.TrimSpace(code); code != "" {
			codes = append(codes, code)
		}
	}
	return codes
}
