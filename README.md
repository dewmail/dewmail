Dewmail
=======

An open-source email parsing microservice for HTTP APIs. Written in Go. Receive email at your project's email address and automatically initiate a JSON formatted POST request to your app's existing API.

## Use ##
Dewmail will receive emails on behalf of your domain and initiate a JSON encoded POST request to your API. Let's look at a specific example to cover some of the implementation details.

First, create the following MX records for your domain ```example.com```

		TYPE|SUBDOMAIN|MAILSERVER HOST|PREF|TTL
		----|---------|---------------|----|---
		MX|api|in1.dewmail.io.|10|300
		MX|api|in2.dewmail.io.|20|300

You may then send an email from any client--for example, consider the following message.

```
To: foo+bar@api.example.com
Subject: mail
Body: unsubscribe
```
		
Dewmail will receive and parse your email generating the following JSON request
 ```{"body":"unsubscribe","from":"client@example.org","sender-IP":"12.34.56.78","spf":"Pass","subject":"mail","time":"2015-01-17T08:00:39Z","to":"foo+bar@api.example.com"}```, which it will submit to ```http://api.example.com/foo/bar```. Your API can then parse the message however you see fit.

## Demo ##

For a demo, send an email to [test@demo.dewmail.io](mailto:test@demo.dewmail.io) and go to [http://dewmail.io/demo.php](http://dewmail.io/demo.php).

## Install ##
1. Add the following MX record to the domain you wish to receive calls from Dewmail.

		TYPE|SUBDOMAIN|MAILSERVER HOST|PREF|TTL
		----|---------|---------------|----|---
		MX|somesubdomain|in1.dewmail.io.|10|300
		MX|somesubdomain|in2.dewmail.io.|20|300

	**Note:** this will be the domain of the email address you must receive messages at and also the domain Dewmail will POST to looking for your API.
2. Pull down Dewmail binary or build from source.
3. Optional, add site to webserver if you want a way to query app status via HTTP. A sample Apache [.htaccess](api/public_html/.htaccess) file is provided to proxy requests for ```/up``` to the Dewmail app and will display a message if the service is up.
4. Make sure port 25 isn't blocked by your firewall.
5. Start Dewmail.
6. Start sending emails.

## Configuration ##

The previous example describes the default setup. There are a few options that may be configured for your particular setup. The following options can be set in the master configuration file [config.go](config.go).

Option|Type|Function
------|----|--------
```OptDomainCheckingOn```|```bool```|If ```true```, will enforce only accepting mail for specific domains in ```OptValidDomains```. Default, off, ```false```.
```OptValidDomains```|```[]string```|Domain(s) to accept emails for. By default, not used.
```OptAPIRoute```|```string```|Route for Dewmail to POST to (with leading and trailing slashes). Default, ```/```.
```OptToHTTPS```|```bool```|Whether to use HTTPS/HTTP for POST request. Default, HTTP, ```false```.
```OptDataStoreUrl```|```string```|URL to backend datastore. All messages for all domains will be POSTed here as well. Default, none, ```""```.
```OptDataStoreCountUrl```|```string```|URL to backend datastore message sent count. All messages for all domains will increment this. Default, none, ```""```.
```OptHTTPPort```|```string```|HTTP port number for serving web requests. Default, ```8111```.
```OptSPFCheck```|```bool```|Whether we should check SPF validation. Default, ```false```.
```OptRequireSPFPass```|```bool```|Whether we should abort on SPF anything but pass. Default, ```false```.
```OptSPFAPI```|```string```|Domain of API to check SPF results from. Default, ```""```.
```OptSPFAPIKey```|```string```|API key for SPF validation system. Default, ```""```.
