/*
Purpose: This package handles updating purchase orders
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
	//"fmt"
	"net/http"
	"regexp"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//04.10.2012 ghh - this function is reponsible for parsing our Json package
//and placing into internal array for us to work with
func (p *POUpdates) ParsePackage(pkg []byte, db *sql.DB) (code int, err error) {
	//09.27.2013 naj - unpack the data and make sure we have the minimum required data.
	err = json.Unmarshal(pkg, &p.PurchaseOrders)

	if err != nil {
		code = http.StatusBadRequest
		return code, err
	}

	for i := 0; i < len(p.PurchaseOrders); i++ {
		c := p.PurchaseOrders[i]

		//09.27.2013 naj - check to see if we have a Merx PO Number
		matched, err := regexp.MatchString("^MERX-.*", c.MerxPO)

		if err != nil {
			code = http.StatusInternalServerError
			return code, err
		}

		if !matched || c.MerxPO == "" {
			code = http.StatusBadRequest
			err = errors.New("Missing MerxPO")
			return code, err
		}

		if c.Status <= 0 {
			code = http.StatusBadRequest
			err = errors.New("Missing Status")
			return code, err
		}

	}

	p.db = db
	return 0, nil
}

//this function actually takes the contents of our structure and processes it
func (p *POUpdates) ProcessPackage() (code int, err error) {
	db := p.db
	//10.04.2013 naj - start a transaction
	transaction, err := db.Begin()

	for x := 0; x < len(p.PurchaseOrders); x++ {
		a := p.PurchaseOrders[x]

		//09.27.2013 naj - get the poid
		poid := strings.TrimPrefix(a.MerxPO, "MERX-")
		poid = strings.TrimLeft(poid, "0")

		//09.27.2013 naj - update each PO
		_, err := transaction.Exec("update PurchaseOrders set Status = ? where POID = ?", a.Status, poid)

		if err != nil {
			code = http.StatusInternalServerError
			//10.04.2013 naj - rollback the transaction
			_ = transaction.Rollback()
			return code, err
		}

		//09.27.2013 naj - now loop through the packages and items and update the tracking data
		for y := 0; y < len(a.Boxes); y++ {
			b := a.Boxes[y]

			//09.30.2013 naj - insert a record for the box
			result, err := transaction.Exec("insert into ShippedBoxes (BoxNumber, TrackingNumber) values (?, ?)", b.BoxNumber, b.TrackingNumber)

			if err != nil {
				code = http.StatusInternalServerError
				//10.04.2013 naj - rollback the transaction
				_ = transaction.Rollback()
				return code, err
			}

			//09.30.2013 naj - get the boxid
			boxid, err := result.LastInsertId()

			if err != nil {
				code = http.StatusInternalServerError
				//10.04.2013 naj - rollback the transaction
				_ = transaction.Rollback()
				return code, err
			}

			//10.01.2013 naj - now loop thorough the shipped items and update the shipping details.
			for z := 0; z < len(b.Items); z++ {
				c := b.Items[z]

				//09.30.2013 naj - get the POItemID
				rows, err := db.Query("select POItemID, Quantity from PurchaseOrderItems where POID = ? and VendorCode = ? and PartNumber = ?", poid, c.VendorCode, c.PartNumber)

				if err != nil {
					code = http.StatusInternalServerError
					//10.04.2013 naj - rollback the transaction
					_ = transaction.Rollback()
					return code, err
				}

				for rows.Next() {
					var poitemid, qtyordered int
					err = rows.Scan(&poitemid, &qtyordered)

					if err != nil {
						code = http.StatusInternalServerError
						//10.04.2013 naj - rollback the transaction
						_ = transaction.Rollback()
						return code, err
					}

					//09.30.2013 naj - insert a record into the PurchaseOrderShipped table and
					_, err := transaction.Exec("insert into PurchaseOrderShipped (POItemID, BoxID, QtyShipped, Cost) values (?, ?, ?, ?)", poitemid, boxid, c.Qty, c.Cost)

					if err != nil {
						code = http.StatusInternalServerError
						//10.04.2013 naj - rollback the transaction
						_ = transaction.Rollback()
						return code, err
					}

					//10.01.2013 naj - Since the po being updated could have already been in the system and have items that were listed on back order
					//we need to check to see how many items were originally ordered and compare that to the number of items currently shipped.
					//If the quantity shipped equals the quantity ordered we need to remove any back order entrys for this item.
					var qtyshipped int

					err = db.QueryRow("select ifnull(sum(QtyShipped),0) from PurchaseOrderShipped where POItemID = ?", poitemid).Scan(&qtyshipped)

					switch {
					case err == sql.ErrNoRows:
						code = http.StatusInternalServerError
						err = errors.New("Could not locate shipped parts, this should never happen.")
						//10.04.2013 naj - rollback the transaction
						_ = transaction.Rollback()
						return code, err
					case err != nil:
						code = http.StatusInternalServerError
						//10.04.2013 naj - rollback the transaction
						_ = transaction.Rollback()
						return code, err
					}

					if qtyordered >= qtyshipped {
						//10.01.2013 naj - we have shipped the full quantity ordered, now we need to remove any back order entries for the current POItemID
						_, err := transaction.Exec("delete from PurchaseOrderBackOrder where POItemID = ?", poitemid)

						if err != nil {
							code = http.StatusInternalServerError
							//10.04.2013 naj - rollback the transaction
							_ = transaction.Rollback()
							return code, err
						}
					}
				}
			}
		}

		//10.01.2013 naj - now loop through the items that have not shipped and update the estimated ship date.
		for y := 0; y < len(a.Pending); y++ {
			b := a.Pending[y]

			//09.30.2013 naj - get the POItemID
			rows, err := db.Query("select POItemID from PurchaseOrderItems where POID = ? and VendorCode = ? and PartNumber = ?", poid, b.VendorCode, b.PartNumber)

			if err != nil {
				code = http.StatusInternalServerError
				//10.04.2013 naj - rollback the transaction
				_ = transaction.Rollback()
				return code, err
			}

			for rows.Next() {
				var poitemid int
				err = rows.Scan(&poitemid)

				if err != nil {
					code = http.StatusInternalServerError
					//10.04.2013 naj - rollback the transaction
					_ = transaction.Rollback()
					return code, err
				}

				//10.01.2013 naj - before we can update the back order items we need to figure out if there is already
				//a record that needs to be updated.
				var query string

				err = db.QueryRow("select POItemID from PurchaseOrderBackOrder where POItemID = ?", poitemid).Scan(&poitemid)

				switch {
				case err == sql.ErrNoRows:
					query = "insert into PurchaseOrderBackOrder (QtyPending, Cost, EstShipDate, POItemID) values (?, ?, ?, ?)"
				case err != nil:
					code = http.StatusInternalServerError
					//10.04.2013 naj - rollback the transaction
					_ = transaction.Rollback()
					return code, err
				default:
					query = "update PurchaseOrderBackOrder set QtyPending = ?, Cost = ?, EstShipDate = ? where POItemID = ?"
				}

				//10.01.2013 naj - update the estimated ship date for the pending items
				_, err := transaction.Exec(query, b.Qty, b.Cost, b.EstShipDate, poitemid)

				if err != nil {
					code = http.StatusInternalServerError
					//10.04.2013 naj - rollback the transaction
					_ = transaction.Rollback()
					return code, err
				}
			}
		}
	}

	err = transaction.Commit()
	if err != nil {
		code = http.StatusInternalServerError
		//10.04.2013 naj - rollback the transaction
		_ = transaction.Rollback()
		return code, err
	}

	code = http.StatusOK
	return code, nil
}
