Dewmail
=======

An easily deployable SMTP server system for integrating with existing API using convenient hooks. Written in Go. Receive email at your project's email address and automatically initiate a JSON formatted POST request to your app's API.

## Use ##
Dewmail will receive emails on behalf of your domain and initiate a JSON encoded POST request to your API. Let's look at a specific example to cover some of the implementation details.

If you create the following MX record for your domain ```example.com```

		TYPE|SUBDOMAIN|MAILSERVER HOST|PREF|TTL
		----|---------|---------------|----|---
		MX|do|dewmail.withaspark.com.|10|300

You may then send an email from any client--for example, the following message.

To: ```foo+bar@do.example.com```
Subject: ```mail```
Body: ```unsubscribe```
		
Dewmail will receive and parse your email generating the following JSON request ```{"from":"client@example.org","to":"foo+bar@do.example.com","subject":"mail","body":"unsubscribe"}```, which it will submit to ```http://do.example.com/foo/bar```. Your API can then parse the message as you see fit.

## Demo ##

For a demo, send an email to [test@dewmail.withaspark.com](mailto:test@dewmail.withaspark.com) and go to [http://dewmail.withaspark.com/](http://dewmail.withaspark.com/).

## Install ##
1. Add the following MX record to the domain you wish to receive calls from Dewmail.

		TYPE|SUBDOMAIN|MAILSERVER HOST|PREF|TTL
		----|---------|---------------|----|---
		MX|somesubdomain|dewmail.withaspark.com.|10|300

	**Note:** this will be the domain of the email address you must receive messages at and also the domain Dewmail will POST to looking for your API.
2. Pull down Dewmail binary or build from source.
3. Optional, add site to webserver if you want a way to query app status via HTTP. A sample Apache site .conf file is provided as an example [dewmail.apache.conf](dewmail.apache.conf), which will listen for requests to ```/up``` and will display a message that the service is up.
4. Make sure port 25 isn't blocked by your firewall.
5. Start Dewmail.
6. Send Dewmail your email.

## Configuration ##

The previous example describes the default setup. There are a few options that may be configured for your particular setup. The following options can be set in the master configuration file [config.go](config.go).

Option|Type|Function
------|----|--------
```OptDomainCheckingOn```|```bool```|If ```true```, will enforce only accepting mail for specific domains in ```OptValidDomains```. Default, off, ```false```.
```OptValidDomains```|```[]string```|Domain(s) to accept emails for. By default, not used.
```OptAPIRoute```|```string```|Route for Dewmail to post to (with leading and trailing slashes). Default, ```/```.
```OptToHTTPS```|```bool```|Whether to use HTTPS/HTTP for post request. Default, HTTPS, ```true```.

