/*
Purpose: 	This package handles incomming purchase orders
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

package packages

import (
	"database/sql"
	"encoding/json"
	"errors"
	"strconv"
//	"strings"

	_ "github.com/go-sql-driver/mysql"
)

//04.10.2012 ghh - this function is reponsible for parsing our Json package
//and placing into internal array for us to work with
func (p *POSend) ParsePackage(pkg []byte, 
										db *sql.DB, 
										bsvkeyid int) error {

	//first we take the json string and convert to a dictionary
	err := json.Unmarshal(pkg, &p.PurchaseOrders)

	if err != nil {
		return err
	}

	//07.09.2013 naj - set the db and bsvkeyid variables
	p.db = db
	p.bsvkeyid = bsvkeyid
	return nil
}

//04.10.2012 ghh - this function deals with making sure the package we received
//has enough data in it to work with.
func (p *POSend) VerifyPackage() error {
	var currentpo, errtext string
	errtext = ""

	for i := 0; i < len(p.PurchaseOrders); i++ {
		c := p.PurchaseOrders[i]

		if c.DealerPONumber == "" {
			currentpo = ""
			errtext += " Missing DealerPONumber\n"
		} else {
			currentpo = c.DealerPONumber
		}

		//07.21.2015 ghh - ( we need to verify that we have a blank
		//item and unit array and if so return error
		if len(c.Items) == 0 && len(c.Units) == 0{
			errtext += currentpo + " PO Has No Parts Or Units\n"
		} 
	}

	if errtext != "" {
		return errors.New(errtext)
	} else {
		return nil
	}
}

//this function actually takes the contents of our structure and processes it
func (p *POSend) ProcessPackage(dealerid int, dealerkey string) ([]byte, error) {
	db := p.db
	//10.04.2013 naj - start a transaction
	transaction, err := db.Begin()
	if err != nil {
		return nil, err
	}

	//06.02.2013 naj - make a slice to hold the purchase orders
	r := make([]AcceptedOrder, 0, len(p.PurchaseOrders))

	//06.05.2015 ghh -because the system has the ability to push more than one purchase
	//order through at the same time it will loop through our array and process each
	//one separately
	for i := 0; i < len(p.PurchaseOrders); i++ {
		//06.02.2013 naj - stick the current PO into a new variable to keep the name short.
		c := p.PurchaseOrders[i]


		//06.02.2013 naj - put the current PONumber into the response
		r = r[0 : len(r)+1]
		r[i].DealerPO = c.DealerPONumber

		//06.10.2014 naj - check to see if the po is already in the system.
		//If it is and it's not processed yet, delete the the po and re-enter it.
		//If it is and it's processed return an error.
		var result sql.Result
		var temppoid int
		var tempstatus int

		//06.02.2015 ghh - first we grab the Ponumber that is being sent to use and we're going to see
		//if it has already been processed by the vendor
		err = transaction.QueryRow(`select ifnull(POID, 0 ), ifnull( Status, 0 ) 
											from PurchaseOrders 
											where DealerID = ? and DealerPONumber = ?`,
					dealerid, c.DealerPONumber).Scan(&temppoid, &tempstatus)

		//case err == sql.ErrNoRows:
		//if we have a PO already there and its not been processed yet by the vendor then we're going
		//to delete it as we're uploading it a second time.
		if temppoid > 0 { 
			if tempstatus == 0 { //has it been processed by vendor yet?
				result, err = transaction.Exec(`delete from PurchaseOrders 
															where DealerID=? 
															and DealerPONumber=? `, dealerid, c.DealerPONumber )
				if err != nil {
					return nil, err
				}

				//now delete the items from the old $_POST[
				result, err = transaction.Exec(`delete from PurchaseOrderItems 
															where POID=? `, temppoid )
				if err != nil {
					return nil, err

					//08.06.2015 ghh - delete units from linked units table
					result, err = transaction.Exec(`delete from PurchaseOrderUnits 
																where POID=? `, temppoid )
					if err != nil {
						return nil, err
					}
				}
			}

			//if we get here then we must have found an existing PO so lets log it and return
			if tempstatus > 0 {
				err = errors.New("Error: 16207 Purchase order already sent and pulled by vendor.")
				return nil, err
			}

			if err != sql.ErrNoRows {
				//if there was an error then return it
				if err != nil {
					return nil, err
				}
			}


			//06.02.2013 naj - create the PO record in the database.
			result, err = transaction.Exec(`insert into PurchaseOrders (
				DealerID, BSVKeyID, DealerPONumber, POReceivedDate, BillToFirstName, BillToLastName, BillToCompanyName, 
				BillToAddress1, BillToAddress2, BillToCity, BillToState, BillToZip, 
				BillToCountry, BillToPhone, BillToEmail, 
				ShipToFirstName, ShipToLastName, ShipToCompanyName, ShipToAddress1,
				ShipToAddress2, ShipToCity, ShipToState, ShipToZip, ShipToCountry, 
				ShipToPhone, ShipToEmail,  
				PaymentMethod, LastFour, ShipMethod) values 
				(?, ?, curdate(), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, 
				?, ?, ?, ?, ?, ?, ? )`, 
				dealerid, c.BSVKeyID, c.DealerPONumber,
				c.BillToFirstName, c.BillToLastName, c.BillToCompanyName, c.BillToAddress1, 
				c.BillToAddress2, c.BillToCity, c.BillToState, c.BillToZip, c.BillToCountry, 
				c.BillToPhone, c.BillToEmail,
				c.ShipToFirstName, c.ShipToLastName, c.ShipToCompanyName, c.ShipToAddress1, 
				c.ShipToAddress2, c.ShipToCity, c.ShipToState, c.ShipToZip, c.ShipToCountry, 
				c.ShipToPhone, c.ShipToEmail, c.PaymentMethod, c.LastFour, c.ShipMethod )

			if err != nil {
				//10.04.2013 naj - rollback transaction
				_ = transaction.Rollback()
				return nil, err
			}

			//06.02.2013 naj - get the POID assigned to the PO
			poid, err := result.LastInsertId()

			//06.02.2013 naj - format the POID and put the assigned POID into the response
			temp := strconv.FormatInt(poid, 10)

			r[i].InternalID = temp
			r[i].DealerKey = dealerkey

			if err != nil {
				//10.04.2013 naj - rollback transaction
				_ = transaction.Rollback()
				return nil, err
			}

			//06.05.2015 ghh - now loop through the items array and insert all the parts for
			//the order
			for j := 0; j < len(c.Items); j++ {
				//06.02.2013 naj - attach the parts to the current PO.
				_, err := transaction.Exec(`insert into PurchaseOrderItems (POID, PartNumber, VendorID, 
													Quantity) value (?, ?, ?, ?)`, 
													poid, c.Items[j].PartNumber, c.Items[j].VendorID, 
													c.Items[j].Qty)
				if err != nil {
					//10.04.2013 naj - rollback transaction
					_ = transaction.Rollback()
					return nil, err
				}

					//08.06.2015 ghh - ( now that we've written the line into the table we need to
					//query a few things in order to build a proper response to send back.  Things
					//we want to know are how many will ship, any supersession or other known info
					//current cost...

			}


			//07.21.2015 ghh - now loop through the list of units and add them to the PO
			for j := 0; j < len(c.Units); j++ {
				//06.02.2013 naj - attach the parts to the current PO.
				_, err := transaction.Exec(`insert into PurchaseOrderUnits (POID, ModelNumber, Year,
													VendorID, OrderCode, Colors, Details 
													Quantity) value (?, ?, ?, ?, ?, ?, ?, ?)`, 
													poid, c.Units[j].ModelNumber, c.Units[j].Year, 
													c.Units[j].VendorID, c.Units[j].OrderCode,
													c.Units[j].Colors, c.Units[j].Details,
													c.Units[j].Qty)
				if err != nil {
					//10.04.2013 naj - rollback transaction
					_ = transaction.Rollback()
					return nil, err
				}
			}
		}
	}

	//06.05.2015 ghh - now we'll take the array and marshal it back into a json
	//array to be returned to client
	if len(r) > 0 {
		//06.02.2013 naj - JSON Encode the response data.
		resp, err := json.Marshal(r)

		if err != nil {
			//10.04.2013 naj - rollback transaction
			_ = transaction.Rollback()
			return nil, err
		}

		//10.04.2013 naj - commit the transaction
		err = transaction.Commit()
		if err != nil {
			//10.04.2013 naj - rollback transaction
			_ = transaction.Rollback()
			return nil, err
		}

		return resp, nil
	} else {
		//10.04.2013 naj - rollback transaction
		_ = transaction.Rollback()
		return nil, errors.New("No valid parts were in the purchase order")
		}

}
