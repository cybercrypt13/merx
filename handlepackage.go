/*
Purpose:	Merx is a service that will listen for requests from dealers
			relay them to suppliers and manufacturers.
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
	"errors"
	"net/http"
	"strings"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"merx/common"
	"merx/packages"
)

//04.09.2012 ghh - this function is threaded and receives each
//incoming request from the webserver.  Technically this is the
//main entry point for the program that gets called from main
//this function takes the url and passes back the required structure object
//that we're going to work with in the main thread
func handlePackage(	req *http.Request, 
							bsvkeyid int, 
							dealerid int, 
							dealerkey string, 
							db *sql.DB) (code int, json []byte, err error) {
	url := req.URL.Path[1:]
	code = http.StatusOK

	//logging.LogError("", "Inside Handle Package function" )

	//05.28.2013 naj - here we need to handle the request type
	//Merx only support GET and POST, so that is all we need to check for.
	if req.Method == "GET" {
		urlquery := req.URL.Query()

		//07.12.2013 ghh - get purchase order status
		if strings.ToLower(url) == "getorderstatus" {
			code, json, err = packages.GetOrderStatus(dealerid, urlquery.Get("InternalID"), db)
			if err != nil {
				return
			}
		}

		//convert vendorid to int

		//07.12.2013 ghh - check inventory
		if strings.ToLower(url) == "inventoryverify" {
			vendorid, _ := strconv.ParseInt(urlquery.Get( "VendorID" ),10,0)

			code, json, err = packages.GetInventoryLocal(dealerid, 
							bsvkeyid, 
							int(vendorid), 
							urlquery.Get("PartNumber"), 
							db)
			if err != nil {
				return
			}
		}
	} else {
		//06.03.2013 naj - This is a POST request which means we will load data
		//from the client into our database.
		//take the received data from req object and store it in a byte array
		//in order to pass it to our Json control later
		//now we create a varable of our interface structure so that we can use
		//it to pass through to the proper place.  This object is in
		//common/interface.go
		var pkg common.Packagecontents
		formdata := []byte(req.FormValue("data"))

		//05.31.2013 naj - figure out which package to load.
		if strings.ToLower(url) == "sendorder" {
			pkg = new(packages.POSend)


		}

		//05.31.2013 naj - make sure we have a package
		if pkg == nil {
			return http.StatusNotFound, nil, errors.New("404 Not Found")
		}

		//now that we have  our interface, we need to next parse the received
		//package contents and load it into our structure that has been chosen
		err = pkg.ParsePackage(formdata, db, bsvkeyid)
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		//06.03.2013 naj - now verify the package contains all the required data
		err = pkg.VerifyPackage()
		if err != nil {
			return http.StatusBadRequest, nil, err
		}

		//06.03.2013 naj - now load the package into the database
		json, err = pkg.ProcessPackage(dealerid, dealerkey)
		if err != nil {
			switch {
			case err.Error() == "No valid parts were in the purchase order":
				return http.StatusBadRequest, nil, err
			case err.Error() == "PO is already being processed":
				return http.StatusNotAcceptable, nil, err
			default:
				return http.StatusInternalServerError, nil, err
			}
		}
	}
	return
}
