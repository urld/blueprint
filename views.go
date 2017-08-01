// Copyright (c) 2017, David Url
// Use of this source code is governed by the
// GNU General Public License Version 2
// which can be found in the LICENSE file.

package main

import "io"

type View interface {
	Title() string
	Description() string
	Dot(w io.Writer, model Model) error
}

type SystemContextView struct {
	title       string
	description string
	Systems     []string
	Personas    []string
}

func (v SystemContextView) Description() string {
	return v.description
}

func (v SystemContextView) Title() string {
	return v.title
}

type ContainerView struct {
	Title       string
	Description string
	System      string
	Containers  []string
	Systems     []string
	Personas    []string
}

type ComponentView struct {
	Title       string
	Description string
	Container   string
	Components  []string
	Containers  []string
	Systems     []string
	Personas    []string
}
