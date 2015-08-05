/*
Purpose:	This is where we will handle the admin functions of Merx.
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
	"errors"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"merx/common"
	"merx/packages"
)

func handleAdmin(db *sql.DB, req *http.Request) (code int, data []byte, err error) {
	var pkg common.AdminPackage

	//09.25.2013 naj - get the path that the client passed.
	url := strings.ToLower(req.URL.Path[1:])

	if req.Method == "GET" {
		//10.01.2013 naj - get the query parameters
		urlquery := req.URL.Query()

		//10.01.2013 naj - now figure out which function we are calling.
		switch {
		case url == "getorders":
			//09.26.2013 naj - vendor is pulling pending orders
			code, data, err := packages.GetOrders(db, urlquery.Get("OrderDate"), urlquery.Get("AllOrders"))
			return code, data, err

		default:
			//log the error message to our log file
			logging.LogError("", "Unsupported request path "+url)
			code = http.StatusBadRequest
			err = errors.New("Unsupported request path " + url)
			return code, nil, err
		}
	} else {
		//10.01.2013 naj - now figure out which function we are calling
		switch {
		case url == "adddealer":
			//09.25.2013 naj - vendor is adding a dealer
			formdata := []byte(req.FormValue("data"))
			pkg = new(packages.AddDealers)

			//09.27.2013 naj - Load the data into the package
			code, err = pkg.ParsePackage(formdata, db)

			if err != nil {
				//log the error message to our log file
				logging.LogError("", err.Error())
				return code, nil, err
			}

			//09.27.2013 naj - Now process the package
			code, err := pkg.ProcessPackage()

			if err != nil {
				//log the error message to our log file
				logging.LogError("", err.Error())
				err = errors.New("Database Error")
				return code, nil, err
			}

			//10.03.2013 naj - return the results
			return code, nil, err

		case url == "deletedealer":
			//09.25.2013 naj - vendor is adding a dealer
			formdata := []byte(req.FormValue("Data"))
			pkg = new(packages.DeleteDealers)

			//09.27.2013 naj - Load the data into the package
			code, err = pkg.ParsePackage(formdata, db)

			if err != nil {
				logging.LogError("", err.Error())
				return code, nil, err
			}

			//09.27.2013 naj - Now process the package
			code, err := pkg.ProcessPackage()

			if err != nil {
				//log the error message to our log file
				logging.LogError("", err.Error())
				err = errors.New("Database Error")
				return code, nil, err
			}

			//10.03.2013 naj - return the results
			return code, nil, err

		case url == "updateorderstatus":
			//09.26.2013 naj - vendor is updating the order status.
			formdata := []byte(req.FormValue("data"))
			pkg = new(packages.POUpdates)

			//09.27.2013 naj - Load the data into the package
			code, err = pkg.ParsePackage(formdata, db)

			if err != nil {
				//log the error message to our log file
				logging.LogError("", err.Error())
				return code, nil, err
			}

			//09.27.2013 naj - Now process the package
			code, err := pkg.ProcessPackage()

			if err != nil {
				//log the error message to our log file
				logging.LogError("", err.Error())
				err = errors.New("Database Error")
				return code, nil, err
			}

			//10.03.2013 naj - return the results
			return code, nil, err

		default:
			//log the error message to our log file
			logging.LogError("", "Unsupported request path "+url)
			code = http.StatusBadRequest
			err = errors.New("Unsupported request path " + url)
			return code, nil, err
		}
	}

	return code, data, err
}
