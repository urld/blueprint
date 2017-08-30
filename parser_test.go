// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"bufio"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func TestParseSystem(t *testing.T) {
	m := newModel()
	path := "test/parsesystem"
	lineno := 1
	value := " Test System | Test Description | tag1,tag2"

	parseSystem(m, path, lineno, value)

	assertEqual(t, 0, len(m.Errors), "0 errors expected")
	assertEqual(t, 1, len(m.Systems), "1 system expected")
	sys := m.Systems["Test System"]
	expectedSys := System{Name: "Test System", Description: "Test Description", Tags: []string{"tag1", "tag2"}}
	assertEqual(t, sys, expectedSys, "system content does not match")
}

func TestParseSystemMissingElemet(t *testing.T) {
	m := newModel()
	path := "test/parsesystem"
	value := " Test System | Test Description"

	parseSystem(m, path, 1, value)

	assertEqual(t, 1, len(m.Errors), "1 error expected")
	expectedErr := parseError{File: path, Line: 1, Msg: "System requires 3 elements: Name | Description | Tags"}
	assertEqual(t, expectedErr, m.Errors[0], "error does not match")

	assertEqual(t, 0, len(m.Systems), "0 systems expected")
}

func TestParseSystemDuplicate(t *testing.T) {
	m := newModel()
	path := "test/parsesystem"
	value := " Test System | Test Description | tag1,tag2"

	parseSystem(m, path, 1, value)
	parseSystem(m, path, 2, value)

	assertEqual(t, 1, len(m.Errors), "1 error expected")
	expectedErr := parseError{File: path, Line: 2, Msg: "System is already defined: Test System"}
	assertEqual(t, expectedErr, m.Errors[0], "error does not match")

	assertEqual(t, 1, len(m.Systems), "1 system expected")
	sys := m.Systems["Test System"]
	expectedSys := System{Name: "Test System", Description: "Test Description", Tags: []string{"tag1", "tag2"}}
	assertEqual(t, sys, expectedSys, "system content does not match")
}

func TestParseMultiLineSystem(t *testing.T) {
	rawValue := " Test System | Test Description spans\\\n multiple \\\nlines | tag1,tag2\nother line"
	s := bufio.NewScanner(strings.NewReader(rawValue))
	s.Scan()
	value, lineCnt := parseLine(s)
	expectedValue := "Test System | Test Description spans multiple lines | tag1,tag2"
	assertEqual(t, value, expectedValue, "multiline system text does not match")
	assertEqual(t, lineCnt, 3, "lineCnt of multiline system does not match")
}

func assertEqual(t *testing.T, a, b interface{}, message string) {
	if reflect.DeepEqual(a, b) {
		return
	}
	msg := fmt.Sprintf("%s: %v != %v", message, a, b)
	t.Error(msg)
}
