/*
Purpose: This package handles retrieving vehicle specifications.
Written:	2014.06.25
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

package packages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

func GetVehicle(db *sql.DB, vin string) (code int, resp []byte, err error) {
	//06.25.2014 naj - first make sure we have a VIN.
	if vin == "" {
		code = http.StatusBadRequest
		resp = nil
		err = errors.New("Missing VIN")
		return
	}

	//06.25.2014 naj - now we need to lookup the vehicle in the database base.
	var r VehicleData

	//TODO query the database

	//TODO loop through the results

	resp, err = json.Marshal(r)
	if err != nil {
		code = http.StatusInternalServerError
		resp = nil
		err = errors.New("Error encoding JSON object")
		return
	}

	code = http.StatusOK
	return
}
