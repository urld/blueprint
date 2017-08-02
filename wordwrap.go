// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"bytes"
	"strings"
)

func wrapWords(text string, limit int) string {
	words := strings.Fields(text)
	var buf bytes.Buffer
	remaining := limit
	for i, word := range words {
		switch {
		case i == 0:
			// first word is special
		case len(word)+1 > remaining:
			buf.WriteString("<BR/>")
			remaining = limit
		default:
			buf.WriteRune(' ')
			remaining--
		}
		buf.WriteString(word)
		remaining -= len(word)
	}
	return buf.String()
}
