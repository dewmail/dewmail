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
	"encoding/json"
	"fmt"
)

// An action structure for pushing to
type action struct {
	url  string
	body []byte
}

// Instantiates a new action object.
// Message used because we push all actions to the same domain as the to
// address in email. E.g., email to foo+add@example.com will push to
// http://example.com/foo/add.
func NewAction(m Message) (*action, error) {
	// Build an action object
	a := new(action)

	// Build url to push to
	a.url = "http"
	if OptToHTTPS {
		a.url += "s"
	}
	a.url += "://" + m.GetDomain() + m.GetPath()

	// Convert message object to JSON
	tempBody, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("Failed to encode message for send. %v", err)
	}
	a.body = tempBody

	return a, nil
}

// Initiates the HTTP request
func (a *action) Do() error {
	// Make call to API
	BuildJSONPost(a.url, a.body)

	// Push record to global datastore
	BuildJSONPost(OptDataStoreUrl, a.body)

	// Get mailsSent counter from global datastore
	sMailsSent, getCountErr := BuildJSONGet(OptDataStoreCountUrl)
	if getCountErr != nil {
		return fmt.Errorf("Error: Failed to get count of mails sent. %v", getCountErr)
	}

	// Get current count of mails sent
	var MailsSent map[string]int
	unmarshErr := json.Unmarshal([]byte(sMailsSent), &MailsSent)
	if unmarshErr != nil {
		return fmt.Errorf("Error: Failed to unmarshal mailsSent. %v", unmarshErr)
	}

	// Update mailsSent count in global datastore, try iMaxRetries times
	var updateCountErr error
	var sUpdateResp string
	var UpdateResp map[string]string
	const iMaxRetries int = 100
	for iRetries := 1; iRetries <= iMaxRetries; iRetries++ {
		// Send patch queries until we have incremented the count
		// The backend datastore will return an error json response if new value isn't greater than old
		sUpdateResp, updateCountErr = BuildJSONPatch(OptDataStoreCountUrl, []byte(fmt.Sprintf("{\"count\": %d}", MailsSent["count"]+iRetries)))
		unmarshErr := json.Unmarshal([]byte(sUpdateResp), &UpdateResp)
		if unmarshErr != nil {
			return fmt.Errorf("Error: Failed to unmarshal mails sent count update. %v", unmarshErr)
		}

		// Was a count response returned
		if _, success := UpdateResp["count"]; success {
			updateCountErr = nil
			break
		} else {
			updateCountErr = fmt.Errorf("Error: Failed to update sent mail count")
		}
	}

	return updateCountErr
}

