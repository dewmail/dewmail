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
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// Builds a json encoded HTTP PATCH request
func BuildJSONPatch(url string, content []byte) (string, error) {
	return BuildJSONRequest("PATCH", url, content)
}

// Builds a json encoded HTTP POST request
func BuildJSONPost(url string, content []byte) (string, error) {
	return BuildJSONRequest("POST", url, content)
}

// Builds an HTTP GET request
func BuildJSONGet(url string) (string, error) {
	return BuildJSONRequest("GET", url, nil)
}

// Builds a generic HTTP request, handles errors and logging
func BuildJSONRequest(reqType string, url string, content []byte) (string, error) {
	// Check url valid
	//TODO: Add better url test
	if len(url) < 1 {
		return "", fmt.Errorf("Invalid URL %s", url)
	}

	// Make request to url containing JSON content
	request, err := http.NewRequest(reqType, url, bytes.NewBuffer(content))
	if err != nil {
		return "", fmt.Errorf("Failed to build request. %v", err)
	}

	// Display/log request for debugging
	log.Printf("Request: [%s] %s", url, content)

	// Set request headers
	request.Header.Set("Content-Type", "application/json")

	// Initiate HTTP client
	clientHandler := &http.Client{
		Transport: &http.Transport{
			// Allow certs with name mismatch
			//FIXME: This needs to be addressed, susceptible to
			//       MITM attacks. Bug 154.
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	response, err := clientHandler.Do(request)

	// Check for errors on response
	if err != nil {
		return "", fmt.Errorf("Failed to get response from %s. %v", url, err)
	}
	defer response.Body.Close()

	// Display/log response
	body, _ := ioutil.ReadAll(response.Body)
	log.Printf("Response: [%s] %s %s", response.Status, response.Header, string(body))

	return string(body), nil
}

