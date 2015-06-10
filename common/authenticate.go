/*
Purpose:	This handles authenticating the client
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

package common

import (
	"database/sql"
	"errors"
	"net/http"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/satori/go.uuid"
)

//05.29.2013 naj - this function will authenticate the client
func AuthenticateClient(db *sql.DB, 
		req *http.Request) (code int, dealerkey string, 
		dealerid int, bsvkeyid int, err error) {
	//06.03.2013 naj - initialize some variables
	var accountnumber, sentdealerkey, bsvkey string
	code = http.StatusOK

	//05.29.2013 naj - first we grab the AccountNumber and DealerKey
	if req.Method == "GET" {
		//first we need to grab the query string from the url so
		//that we can retrieve our variables
		temp := req.URL.Query()
		accountnumber = temp.Get("accountnumber")
		sentdealerkey = temp.Get("dealerkey")
		bsvkey = temp.Get("bsvkey")
	} else {
		accountnumber = req.FormValue("accountnumber")
		sentdealerkey = req.FormValue("dealerkey")
		bsvkey = req.FormValue("bsvkey")
	}


	//if we don't get back a BSV key then we need to bail as
	//its a requirement. 
	if bsvkey == "" {
		err = errors.New("Missing BSV Key In Package")
		code = http.StatusUnauthorized
		return
	}

	//if we didn't get an account number for the customer then we need to
	//also bail
	if accountnumber == "" {
		err = errors.New("Missing account number")
		code = http.StatusUnauthorized
		return
	}

	//06.03.2013 naj - validate the BSVKey to make sure the the BSV has been certified for MerX
	err = db.QueryRow(`select BSVKeyID from AuthorizedBSVKeys 
							where BSVKey = ?`, bsvkey).Scan(&bsvkeyid)

	switch {
		case err == sql.ErrNoRows:
			//05.06.2015 ghh - here, instead of sending back an invalid key just yet we need
			//to first query the merx keyvalidator service at beta.nizex.com to see if the
			//key has in fact been registered and if so we'll just add it to our local tables
			//for future use.
			err = errors.New("Invalid BSV Key")
			code = http.StatusUnauthorized
			return
		case err != nil:
			code = http.StatusInternalServerError
			return
		}

	//05.29.2013 naj - check to see if the supplied credentials are correct.
	//06.24.2014 naj - new format of request allows for the dealer to submit a request without a dealerkey on the first request to merX.
	err = db.QueryRow(`select DealerID, ifnull(DealerKey, '') as DealerKey 
							from DealerCredentials where AccountNumber = ? 
							and Active = 1 `, 
							accountnumber).Scan(&dealerid, &dealerkey )

	switch {
		case err == sql.ErrNoRows:
			err = errors.New("Account not found")
			code = http.StatusUnauthorized
			return
		case err != nil:
			code = http.StatusInternalServerError
			return
	}

	//05.06.2015 ghh - now we check to see if we have a valid key for the dealer
	//already.  If they don't match then we get out. Keep in mind they could send
	//a blank key on the second attempt after we've generated a key and we need
	//to not allow that.
	if sentdealerkey != dealerkey {
		err = errors.New("Access Key Is Not Valid" )
		code = http.StatusUnauthorized
		return
	}

	//06.03.2013 naj - parse the RemoteAddr and update the client credentials
	address := strings.Split(req.RemoteAddr, ":")

	//06.24.2014 naj - If we got this far then we have a dealerid, now we need to see if 
	//they dealerkey is empty, if so create a new key and update the dealer record.
	if dealerkey == "" {
		dealerkey = uuid.NewV1().String()

		_, err = db.Exec(`update DealerCredentials set DealerKey = ?,
								LastIPAddress = inet_aton(?),
								AccessedDateTime = now()
								where DealerID = ?`, dealerkey, address[0], dealerid)

		if err != nil {
			code = http.StatusInternalServerError
			return
		}
	}

	_, err = db.Exec(`update DealerCredentials set LastIPAddress = inet_aton(?), 
						AccessedDateTime = now() 
						where DealerID = ?`, address[0], dealerid)

	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	return
}
