# blueprint

[![Build Status](https://travis-ci.org/urld/blueprint.svg?branch=master)](https://travis-ci.org/urld/blueprint)
[![Go Report Card](https://goreportcard.com/badge/github.com/urld/blueprint)](https://goreportcard.com/report/github.com/urld/blueprint)
[![GoDoc](https://godoc.org/github.com/urld/blueprint?status.svg)](https://godoc.org/github.com/urld/blueprint)

`blueprint` is a tool to document and visualize your software architecture, based on the [C4 software architecture model](https://c4model.com).

## Installation

First install [graphviz](http://graphviz.org/Download.php) for your OS, then

	go get github.com/urld/blueprint/cmd/blueprint

## Usage
Run blueprint for a specific project directory to launch an interactive http server:

	blueprint test/ok

You can click on each element to change views.

![Example](https://github.com/urld/blueprint/blob/master/test/example.png)

### Syntax 

A project directory can contain multiple textfiles
containing an architecture description of the following syntax (currently similar to [structurizr express](https://structurizr.com/express)):

	Persona = Name | Description | Tags
	System = Name | Description | Tags
	Container = System Name | Name | Description | Technology | Tags
	Component = Container Name | Name | Description | Technology | Tags

	Relationship = Source Name | Description | Technology | Destination Name | Tags

	SystemContext = CoreSystems | ExternalSystems | Name | Description

`Tags`, `CoreSystems` and `ExternalSystems` accept comma separated lists of values.

An example can be found within `test/ok`.




## TODO

* html file export so there is no need for a http server running all the time
* improved layout
* more expressive syntax?
* more detailed usage documentation
