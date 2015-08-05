/*
Purpose: 	This package handles invenory verification requests
Written:	07.11.2013
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
	"merx/common"
)

func GetInventoryLocal(dealerid int, bsvkeyid int, vendorcode string, partnum string, db *sql.DB) (code int, resp []byte, err error) {
	//06.05.2013 naj - make sure we have a merxpo
	if partnum == "" {
		code = http.StatusBadRequest
		err = errors.New("Missing Part Number")
		return
	}

	vendorid, err := common.GetVendorID(db, bsvkeyid, vendorcode)
	if err != nil {
		if err.Error() == "Unable to locate Vendor" {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}
		return
	}

	itemid, superseded, nla, err := common.GetPart(db, vendorid, partnum)
	if err != nil {
		if err.Error() == "Unable to locate part number" {
			code = http.StatusNotFound
		} else {
			code = http.StatusInternalServerError
		}
		return
	}

	//07.12.2013 naj - setup response object
	var r inventory
	r.PartNumber = partnum
	r.VendorCode = vendorcode

	err = db.QueryRow("select Description, Cost, List, QtyAvail, Category from Items where ItemID = ?", itemid).Scan(&r.Description, &r.Cost, &r.MSRP, &r.Stock, &r.Category)

	switch {
	case err == sql.ErrNoRows:
		code = http.StatusNotFound
		err = errors.New("Part Not Found")
		return
	case err != nil:
		code = http.StatusInternalServerError
		return
	default:
		if nla {
			r.NLA = 1
		} else {
			r.NLA = 0
		}

		if superseded != "" {
			r.PartNumber = superseded
			r.Superseded = 1
		}
	}

	//07.12.2013 naj - JSON encode the response data
	resp, err = json.Marshal(r)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	code = http.StatusOK
	return
}
