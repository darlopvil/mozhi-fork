package pages

import (
	"strings"

	"codeberg.org/aryak/libmozhi"
	"codeberg.org/aryak/mozhi/utils"
	"github.com/gofiber/fiber/v2"
)

// jaDictView is the Japanese dictionary panel under a translation: the word
// looked up and the sources that have it.
type jaDictView struct {
	Word    string
	Results []libmozhi.JaDict
}

func isJapanese(lang string) bool {
	l := strings.ToLower(lang)
	return l == "ja" || strings.HasPrefix(l, "ja-")
}

// jaDictLookup looks word up in every enabled source, in the background.
func jaDictLookup(word string) <-chan jaDictView {
	sources := utils.JaDictSources()
	if len(sources) == 0 {
		return nil
	}
	ch := make(chan jaDictView, 1)
	go func() {
		v := jaDictView{Word: word}
		for _, r := range libmozhi.JapaneseDictionary(word, sources) {
			if len(r.Entries) > 0 {
				v.Results = append(v.Results, r)
			}
		}
		ch <- v
	}()
	return ch
}

// jaDictStart looks up a Japanese source text right away, while it is being
// translated.
func jaDictStart(from, text string) <-chan jaDictView {
	if !isJapanese(from) {
		return nil
	}
	if w, ok := libmozhi.JaDictWord(text); ok {
		return jaDictLookup(w)
	}
	return nil
}

// jaDictFinish returns the panel: the lookup jaDictStart started, or one of
// the Japanese side known only after translating: a detected source, or the
// translation into Japanese of a single engine.
func jaDictFinish(started <-chan jaDictView, from, to, text string, tr libmozhi.LangOut, single bool) jaDictView {
	if started == nil {
		w, ok := "", false
		switch {
		case from == "auto" && isJapanese(tr.AutoDetect):
			w, ok = libmozhi.JaDictWord(text)
		case single && !isJapanese(from) && isJapanese(to):
			w, ok = libmozhi.JaDictWord(tr.OutputText)
		}
		if !ok {
			return jaDictView{}
		}
		if started = jaDictLookup(w); started == nil {
			return jaDictView{}
		}
	}
	return <-started
}

// HandleJaDict godoc
//
//	@Summary		Look up a Japanese word in the Japanese dictionaries
//	@Description	Jisho, jpdb, MOJi and Weblio, at once. A source that doesn't have the word answers no entries; one that fails answers its error.
//	@Param			word	query	string	true	"Japanese word or short phrase (at most 15 characters)"
//	@Param			sources	query	string	false	"Comma-separated sources (default: every enabled one)"
//	@Success		200		{array}	libmozhi.JaDict
//	@Router			/api/jadict [get]
func HandleJaDict(c *fiber.Ctx) error {
	word, ok := libmozhi.JaDictWord(c.Query("word"))
	if !ok {
		return fiber.NewError(fiber.StatusBadRequest, "word must be one Japanese word or a short phrase (at most 15 characters)")
	}
	sources := utils.JaDictSources()
	if s := c.Query("sources"); s != "" {
		sources = nil
		for _, id := range strings.Split(s, ",") {
			id = strings.TrimSpace(id)
			if !utils.JaDictEnabled(id) {
				return fiber.NewError(fiber.StatusBadRequest, "dictionary "+id+" is not enabled")
			}
			sources = append(sources, id)
		}
	}
	return c.JSON(libmozhi.JapaneseDictionary(word, sources))
}
