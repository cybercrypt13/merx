/*
Purpose: This package handles deleting dealers from Merx.
Written:	09.27.2013
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
	"strconv"

	_ "github.com/go-sql-driver/mysql"
)

//04.10.2012 ghh - this function is reponsible for parsing our Json package
//and placing into internal array for us to work with
func (p *DeleteDealers) ParsePackage(pkg []byte, db *sql.DB) (code int, err error) {
	err = json.Unmarshal(pkg, &p.Dealers)

	if err != nil {
		code = http.StatusBadRequest
		err = errors.New("Error parsing JSON data")
		return code, err
	}

	//09.26.2013 naj - validate the data.
	var errtext string
	errtext = ""

	if len(p.Dealers) > 0 {
		for i := 0; i < len(p.Dealers); i++ {
			c := p.Dealers[i]

			if c.AccountNumber == "" {
				errtext += "Record " + strconv.FormatInt(int64(i), 10) + " Missing Account Number\n"
			}
		}
	} else {
		errtext = "No Dealer Records Submitted\n"
	}

	if errtext != "" {
		code = http.StatusBadRequest
		err = errors.New(errtext)
		return code, err
	}

	p.db = db
	return 0, nil
}

//this function actually takes the contents of our structure and processes it
func (p *DeleteDealers) ProcessPackage() (code int, err error) {
	db := p.db

	for i := 0; i < len(p.Dealers); i++ {
		c := p.Dealers[i]

		_, err := db.Exec("update DealerCredentials set Active = 0 where AccountNumber = ?", c.AccountNumber)

		if err != nil {
			return http.StatusInternalServerError, err
		}
	}
	return http.StatusOK, nil
}
