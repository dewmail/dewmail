/*
 * The MIT License (MIT)
 * 
 *  Copyright (c) 2014 Stephen Parker (withaspark.com)
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *  
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"encoding/json"
	"mime"
	"mime/multipart"
	"net/mail"
	"strings"
	"time"

	"code.google.com/p/go-smtpd/smtpd"
)

// A message structure
type Message struct {
	From     string `json:"from"`
	To       string `json:"to"`
	path     string
	domain   string
	Subject  string `json:"subject"`
	Body     string `json:"body"`
	Time     string `json:"time"`
	SPF      string `json:"spf"`
	IP       string `json:"sender-IP"`
	spath    string
	sdomain  string
	encoding string
	buffer   bytes.Buffer
}

// Getter of path
func (m *Message) GetPath() string {
	return m.spath
}

// Getter of domain
func (m *Message) GetDomain() string {
	return m.sdomain
}

// Add recipient method for Message types
// Required by package
func (m *Message) AddRecipient(to smtpd.MailAddress) error {
	m.To = strings.ToLower(to.Email())

	// Check if valid domain to receive mail for
	if OptDomainCheckingOn {
		_, dom := SplitToAddress(m.To)
		err := DomainCheck(OptValidDomains, dom)
		if err != nil {
			return fmt.Errorf("%v", err)
		}
	}

	return nil
}

// Add BeginData method for Message types
// Required by package
func (m *Message) BeginData() error {
	return nil
}

// Add parse method for Message types
func (m *Message) parse(r io.Reader) error {
	var err error
	message, _ := mail.ReadMessage(r)

	// Get headers
	m.To = message.Header.Get("To")
	m.encoding = message.Header.Get("Content-Type")
	var sReceived string = message.Header.Get("Received")

	// Verify sender before continuing
	if OptSPFCheck {
		var sSPFRes string
		sSPFRes, err = BuildJSONGet(OptSPFAPI, []byte(fmt.Sprintf(
			`{"apiKey": "%s", "email": "%s", "received": "%s"}`,
			OptSPFAPIKey, m.From, sReceived)))
		if err != nil {
			m.SPF = "TempError"
		} else {
			var SPFRes map[string]string
			err = json.Unmarshal([]byte(sSPFRes), &SPFRes)
			if err != nil {
				m.SPF = "TempError"
			} else {
				m.SPF = SPFRes["result"]
				m.IP = SPFRes["sender-IP"]
			}
		}
	} else {
		m.SPF = "None"
	}

	// If configured to require SPF pass only and this message doesn't pass, stop with error
	if OptRequireSPFPass && m.SPF != "Pass" {
		return fmt.Errorf("Error: SPF not pass for sender %s", m.From)
	}

	// Set time
	m.Time = time.Now().UTC().Format(time.RFC3339)

	// Parse to
	m.spath, m.sdomain = SplitToAddress(m.To)

	// Get subject
	m.Subject = message.Header.Get("Subject")

	// Get body
	tempbuf := new(bytes.Buffer)
	tempbuf.ReadFrom(message.Body)
	mediaType, args, _ := mime.ParseMediaType(m.encoding)
	if strings.HasPrefix(mediaType, "multipart") {
		part := multipart.NewReader(tempbuf, args["boundary"])
		for {
			nPart, err := part.NextPart()
			// If reached end of message, stop looping
			if err == io.EOF {
				return nil
			// Pass errors
			} else if err != nil {
				return err
			// Grab the text/plain version
			} else if strings.Contains(nPart.Header.Get("Content-Type"), "text/plain") {
				tempPart, _ := ioutil.ReadAll(nPart)
				m.Body = strings.Replace(strings.Replace(strings.TrimSpace(string(tempPart)), "\r\n", " ", -1), "\n", " ", -1)
			// Message had no text/plain formatting
			//TODO: One day add html parsing to strip tags
			} else {
				return errors.New("Error: No text/plain formatting of message.")
			}
		}
	}

	return nil
}

// Add Close method for Message types
// Required by package
func (m *Message) Close() error {
	// Do some processing here
	m.parse(&m.buffer)

	// Display what we received
	log.Printf("Received message from %s: [%s] %s\n", m.From, m.Subject, m.Body)

	// Send to application for processing
	a, err := NewAction(*m)
	if err != nil {
		log.Printf("Error: Failed to parse action. %v", err)
		return fmt.Errorf("Error: Failed to parse action. %v", err)
	}
	err = a.Do()
	if err != nil {
		log.Printf("Error: Failed to run action. %v", err)
		return fmt.Errorf("Error: Failed to run action. %v", err)
	}

	return nil
}

// Add Write method for Message types
// Required by package
func (m *Message) Write(line []byte) error {
	m.buffer.Write(line)
	return nil
}

// Extract project and domain from email address
func SplitToAddress(sTo string) (string, string) {
	var sMailbox, sDomain string

	// Separate address into mailbox@domain
	sToPieces := strings.Split(sTo, "@")
	sMailbox = sToPieces[0]
	sDomain = sToPieces[1]

	// If + in mailbox, use those to build path to API dest
	sMailbox = OptAPIRoute + strings.Replace(sMailbox, "+", "/", -1)

	return sMailbox, sDomain
}

// Determines if we should handle this domain
func DomainCheck(domains []string, domain string) error {
	var bMatch bool = false

	// Iterate over domains set in config
	for _, el := range domains {
		if strings.HasSuffix(domain, el) {
			bMatch = true
			break
		}
	}
	if !bMatch {
		return fmt.Errorf("Not accepting email for domain: %s", domain)
	}
	return nil
}
