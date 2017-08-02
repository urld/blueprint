// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package blueprint

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type ParseError struct {
	Msg  string
	File string
	Line int
}

func (e ParseError) Error() string {
	return fmt.Sprintf("%s:%d: %s", e.File, e.Line, e.Msg)
}

func Parse(path string) (Model, error) {
	m := newModel()
	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil

		}
		return parseFile(path, m)
	})

	return *m, err
}

func parseFile(path string, m *Model) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	lineno := 0
	s := bufio.NewScanner(file)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		lineno++
		if line == "" {
			continue
		}

		i := strings.Index(line, "=")
		if i == -1 {
			// TODO error?
			continue
		}
		key := strings.TrimSpace(line[:i])
		value := strings.TrimSpace(line[i+1:])
		switch key {
		case "Persona", "Person":
			parsePersona(m, path, lineno, value)
		case "System", "SoftwareSystem":
			parseSystem(m, path, lineno, value)
		case "Container":
			parseContainer(m, path, lineno, value)
		case "Component":
			parseComponent(m, path, lineno, value)
		case "Relationship":
			parseRelationship(m, path, lineno, value)
		case "SystemContext":
			parseSystemContext(m, path, lineno, value)
		default:
			m.addErr(path, lineno, "unknown keyword: "+key)
		}

	}

	return nil
}

func parsePersona(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 3 {
		m.addErr(path, lineno, "Persona requires 3 elements: Name | Description | Tags")
		return
	}
	if len(fields) > 3 {
		m.addErr(path, lineno, "Persona requires 3 elements: Name | Description | Tags")
	}

	name := strings.TrimSpace(fields[0])
	description := strings.TrimSpace(fields[1])
	tags := parseTags(fields[2])

	if _, ok := m.Personas[name]; ok {
		m.addErr(path, lineno, "Persona is already defined: "+name)
		return
	}
	m.Personas[name] = Persona{Name: name, Description: description, Tags: tags}
}

func parseSystem(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 3 {
		m.addErr(path, lineno, "System requires 3 elements: Name | Description | Tags")
		return
	}
	if len(fields) > 3 {
		m.addErr(path, lineno, "System requires 3 elements: Name | Description | Tags")
	}

	name := strings.TrimSpace(fields[0])
	description := strings.TrimSpace(fields[1])
	tags := parseTags(fields[2])

	if _, ok := m.Systems[name]; ok {
		m.addErr(path, lineno, "System is already defined: "+name)
		return
	}
	m.Systems[name] = System{Name: name, Description: description, Tags: tags}
}

func parseContainer(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 5 {
		m.addErr(path, lineno, "Container requires 5 elements: System | Name | Description | Technology | Tags")
		return
	}
	if len(fields) > 5 {
		m.addErr(path, lineno, "Container requires 5 elements: System | Name | Description | Technology | Tags")
	}

	system := strings.TrimSpace(fields[0])
	name := strings.TrimSpace(fields[1])
	description := strings.TrimSpace(fields[2])
	technology := strings.TrimSpace(fields[3])
	tags := parseTags(fields[4])

	if _, ok := m.Containers[name]; ok {
		m.addErr(path, lineno, "Container is already defined: "+name)
		return
	}
	m.Containers[name] = Container{Name: name, System: system, Description: description, Technology: technology, Tags: tags}
}

func parseComponent(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 5 {
		m.addErr(path, lineno, "Component requires 5 elements: Container | Name | Description | Technology | Tags")
		return
	}
	if len(fields) > 5 {
		m.addErr(path, lineno, "Component requires 5 elements: Container | Name | Description | Technology | Tags")
	}

	container := strings.TrimSpace(fields[0])
	name := strings.TrimSpace(fields[1])
	description := strings.TrimSpace(fields[2])
	technology := strings.TrimSpace(fields[3])
	tags := parseTags(fields[4])

	if _, ok := m.Components[name]; ok {
		m.addErr(path, lineno, "Component is already defined: "+name)
		return
	}
	m.Components[name] = Component{Name: name, Container: container, Description: description, Technology: technology, Tags: tags}
}

func parseRelationship(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 5 {
		m.addErr(path, lineno, "Relationship requires 5 elements: Source | Description | Technology | Destination | Tags")
		return
	}
	if len(fields) > 5 {
		m.addErr(path, lineno, "Relationship requires 5 elements: Source | Description | Technology | Destination | Tags")
	}

	source := strings.TrimSpace(fields[0])
	description := strings.TrimSpace(fields[1])
	technology := strings.TrimSpace(fields[2])
	destination := strings.TrimSpace(fields[3])
	tags := parseTags(fields[4])

	m.Relationships = append(m.Relationships,
		Relationship{Source: source, Description: description, Technology: technology, Destination: destination, Tags: tags})
}

func parseSystemContext(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 4 {
		m.addErr(path, lineno, "SystemContext requires 4 elements: CoreSystems | ExternalSystems | Name | Description")
		return
	}
	if len(fields) > 4 {
		m.addErr(path, lineno, "SystemContext requires 4 elements: CoreSystems | ExternalSystems | Name | Description")
	}

	coreSys := parseTags(fields[0])
	extSys := parseTags(fields[1])
	name := strings.TrimSpace(fields[2])
	description := strings.TrimSpace(fields[3])

	if _, ok := m.SystemContexts[name]; ok {
		m.addErr(path, lineno, "View is already defined: "+name)
		return
	}
	m.SystemContexts[name] = SystemContext{Name: name, Description: description, CoreSystems: coreSys, ExternalSystems: extSys}
}

func parseTags(s string) []string {
	tags := strings.Split(s, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags
}

func (m *Model) addErr(path string, lineno int, msg string) {
	err := ParseError{File: path, Line: lineno, Msg: msg}
	m.Errors = append(m.Errors, err)
}
