/*
Purpose: This package handles retrieving purchase order statuses for the client
Written:	05.28.2013
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
//	"regexp"
//	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//07.22.2015 ghh - we're only looking for the internalid that was returned
//when they placed their order.  So we're just going to look it up and
//make sure its linked to the proper dealer account to prevent checking for
//orders you don't own.
func GetOrderStatus(dealerid int, internalid string, db *sql.DB) (code int, resp []byte, err error) {
	//06.05.2013 naj - make sure we have an internalid
	if internalid == "" {
		code = http.StatusBadRequest
		err = errors.New("Missing InternalID field")
		return
	}

	//matched, err := regexp.MatchString("^MERX-.*", merxpo)

	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	//if !matched {
	//	code = http.StatusBadRequest
	//	err = errors.New("Invalid PO format")
	//	return
	//}

	//06.03.2013 naj - initialize some variables
	var r POStatus
	code = http.StatusOK
	r.InternalID = internalid

	//if temp is not null the we will convert it to a string and load it into r.EstShipDate
	err = db.QueryRow(`select DealerPONumber from PurchaseOrders 
							where DealerID = ? and POID = ?`,
		dealerid, internalid).Scan(&r.DealerPO)

	switch {
	case err == sql.ErrNoRows:
		code = http.StatusNotFound
		err = errors.New("Could locate PO: " + internalid)
		return
	case err != nil:
		code = http.StatusInternalServerError
		return
	}

	//06.03.2013 naj - if the status is not 4 then we have no shipment data yet, so we can just return what data we do have.
	//07.22.2015 ghh - removed because new flow will allow requesting order updates anytime you want.
	//had to do this because its possible some parts of this order have already shipped and we still
	//need to allow retrieving information on remaining parts
	//if r.Status < 4 {
	//	resp, err = json.Marshal(r)
	//	if err != nil {
	//		code = http.StatusInternalServerError
	//		return
	//	}
	//	return
	//}

	//06.03.2013 naj - now get all parts and the associated boxes
	rows, err := db.Query("select ifnull(a.VendorID, ''), ifnull(a.PartNumber, ''), ifnull(b.BoxID, 0), "+
		"ifnull(b.QtyShipped, 0), ifnull(c.QtyPending, 0), ifnull(c.EstShipDate, ''), b.Cost "+
		"from PurchaseOrderItems a left outer join PurchaseOrderShipped b on a.POItemID = b.POItemID "+
		"left outer join PurchaseOrderBackOrder c on a.POItemID = c.POItemID where "+
		"a.POID = ? order by b.BoxID", internalid)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}

	//06.03.2013 naj - initialize the boxid and counter variables and a slice to hold our items
	var currentboxid int = 0
	x := 0
	y := 0
	z := 0
	items := make([]shipitem, 0, 100)
	boxes := make([]box, 0, 100)
	pend := make([]penditem, 0, 100)

	//07.10.2013 naj - Loop through the results
	for rows.Next() {
		var partnumber, estship string
		var vendorid, boxid, qtyshipped, qtypending int
		var cost float32

		err = rows.Scan(&vendorid, &partnumber, &boxid, &qtyshipped, &qtypending, &estship, &cost)

		if err != nil {
			code = http.StatusInternalServerError
			return
		}

		//07.12.2013 naj - this is where we organized the data in to the appropriate arrays
		switch {
		case currentboxid != boxid:
			//06.05.2013 naj - if this is not the first box then we need to increment the y counter,
			//reset the x counter, and reset the items slice
			if currentboxid > 0 {
				y++
				x = 0
				items = make([]shipitem, 1, 100)
			}

			//06.03.2013 naj - update currentboxid
			currentboxid = boxid

			//06.03.2013 naj - add the current parts data to the items slice
			if qtyshipped > 0 {
				items = items[0 : len(items)+1]
				items[x].VendorID = vendorid
				items[x].PartNumber = partnumber
				items[x].Qty = qtyshipped
				items[x].Cost = cost
			}

			//07.10.2013 naj - check to see if we have any pending items.
			if qtypending > 0 {
				//07.10.2013 naj - first make sure the slice is not full
				if z >= cap(pend) {
					newPend := make([]penditem, len(pend), (cap(pend)+1)*2)
					copy(newPend, pend)
					pend = newPend
				}

				pend = pend[0 : len(pend)+1]
				pend[z].VendorID = vendorid
				pend[z].PartNumber = partnumber
				pend[z].Qty = qtypending
				pend[z].EstShipDate = estship
				pend[z].Cost = cost
				z++
			}

			//06.03.2013 naj - first check to make sure the slice is not full
			if y >= cap(boxes) {
				//06.03.2013 naj - our slice is too small to handle anymore data
				//So we will make a new slice with and double the cap.
				newBoxes := make([]box, len(boxes), (cap(boxes)+1)*2)
				copy(newBoxes, boxes)
				boxes = newBoxes
			}

			//07.10.2013 naj - If the slice length is less then the counter we will need to expand it
			boxes = boxes[0 : len(boxes)+1]
			boxes[y].Items = items

			//06.03.2013 naj - Since this is a new box get the tracking and box numbers from the database
			err = db.QueryRow("select BoxNumber, TrackingNumber from ShippedBoxes where BoxID = ?", boxid).Scan(&boxes[y].BoxNumber, &boxes[y].TrackingNumber)
			switch {
			case err == sql.ErrNoRows:
				code = http.StatusNotFound
				err = errors.New("Could not locate boxes for PO:" + internalid)
				return
			case err != nil:
				code = http.StatusInternalServerError
				return
			}

		default:
			//06.05.2013 naj - increment the x counter
			x++
			//06.03.2013 naj - first check to make sure the slice is not full
			if x >= cap(items) {
				//06.03.2013 naj - our slice is too small to handle anymore data
				//So we will make a new slice with and double the cap.
				newItems := make([]shipitem, len(items), (cap(items)+1)*2)
				copy(newItems, items)
				items = newItems
			}

			//06.03.2013 naj - add the current parts data to the items slice
			if qtyshipped > 0 {
				items = items[0 : len(items)+1]
				items[x].VendorID = vendorid
				items[x].PartNumber = partnumber
				items[x].Qty = qtyshipped
				items[x].Cost = cost
			}

			//07.10.2013 naj - check to see if we have any pending items.
			if qtypending > 0 {
				//07.10.2013 naj - first make sure the slice is not full
				if z >= cap(pend) {
					newPend := make([]penditem, len(pend), (cap(pend)+1)*2)
					copy(newPend, pend)
					pend = newPend
				}

				pend = pend[0 : len(pend)+1]
				pend[z].VendorID = vendorid
				pend[z].PartNumber = partnumber
				pend[z].Qty = qtypending
				pend[z].EstShipDate = estship
				pend[z].Cost = cost
				z++
			}

			//06.03.2013 naj - save the items into the current boxes slice
			if len(items) > 0 {
				boxes = boxes[0 : len(boxes)+1]
				boxes[y].Items = items
			}
		}
	}

	//06.03.2013 naj - now add our boxes to the POStatus structure and JSON encode our response.
	if len(boxes) > 0 {
		r.Boxes = boxes
	}

	//07.10.2013 naj - add the pending items to the response
	if len(pend) > 0 {
		r.Pending = pend
	}

	resp, err = json.Marshal(r)
	if err != nil {
		code = http.StatusInternalServerError
		return
	}
	return
}
