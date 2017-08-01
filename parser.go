// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

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

func newModel() *Model {
	m := new(Model)
	m.personas = make(map[string]Persona)
	m.systems = make(map[string]System)
	m.containers = make(map[string]Container)
	m.components = make(map[string]Component)
	m.relationships = make([]Relationship, 0)
	m.errors = make([]error, 0)
	return m
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
		default:
		}

	}

	return nil
}

func parsePersona(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 3 {
		addErr(m, path, lineno, "Persona requires 3 elements: Name | Description | Tags")
		return
	}
	if len(fields) > 3 {
		addErr(m, path, lineno, "Persona requires 3 elements: Name | Description | Tags")
	}

	name := strings.TrimSpace(fields[0])
	description := strings.TrimSpace(fields[1])
	tags := parseTags(fields[2])

	if _, ok := m.personas[name]; ok {
		addErr(m, path, lineno, "Persona is already defined: "+name)
		return
	}
	m.personas[name] = Persona{Name: name, Description: description, Tags: tags}
}

func parseSystem(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 3 {
		addErr(m, path, lineno, "System requires 3 elements: Name | Description | Tags")
		return
	}
	if len(fields) > 3 {
		addErr(m, path, lineno, "System requires 3 elements: Name | Description | Tags")
	}

	name := strings.TrimSpace(fields[0])
	description := strings.TrimSpace(fields[1])
	tags := parseTags(fields[2])

	if _, ok := m.systems[name]; ok {
		addErr(m, path, lineno, "System is already defined: "+name)
		return
	}
	m.systems[name] = System{Name: name, Description: description, Tags: tags}
}

func parseContainer(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 5 {
		addErr(m, path, lineno, "Container requires 5 elements: System | Name | Description | Technology | Tags")
		return
	}
	if len(fields) > 5 {
		addErr(m, path, lineno, "Container requires 5 elements: System | Name | Description | Technology | Tags")
	}

	system := strings.TrimSpace(fields[0])
	name := strings.TrimSpace(fields[1])
	description := strings.TrimSpace(fields[2])
	technology := strings.TrimSpace(fields[3])
	tags := parseTags(fields[4])

	if _, ok := m.containers[name]; ok {
		addErr(m, path, lineno, "Container is already defined: "+name)
		return
	}
	m.containers[name] = Container{Name: name, System: system, Description: description, Technology: technology, Tags: tags}
}

func parseComponent(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 5 {
		addErr(m, path, lineno, "Component requires 5 elements: Container | Name | Description | Technology | Tags")
		return
	}
	if len(fields) > 5 {
		addErr(m, path, lineno, "Component requires 5 elements: Container | Name | Description | Technology | Tags")
	}

	container := strings.TrimSpace(fields[0])
	name := strings.TrimSpace(fields[1])
	description := strings.TrimSpace(fields[2])
	technology := strings.TrimSpace(fields[3])
	tags := parseTags(fields[4])

	if _, ok := m.components[name]; ok {
		addErr(m, path, lineno, "Component is already defined: "+name)
		return
	}
	m.components[name] = Component{Name: name, Container: container, Description: description, Technology: technology, Tags: tags}
}

func parseRelationship(m *Model, path string, lineno int, value string) {
	fields := strings.Split(value, "|")
	if len(fields) < 5 {
		addErr(m, path, lineno, "Relationship requires 5 elements: Source | Description | Technology | Destination | Tags")
		return
	}
	if len(fields) > 5 {
		addErr(m, path, lineno, "Relationship requires 5 elements: Source | Description | Technology | Destination | Tags")
	}

	source := strings.TrimSpace(fields[0])
	description := strings.TrimSpace(fields[1])
	technology := strings.TrimSpace(fields[2])
	destination := strings.TrimSpace(fields[3])
	tags := parseTags(fields[4])

	m.relationships = append(m.relationships,
		Relationship{Source: source, Description: description, Technology: technology, Destination: destination, Tags: tags})
}

func parseTags(s string) []string {
	tags := strings.Split(s, ",")
	for i, tag := range tags {
		tags[i] = strings.TrimSpace(tag)
	}
	return tags
}

func addErr(m *Model, path string, lineno int, msg string) {
	err := ParseError{File: path, Line: lineno, Msg: msg}
	m.errors = append(m.errors, err)
}
