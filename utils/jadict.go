package utils

import (
	"strings"

	"codeberg.org/aryak/libmozhi"
)

// JaDictEnabled reports whether id is a Japanese dictionary source (jisho,
// jpdb, moji, weblio) that is on: all are, unless MOZHI_<ID>_ENABLED=false.
func JaDictEnabled(id string) bool {
	for _, s := range libmozhi.JaDictSources() {
		if s == id {
			return EnvTrueNoExist("MOZHI_" + strings.ToUpper(id) + "_ENABLED")
		}
	}
	return false
}

// JaDictSources returns the enabled Japanese dictionary sources, in display
// order.
func JaDictSources() []string {
	out := []string{}
	for _, s := range libmozhi.JaDictSources() {
		if JaDictEnabled(s) {
			out = append(out, s)
		}
	}
	return out
}
