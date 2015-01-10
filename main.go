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
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"time"

	"code.google.com/p/go-smtpd/smtpd"
)

func main() {
	// Setup logfile
	tDate := time.Now()
	var sLogFile string = "logs/dewmail-" + tDate.Format("2006-01-02") + ".log"
	fpLog, err := os.OpenFile(sLogFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal("Error opening/creating logfile %v", err)
	}
	defer fpLog.Close()
	log.SetOutput(fpLog)

	// Start listener on HTTP port so we can use pinging services to verify service up
	http.HandleFunc("/", HandleHTTP)
	go http.ListenAndServe(":" + OptHTTPPort, nil)

	// Start listener on SMTP port
	var SMTPOptions = flag.String("smtp", ":25", "")
	SMTPListener, eSMTPError := net.Listen("tcp", *SMTPOptions)
	if eSMTPError != nil {
		log.Fatal(fmt.Errorf("Error listening for SMTP %v", eSMTPError))
	} else {
		fmt.Println("Listening for SMTP...")
	}

	// Start SMTP server
	SMTPServer := &smtpd.Server{
		OnNewConnection: func(conn smtpd.Connection) error {
			return nil
		},
		OnNewMail: func(conn smtpd.Connection, from smtpd.MailAddress) (smtpd.Envelope, error) {
			message := &Message{
				From: from.Email(),
			}
			return message, nil
		},
	}
	SMTPServer.Serve(SMTPListener)
}

// Serve a (malformed) HTML response
func HandleHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Dewmail</h1><p>Service is up</p>"))
}
