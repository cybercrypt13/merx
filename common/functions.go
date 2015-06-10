/*
Purpose:	This handles getting the VendorID for a given bsvcode and VendorCode
Written:	07.09.2013
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

package common

import (
	"database/sql"
	"errors"

	_ "github.com/go-sql-driver/mysql"
)

//07.09.2013 naj - This function will retrieve a VendorID from the database based on bsvcode and vendorcode
func GetVendorID(db *sql.DB, bsvkeyid int, vendorcode string) (vendorid int, err error) {
	if vendorcode == "" {
		return 0, errors.New("No VendorCode found")
	}

	//07.09.2013 naj - check to see if we are using a bsv specific code or the standard code.
	if bsvkeyid > 0 {
		//07.09.2013 naj - BSV specific code.
		err = db.QueryRow("select VendorID from BSVVendorCodeLinks where BSVKeyID = ? and BSVVendorCode = ?", bsvkeyid, vendorcode).Scan(&vendorid)
	} else {
		err = db.QueryRow("select VendorID from VendorCodes where VendorCode = ?", vendorcode).Scan(&vendorid)
	}

	switch {
	case err == sql.ErrNoRows:
		err = errors.New("Unable to locate Vendor")
		return
	case err != nil:
		return
	default:
		return
	}

}

//07.12.2013 naj - This function will lookup a part and return details about the part.
func GetPart(db *sql.DB, vendorid int, partnumber string) (itemid int, superseded string, nla bool, err error) {
	var super int
	err = db.QueryRow("select ItemID, SupersessionID, NLA from Items where VendorID = ? and "+
		"PartNumber = ?", vendorid, partnumber).Scan(&itemid, &super, &nla)
	switch {
	case err == sql.ErrNoRows:
		return 0, "", false, errors.New("Unable to locate part number")
	case err != nil:
		return 0, "", false, err
	default:
		if super > 0 {
			var tempid int
			var tempnla bool
			err = db.QueryRow("select ItemID, PartNumber, NLA from Items where ItemID = ?", super).Scan(&tempid, &superseded, &tempnla)
			switch {
			case err == sql.ErrNoRows:
				//07.12.2013 naj - could not locate the superseded part so we will return treat it as if it was not superseded.
				superseded = ""
				err = nil
			case err != nil:
				return 0, "", false, err
			default:
				//07.12.2013 naj - we found the supersession so we set our itemid and nla variable to match the superseded part
				itemid = tempid
				nla = tempnla
			}
		}
	}

	return
}

//07.15.2014 naj - This function will retrieve the DealerKey based on the dealers ID.
//If a dealer has successfully authenticated with the server we always return the DealerKey to the client.
func GetDealerKey(db *sql.DB, dealerid int) (dealerkey string, err error) {
	dealerkey = ""

	if dealerid < 1 {
		err = errors.New("Invalid Dealer Key")
		return
	}

	err = db.QueryRow("select DealerKey from DealerCredentials where DealerID = ?", dealerid).Scan(&dealerkey)
	switch {
	case err == sql.ErrNoRows:
		err = errors.New("Unable to locate DealerKey")
		return
	case err != nil:
		return
	default:
		return
	}
}
