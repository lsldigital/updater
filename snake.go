package updater

import (
	"strings"

	"github.com/martingallagher/runes"
)

// Credit: https://github.com/azer/snakecase
func toSnakeCase(str string) string {
	var sb strings.Builder

	in := []rune(str)
	lenIn := len(in)

	isLower := func(idx int) bool {
		return idx >= 0 && idx < lenIn && runes.IsLower(in[idx])
	}

	for i, r := range in {
		if runes.IsSpace(r) {
			if i+1 < lenIn {
				in[i+1] = runes.ToUpper(in[i+1])
			}
			continue
		}
		if runes.IsUpper(r) {
			r = runes.ToLower(r)
			if i > 0 && in[i-1] != '_' && (isLower(i-1) || isLower(i+1)) {
				sb.WriteRune('_')
			}
		}
		sb.WriteRune(r)
	}

	return sb.String()
}
