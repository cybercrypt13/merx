/*
Purpose:	This file contains the code that sets up the service and listens for
			incoming requests coming from customers.  it is also responsible for
			setting up the database connection and logfile processes

Written:	05.28.2013
By:		Glenn Hancock <ghancock@nizex.com>
URL:		www.nizex.com
Version:	1.0.0.1

The MIT License (MIT)

Copyright (c) 2013 Nizex Inc.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
*/

package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
	common "merx/common"
)

//05.28.2013 naj - declare our global variables for logging and database connection info.
var logging *common.LogFile
var DBConnect string

func main() {
	servicename := "Merx"

	//06.17.2014 naj - handle command line options
	confpath := flag.String("conf", "/etc/merx/merx.conf", "Path of configuration file")
	flag.Parse()

	//05.24.2013 naj - read the config file
	err := common.OpenConfigFile(*confpath)

	if err != nil {
		fmt.Println(err)
		return
	}

	//05.24.2013 naj - get the the log file path
	logfile := common.GetConfig("Defaults", "LogFile", "/var/log/merx.log")

	//05.24.2013 naj - start the logging.
	logging, err = common.NewLog(logfile)

	if err != nil {
		fmt.Println(err)
		return
	}

	//05.24.2013 naj - get the network settings
	netip := common.GetConfig("Network", "Listen", "0.0.0.0")
	netport := common.GetConfig("Network", "Port", "8000")
	sslcert := common.GetConfig("Network", "SSLCert", "")
	sslkey := common.GetConfig("Network", "SSLKey", "")

	//09.25.2013 naj - get the admin interface settings
	adminip := common.GetConfig("Network", "AdminListen", netip)
	adminport := common.GetConfig("Network", "AdminPort", "9000")

	//09.25.2013 naj - make sure that the IP address and port number
	//of the admin interface is not the same as the regular request interface.
	if netip == adminip && netport == adminport {
		fmt.Println("Admin Interface cannot be on the same IP address and port")
		return
	}

	//05.24.2013 naj - get the database connection info
	dbName := common.GetConfig("Database", "Name", "merx")
	dbHost := common.GetConfig("Database", "Host", "localhost")
	dbPort := common.GetConfig("Database", "Port", "3306")
	dbUser := common.GetConfig("Database", "User", "")
	dbPass := common.GetConfig("Database", "Password", "")
	dbProto := common.GetConfig("Database", "Protocol", "tcp")

	//05.24.2013 naj - make sure we have a database user and password
	if dbUser == "" || dbPass == "" {
		fmt.Println("Database user or password is missing!\nPlease check the config file.")
		return
	}

	//05.28.2013 naj - Now build the database connection string.
	DBConnect = dbUser + ":" + dbPass + "@"

	if dbProto == "unix" {
		dbSocket := common.GetConfig("Database", "Socket", "/tmp/mysql.sock")
		DBConnect = DBConnect + "(" + dbSocket + ")/" + dbName + "?charset=utf8&autocommit=true"
	} else {
		DBConnect = DBConnect + "(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8&autocommit=true"
	}

	//05.28.2013 naj - start the signal catcher.
	go common.SignalCatcher(servicename, logging)

	//05.28.2013 naj - start the http server
	logging.LogError("", "Starting Merx Server")
	fmt.Println("Merx Server Listening on " + netip + ":" + netport)

	//05.06.2015 ghh - setup the main listener that will catch all of the incoming requests
	// and then call the publicinterface routines to determine what needs to be done with
	//them
	public := &http.Server{
		Addr:           netip + ":" + netport,
		Handler:        http.HandlerFunc(publicInterface),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//05.28.2013 naj - see if we are using ssl or not.
	if sslcert == "" && sslkey == "" {
		go public.ListenAndServe()
	} else {
		go public.ListenAndServeTLS(sslcert, sslkey)
	}

	//09.25.2013 naj - start the admin interface
	fmt.Println("Merx Server Admin Listening on " + adminip + ":" + adminport)

	//05.06.2015 ghh - set up the main admin listener that allows someone to do
	//administration level tasks for maintenance of the server
	admin := &http.Server{
		Addr:           adminip + ":" + adminport,
		Handler:        http.HandlerFunc(adminInterface),
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	//05.28.2013 naj - see if we are using ssl or not.
	if sslcert == "" && sslkey == "" {
		err = admin.ListenAndServe()
	} else {
		err = admin.ListenAndServeTLS(sslcert, sslkey)
	}

	//05.28.2013 naj - check to see if we got an error from ListenAndServe
	if err != nil {
		logging.LogError("", "Error starting Merx Server Admin: "+err.Error())
		fmt.Println("Error starting Merx Server Admin: " + err.Error())
		return
	}
}
