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

import ()

// Domain validation checks.
// If true, will enforce only accepting mail for specific domains in
// OptValidDomains. I am neglecting this because if you control the nameservers,
// I'll assume you're the boss and can handle API calls.
const OptDomainCheckingOn bool = false

// Domain(s) to accept emails for
var OptValidDomains []string = []string{
	"example.com",
	"do.example.com",
	"api.example.org",
}

// Route for Dewmail to POST to (with leading and trailing slashes)
const OptAPIRoute string = "/"

// Use HTTPS for POST request to API
const OptToHTTPS bool = false

// URL to backend datastore. All messages for all domains will be POSTed here as well.
const OptDataStoreUrl string = ""

// URL to backend datastore message sent count. All messages for all domains will increment this.
const OptDataStoreCountUrl string = ""

// HTTP port number for serving web requests
const OptHTTPPort string = "8111"

// Whether we should check SPF validation
const OptSPFCheck bool = false

// Whether we should abort on SPF anything but pass
const OptRequireSPFPass bool = false

// Domain of API to check SPF results from
const OptSPFAPI string = ""

// API key for SPF validation system
const OptSPFAPIKey string = ""
