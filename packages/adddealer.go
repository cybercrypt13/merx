/*
Purpose: This package handles adding dealers to Merx.
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

package packages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
)

//04.10.2012 ghh - this function is reponsible for parsing our Json package
//and placing into internal array for us to work with
func (p *AddDealers) ParsePackage(pkg []byte, db *sql.DB) (code int, err error) {
	err = json.Unmarshal(pkg, &p.Dealers)

	if err != nil {
		code = http.StatusBadRequest
		err = errors.New("Error parsing JSON data")
		return code, err
	}

	//09.26.2013 naj - validate the data.
	var errtext string
	errtext = ""

	for i := 0; i < len(p.Dealers); i++ {
		c := p.Dealers[i]

		if c.AccountNumber == "" {
			errtext += " Missing Account Number\n"
		}

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
func (p *AddDealers) ProcessPackage() (code int, err error) {
	db := p.db

	for i := 0; i < len(p.Dealers); i++ {
		c := p.Dealers[i]

		if c.IPAddress == "" {
			c.IPAddress = "NULL"
		}

		//10.03.2013 naj - first see if the dealer has already been added.
		var count int
		err := db.QueryRow("select count(*) from DealerCredentials where Active = 1 and AccountNumber = ?", c.AccountNumber).Scan(&count)

		switch {
		case err == sql.ErrNoRows:
			count = 0
		case err != nil:
			code = http.StatusInternalServerError
			return code, err
		}

		if count == 0 {
			_, err = db.Exec("insert into DealerCredentials (AccountNumber, IPAddress, Active, CreatedDateTime) "+
				"values (?, inet_aton(?), 1, now())", c.AccountNumber, c.IPAddress)

			if err != nil {
				code = http.StatusInternalServerError
				return code, err
			}
		}
	}
	code = http.StatusCreated
	return code, nil
}
