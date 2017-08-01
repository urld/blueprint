// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

import (
	"bytes"
	"strings"
)

func WrapWords(text string, lineSep string, limit int) string {
	words := strings.Fields(text)
	var buf bytes.Buffer
	remaining := limit
	for i, word := range words {
		if i == 0 {
			// first word is special
		} else if len(word)+1 > remaining {
			buf.WriteString(lineSep)
			remaining = limit
		} else {
			buf.WriteRune(' ')
			remaining--
		}
		buf.WriteString(word)
		remaining -= len(word)
	}
	return buf.String()
}
