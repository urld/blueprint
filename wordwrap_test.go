// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"testing"
)

func TestWrapLastWord(t *testing.T) {
	wrapped := wrapWords("Foo Bar Baz", 9)
	assertEqual(t, "Foo Bar<BR/>Baz", wrapped, "")
}

func TestWrapLastWordExact(t *testing.T) {
	wrapped := wrapWords("Foo Bar Baz", 8)
	assertEqual(t, "Foo Bar<BR/>Baz", wrapped, "")
}

func TestWrapLastTwoWords(t *testing.T) {
	wrapped := wrapWords("Foo Bar Ba", 6)
	assertEqual(t, "Foo<BR/>Bar Ba", wrapped, "")
}

func TestWrapMultipleWords(t *testing.T) {
	wrapped := wrapWords("Foo Bar Baz", 6)
	assertEqual(t, "Foo<BR/>Bar<BR/>Baz", wrapped, "")
}

func TestDontWrapTooLongWord(t *testing.T) {
	wrapped := wrapWords("Short TooLong", 7)
	assertEqual(t, "Short<BR/>TooLong", wrapped, "")
}

func TestDontWrapTooLongWordAtStart(t *testing.T) {
	wrapped := wrapWords("TooLong", 7)
	assertEqual(t, "TooLong", wrapped, "")
}
