// SPDX-FileCopyrightText: Amolith <amolith@secluded.site>
//
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// Link returns an HTML, HTTPS link of a given URI
func link(l string) string {
	return "<a href='https://" + l + "'>" + l + "</a>"
}

// parseIndex parses the index template and returns a template struct.
func (m *model) parseIndex() {
	m.index = nil
	tmpl, err := template.ParseFiles(*flagIndex)
	if err != nil {
		log.Fatal(err)
	}
	m.index = tmpl
}

// List parses the list of members, appends the data to a slice of type list,
// then returns the slice
func (m *model) parseList() {
	file, err := ioutil.ReadFile(*flagMembers)
	if err != nil {
		log.Fatal("Error while loading list of webring members: ", err)
	}
	lines := strings.Split(string(file), "\n")
	m.ring = nil
	for _, line := range lines[:len(lines)-1] {
		fields := strings.Fields(line)
		m.ring = append(m.ring, ring{handle: fields[0], url: fields[1]})
	}
}

func is200(site string) bool {
	resp, err := http.Get(site)
	if err != nil {
		log.Println(err)
		return false
	}
	if resp.StatusCode == http.StatusOK {
		return true
	}
	return false
}
