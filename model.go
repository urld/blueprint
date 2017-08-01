// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

type Model struct {
	personas      map[string]Persona
	systems       map[string]System
	containers    map[string]Container
	components    map[string]Component
	relationships []Relationship
	errors        []error
}

type Persona struct {
	Name        string
	Description string
	Tags        []string
}

type System struct {
	Name        string
	Description string
	Tags        []string
}

type Container struct {
	System      string
	Name        string
	Description string
	Technology  string
	Tags        []string
}

type Component struct {
	Container   string
	Name        string
	Description string
	Technology  string
	Tags        []string
}

type Relationship struct {
	Source      string
	Description string
	Technology  string
	Destination string
	Tags        []string
}

func (m Model) Relationships(source string) []Relationship {
	rels := make([]Relationship, 0)
	for _, r := range m.relationships {
		if r.Source != source {
			continue
		}
		rels = append(rels, r)
	}
	return rels
}
