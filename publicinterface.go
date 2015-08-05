/*
Purpose:	This file receives the incoming request and determines what should
			be done with it.  This is the main file called by the merx.go 
			processes to get things kicked off.

Written:	05.28.2013
By:		Glenn Hancock <ghancock@nizex.com>
URL:		www.nizex.com

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
	"database/sql"
	"net/http"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"merx/common"
)

//05.06.2015 ghh - this is the main function that is called whenever
//a new HTTP/S request is received.  Its main goal is to figure out
//what type of request has been received and then to call the necessary
//routines to have it handled
func publicInterface(resp http.ResponseWriter, req *http.Request) {
	//05.28.2013 naj - log request.
	logging.LogError("", "Start of received data for "+string(req.URL.Path))

	//04.11.2012 ghh - now we're going to create a db handle that will be used
	//for the remainder of the functions so that they can all use a single db
	//connection
	db, err := sql.Open("mysql", DBConnect)

	//05.06.2015 ghh - if we can't get a db connection then log error
	//and return as there is no point in continuing.
	if err != nil {
		//write the header out to the client system
		resp.WriteHeader(http.StatusInternalServerError)
		//write the error contents to the client system
		resp.Write([]byte("Database connection error"))
		//log the error message to our log file
		logging.LogError("", err.Error())
		return
	}

	//05.06.2015 ghh - make sure when this function has completed and returned
	//that we close up our database connections.
	defer db.Close()

	//05.28.2013 naj - confirm database connection
	err = db.Ping()

	if err != nil {
		//write the header out to the client system
		resp.WriteHeader(http.StatusInternalServerError)
		//write the error contents to the client system
		resp.Write([]byte("Database Connection Error"))
		//log the error message to our log file
		logging.LogError("", err.Error())
		return
	}

	//05.29.2013 naj - make sure that we have a GET or a POST.
	if req.Method != "GET" && req.Method != "POST" {
		//write the header out to the client system
		resp.WriteHeader(http.StatusBadRequest)
		//write the error contents to the client system
		resp.Write([]byte("Not a GET or POST request"))
		//log the error message to our log file
		logging.LogError("", "Not a GET or POST request")
		return
	}

	//05.29.2013 naj - authenticate the client
	code, dealerkey, dealerid, bsvkeyid, err := common.AuthenticateClient(db, req)

	if err != nil {
		//write the header out to the client system
		resp.WriteHeader(code)
		//write the error contents to the client system
		if code == http.StatusInternalServerError {
			resp.Write([]byte("500 Internal Server Error\n" + string(code)))
		} else {
			//05.29.2013 naj - make the connection wait for 5 seconds when an authentication failure occurs.
			time.Sleep(5000 * time.Millisecond)
			resp.Write([]byte("401 Unauthorized\n" + err.Error()))
		}
		//log the error message to our log file
		logging.LogError("", err.Error())
		return
	}


	//06.04.2013 naj - now call the handlePackage function which will decide what to do
	//based on the request method and url path
	code, json, err := handlePackage(req, bsvkeyid, dealerid, dealerkey, db )

	//05.31.2013 naj - see if we got back an error
	if err != nil {
		//write the header out to the client system
		resp.WriteHeader(code)
		//write the error contents to the client system
		if code == http.StatusInternalServerError {
			resp.Write([]byte("500 Internal Server Error\n"))
		} else {
			resp.Write([]byte(err.Error() + "\n"))
		}
		//log the error message to our log file
		logging.LogError("", err.Error() + dealerkey )
	} else {
		//06.03.2013 naj - if we got here then we have a successful request
		//we will not write the header and json response text back to the client.
		//write the header out to the client system
		resp.WriteHeader(code)
		//write the error contents to the client system
		resp.Write(json)
		resp.Write([]byte("\n"))
		//log the error message to our log file
		logging.LogError("", "Request Completed")
	}
	return
}
