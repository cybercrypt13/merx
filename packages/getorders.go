/*
Purpose: This package handles retrieving purchase orders for processing.
Written:	10.01.2013
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
	//"fmt"
	"net/http"
	"regexp"
//	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func GetOrder(db *sql.DB, date string, allorders string) (code int, resp []byte, err error) {
	//10.01.2013 naj - validate the date parameter
	if date != "" {
		matched, err := regexp.MatchString("^2[0-1][0-9]{2}-(0[1-9])|(1[0-2])-(0[1-9])|(1[0-9])|(2[0-9])|(3[0-1])", date)

		if err != nil {
			code = http.StatusInternalServerError
			return code, nil, err
		}

		if !matched {
			code = http.StatusBadRequest
			err = errors.New("Bad Date Format. Valid format is YYYY-MM-DD.")
			return code, nil, err
		}
	}

	//10.04.2013 naj - Start a transaction
	transaction, err := db.Begin()

	if err != nil {
		code = http.StatusInternalServerError
		return code, nil, err
	}

	//10.01.2013 naj - setup the query to retrieve the purchase orders
	var rows *sql.Rows

	if date != "" && allorders == "1" {
		//10.01.2013 naj - the vendors has requested all order regardless of whether or not they have already
		//been retrieved. We can only do this if they provided a date.
		rows, err = db.Query("select a.POID, a.DealerPONumber, a.POReceivedDate, a.BillToCompanyName, a.BillToAddress1, a.BillToAddress2, "+
			"a.BillToCity, a.BillToState, a.BillToZip, a.BillToCountry, a.ShipToCompanyName, a.ShipToAddress1, "+
			"a.ShipToAddress2, a.ShipToCity, a.ShipToState, a.ShipToZip, a.ShipToCountry, a.PaymentMethod, "+
			"a.LastFour, b.AccountNumber from PurchaseOrders a, DealerCredentials b where "+
			"a.DealerID = b.DealerID and a.POReceivedDate = ?", date)
	} else {
		//10.01.2013 naj - the vendor is requesting all orders that have not yet been retrieved.
		query := "select a.POID, a.DealerPONumber, a.POReceivedDate, a.BillToCompanyName, a.BillToAddress1, a.BillToAddress2, " +
			"a.BillToCity, a.BillToState, a.BillToZip, a.BillToCountry, a.ShipToCompanyName, a.ShipToAddress1, " +
			"a.ShipToAddress2, a.ShipToCity, a.ShipToState, a.ShipToZip, a.ShipToCountry, a.PaymentMethod, " +
			"a.LastFour, b.AccountNumber from PurchaseOrders a, DealerCredentials b where " +
			"a.DealerID = b.DealerID and a.Status = 0"

		//10.01.2013 naj - figure out if the vendor specified a date.
		if date != "" {
			rows, err = db.Query(query+" and a.POReceivedDate = ?", date)
		} else {
			rows, err = db.Query(query)
		}
	}

	if err != nil {
		code = http.StatusInternalServerError
		//10.04.2013 naj - rollback the transaction
		_ = transaction.Rollback()
		return code, nil, err
	}

	//10.03.2013 naj - create a slice to hold all of the purchase orders.
	purchaseorders := make([]PO, 0, 100)
	x := 0

	for rows.Next() {
		var p PO
		var poid string
		err = rows.Scan(&poid, &p.DealerPONumber, &p.PODate, &p.BillToCompanyName, &p.BillToAddress1, &p.BillToAddress2,
			&p.BillToCity, &p.BillToState, &p.BillToZip, &p.BillToCountry, &p.ShipToCompanyName, &p.ShipToAddress1,
			&p.ShipToAddress2, &p.ShipToCity, &p.ShipToState, &p.ShipToZip, &p.ShipToCountry, &p.PaymentMethod,
			&p.LastFour, &p.AccountNumber)

		if err != nil {
			code = http.StatusInternalServerError
			//10.04.2013 naj - rollback the transaction
			_ = transaction.Rollback()
			return code, nil, err
		}

		//10.03.2013 naj - format the poid and turn it into a merx po number
		p.POID = poid

		//10.04.2013 naj - flag the current po as processing.
		_, err = transaction.Exec("update PurchaseOrders set Status = 1 where POID = ?", poid)

		if err != nil {
			code = http.StatusInternalServerError
			//10.04.2013 naj - rollback the transaction
			_ = transaction.Rollback()
			return code, nil, err
		}

		//10.03.2013 naj - now we need to get all the parts on this po and added it to the po structure.
		rows2, err := db.Query("select VendorID, PartNumber, Quantity from PurchaseOrderItems where POID = ?", poid)

		if err != nil {
			code = http.StatusInternalServerError
			//10.04.2013 naj - rollback the transaction
			_ = transaction.Rollback()
			return code, nil, err
		}

		//10.03.2013 naj - create a slice to hold all of the parts on the order
		items := make([]item, 0, 100)
		y := 0

		for rows2.Next() {
			var i item

			err = rows2.Scan(&i.VendorID, &i.PartNumber, &i.Qty)

			if err != nil {
				code = http.StatusInternalServerError
				//10.04.2013 naj - rollback the transaction
				_ = transaction.Rollback()
				return code, nil, err
			}

			//10.03.2013 naj - make sure that the items slice is not full, if it is double it's size
			if y >= cap(items) {
				ni := make([]item, len(items), (cap(items)+1)*2)
				copy(ni, items)
				items = ni
			}

			//10.03.2013 naj - now expand the slice and add the current item to it
			items = items[0 : len(items)+1]
			items[y] = i
			y++
		}
		//10.03.2013 naj - since we just finished the retrieving all the items we need to add them to the current po
		p.Items = items

		//10.03.2013 naj - make sure the slice is not full, if it is double it's size
		if x >= cap(purchaseorders) {
			np := make([]PO, len(purchaseorders), (cap(purchaseorders)+1)*2)
			copy(np, purchaseorders)
			purchaseorders = np
		}

		//10.03.2013 naj - now expand the slice and add the current purchase order to it.
		purchaseorders = purchaseorders[0 : len(purchaseorders)+1]
		purchaseorders[x] = p
		x++
	}

	//10.04.2013 naj - if we have no purchase orders return no content.
	if len(purchaseorders) == 0 {
		code = http.StatusNoContent
		resp = []byte("No Pending Orders")
		err = errors.New("No Pending Orders")
		return code, resp, err
	}

	//10.03.2013 naj - JSON encode the purchase orders
	resp, err = json.Marshal(purchaseorders)

	if err != nil {
		code = http.StatusInternalServerError
		//10.04.2013 naj - rollback the transaction
		_ = transaction.Rollback()
		return code, nil, err
	}

	//10.04.2013 naj - commit the transaction and return
	err = transaction.Commit()

	if err != nil {
		code = http.StatusInternalServerError
		_ = transaction.Rollback()
		return code, nil, err
	}

	code = http.StatusOK
	return code, resp, err
}
