// SPDX-FileCopyrightText: Amolith <amolith@secluded.site>
//
// SPDX-License-Identifier: BSD-2-Clause

package main

import (
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// Loops through the list of members, checks whether they're up or down, whether
// they contain the requisite webring links, and appends any errors to the
// user-provided validation log
func (m *model) validateMembers() {
	// Get today's date with hours, minutes, and seconds
	today := time.Now().Format("2006-01-02")

	// Check the log header to see if we've already validated today
	logFile, err := os.Open(*flagValidationLog)
	if err != nil {
		fmt.Println("Error opening validation log:", err)
		logFile.Close()
		return
	}

	// Only read the most recent header, which is always 65 bytes long
	logHeader, err := io.ReadAll(io.LimitReader(logFile, 65))
	if err != nil {
		fmt.Println("Error reading validation log:", err)
		logFile.Close()
		return
	}

	if strings.Contains(string(logHeader), today) {
		logFile.Close()
		return
	}

	// Close the file so it's not locked while we're checking the members
	logFile.Close()

	// If any errors were found, write a report to the validation log
	errors := false

	// Start the report with a header
	report := "===== BEGIN VALIDATION REPORT FOR " + today + " =====\n"

	for _, r := range m.ring {
		errorMember := false
		reportMember := ""
		resp, err := http.Get("https://" + r.url)
		if err != nil {
			fmt.Println("Error checking", r.handle, "at", r.url, ":", err)
			reportMember += "  - Error with site: " + err.Error() + "\n"
			if !errors {
				errors = true
			}
			report += "- " + r.handle + " needs to fix the following issues on " + r.url + ":\n"
			report += reportMember
			continue
		}

		if resp.StatusCode != http.StatusOK {
			reportMember += "  - Site is not returning a 200 OK\n"
			if !errors {
				errors = true
			}
			report += "- " + r.handle + " needs to fix the following issues on " + r.url + ":\n"
			report += reportMember
			continue
		}

		// Read the response body into a string
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading webpage body:", err)
			continue
		}

		requiredLinks := []string{
			"https://" + *flagHost + "/next?host=" + url.QueryEscape(r.url),
			"https://" + *flagHost + "/previous?host=" + url.QueryEscape(r.url),
			"https://" + *flagHost,
		}

		decodedBody := html.UnescapeString(string(body));
		for _, link := range requiredLinks {
			if !strings.Contains(decodedBody, link) {
				reportMember += "  - Site is missing " + link + "\n"
				if err != nil {
					fmt.Println("Error writing to validation log:", err)
					continue
				}
				if !errors {
					errors = true
				}
				if !errorMember {
					errorMember = true
				}
			}
		}
		if errorMember {
			report += "- " + r.handle + " needs to fix the following issues on " + r.url + ":\n"
			report += reportMember
		}
	}

	report += "====== END VALIDATION REPORT FOR " + today + " ======\n\n"

	if errors {
		// Write the report to the beginning of the validation log
		f, err := os.OpenFile(*flagValidationLog, os.O_RDWR, 0o644)
		if err != nil {
			fmt.Println("Error opening validation log:", err)
			return
		}
		defer f.Close()

		logContents, err := io.ReadAll(f)
		if err != nil {
			fmt.Println("Error reading validation log:", err)
			return
		}

		if _, err := f.Seek(0, 0); err != nil {
			fmt.Println("Error seeking to beginning of validation log:", err)
			return
		}

		if _, err := f.Write([]byte(report)); err != nil {
			fmt.Println("Error writing to validation log:", err)
			return
		}

		if _, err := f.Write(logContents); err != nil {
			fmt.Println("Error writing to validation log:", err)
			return
		}
		fmt.Println("Validation report for " + today + " written")
	}
}
