/*
Purpose:	This is where we will dealer with the admin interface for Merx
Written:	09.26.2013
By:		Noel Jacques <njacques@nizex.com>
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

	_ "github.com/go-sql-driver/mysql"
)

func adminInterface(resp http.ResponseWriter, req *http.Request) {
	//09.25.2013 naj - log request.
	logging.LogError("", "Start of received admin request for "+req.URL.Path[1:])

	//09.25.2013 naj - open database connection
	db, err := sql.Open("mysql", DBConnect)

	if err != nil {
		//write the header out to the client system
		resp.WriteHeader(http.StatusInternalServerError)
		//write the error contents to the client system
		resp.Write([]byte("Database connection error"))
		//log the error message to our log file
		logging.LogError("", err.Error())
		return
	}
	defer db.Close()

	//09.25.2013 naj - confirm database connection
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

	//09.25.2013 naj - make sure that we have a GET or a POST.
	getpost := false

	if req.Method == "GET" {
		getpost = true
	}

	if req.Method == "POST" {
		getpost = true
	}

	if !getpost {
		//write the header out to the client system
		resp.WriteHeader(http.StatusBadRequest)
		//write the error contents to the client system
		resp.Write([]byte("Not a GET or POST request"))
		//log the error message to our log file
		logging.LogError("", "Not a GET or POST request")
		return
	}

	//10.01.2013 naj - call the handleAdmin function
	code, data, err := handleAdmin(db, req)

	if err != nil {
		//write the header out to the client system
		resp.WriteHeader(code)

		//write the error contents to the client system
		if code == http.StatusInternalServerError {
			resp.Write([]byte("Internal Server Error"))
		} else {
			resp.Write([]byte(err.Error()))
		}

		logging.LogError("", err.Error())
		return
	}

	//10.01.2013 naj - if we made it this far we have a successful request to write the response to the client
	resp.WriteHeader(code)
	resp.Write(data)
	logging.LogError("", "Admin Request Complete")
}
