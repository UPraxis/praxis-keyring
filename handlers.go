// SPDX-FileCopyrightText: Amolith <amolith@secluded.site>
//
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"html/template"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// Serves the webpage created by createRoot()
func (m model) root(writer http.ResponseWriter, request *http.Request) {
	var table string
	for _, member := range m.ring {
		table = table + "  <tr>\n"
		table = table + "    <td>" + member.handle + "</td>\n"
		table = table + "    <td>" + link(member.url) + "</td>\n"
		table = table + "  </tr>\n"
	}
	err := m.index.Execute(writer, template.HTML(table))
	if err != nil {
		log.Println("Error executing template: " + err.Error())
		http.Error(writer, "Internal server error", 500)
	}
}

// Redirects the visitor to the next member, wrapping around the list if the
// next would be out-of-bounds, and ensuring the destination returns a 200 OK
// status before performing the redirect.
func (m model) next(writer http.ResponseWriter, request *http.Request) {
	host := request.URL.Query().Get("host")
	scheme, success := "https://", false
	length := len(m.ring)
	for i, item := range m.ring {
		if item.url == host {
			for j := i + 1; j < length+i; j++ {
				dest := scheme + m.ring[j%length].url
				log.Println("Checking '" + dest + "'")
				if is200(dest) {
					log.Println("Redirecting visitor to '" + dest + "'")
					http.Redirect(writer, request, dest, http.StatusFound)
					success = true
					break
				}
				log.Println("Something went wrong accessing '" + dest + "', skipping site")
			}
			http.Error(writer, `It would appear that either none of the ring members are accessible
(unlikely) or the backend is broken (more likely). In either case,
please `+*flagContactString, 500)
		}
	}
	if !success {
		http.Error(writer, "Ring member '"+host+"' not found.", http.StatusNotFound)
	}
}

// Redirects the visitor to the previous member, wrapping around the list if the
// next would be out-of-bounds, and ensuring the destination returns a 200 OK
// status before performing the redirect.
func (m model) previous(writer http.ResponseWriter, request *http.Request) {
	host := request.URL.Query().Get("host")
	scheme := "https://"
	length := len(m.ring)
	for index, item := range m.ring {
		if item.url == host {
			// from here to start of list
			for i := index - 1; i > 0; i-- {
				dest := scheme + m.ring[i].url
				if is200(dest) {
					log.Println("Redirecting visitor to '" + dest + "'")
					http.Redirect(writer, request, dest, http.StatusFound)
					return
				}
			}
			// from end of list to here
			for i := length - 1; i > index; i-- {
				dest := scheme + m.ring[i].url
				if is200(dest) {
					log.Println("Redirecting visitor to '" + dest + "'")
					http.Redirect(writer, request, dest, http.StatusFound)
					return
				}
			}
			http.Error(writer, `It would appear that either none of the ring members are accessible
(unlikely) or the backend is broken (more likely). In either case,
please `+*flagContactString, 500)
			return
		}
	}
	http.Error(writer, "Ring member '"+host+"' not found.", http.StatusNotFound)
}

// Redirects the visitor to a random member
func (m model) random(writer http.ResponseWriter, request *http.Request) {
	rand.Seed(time.Now().Unix())
	dest := "https://" + m.ring[rand.Intn(len(m.ring)-1)].url
	http.Redirect(writer, request, dest, http.StatusFound)
}

// Serves the log at *flagValidationLog
func (m model) validationLog(writer http.ResponseWriter, request *http.Request) {
	http.Header.Add(writer.Header(), "Content-Type", "text/plain")
	http.ServeFile(writer, request, *flagValidationLog)
}
