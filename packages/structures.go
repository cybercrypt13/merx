/*
Purpose: This file contains all the structure definitions for all packages.
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
)

//05.29.2013 naj - main structure for the purchase orders
type POSend struct {
	PurchaseOrders []PO
	bsvkeyid       int
	db             *sql.DB
}

type PO struct {
	POID 	         		string //only used when the vendor retrieves a PO
	DealerID					int
	Status					int
	AccountNumber   		string //only used when the vendor retrieves a PO
	PODate          		string //only used when the vendor retrieves a PO
	DealerPONumber  		string
	VendorPONumber  		string
	BillToFirstName 		string
	BillToLastName  		string
	BillToCompanyName		string
	BillToPhone     		string
	BillToEmail     		string
	BillToAddress1  		string
	BillToAddress2  		string
	BillToCity      		string
	BillToState     		string
	BillToZip       		string
	BillToCountry   		string
	ShipToFirstName 		string
	ShipToLastName  		string
	ShipToCompanyName		string
	ShipToPhone     		string
	ShipToEmail     		string
	ShipToAddress1  		string
	ShipToAddress2  		string
	ShipToCity      		string
	ShipToState     		string
	ShipToZip       		string
	ShipToCountry   		string
	PaymentMethod   		int
	LastFour        		string
	ShipMethod      		string
	Items           		[]item
	Units						[]unit
}

type item struct {
	VendorCode string
	PartNumber string
	Qty        int
}

type unit struct {
	VendorCode 	string
	OrderCode 	string
	ModelNumber string
	Year			int
	Colors 		string
	Details 		string
	Qty        	int
	ForCustomer	int
}

type AcceptedOrder struct {
	DealerPO  string
	InternalID    string
	DealerKey string
	ItemNotes []ItemNote
}

type ItemNote struct {
	VendorCode string
	PartNumber string
	Superceded int
	NLA        int
	Note       string
}

type Parts struct {
	ItemID     int
	VendorCode string
	PartNumber string
	Qty        int
}

//05.29.2013 naj - main structure for adding dealers
type AddDealers struct {
	Dealers []DealerData
	db      *sql.DB
}

type DealerData struct {
	AccountNumber string
	IPAddress     string //
}

//05.29.2013 naj - main structure for adding dealers
type DeleteDealers struct {
	Dealers []DealerData
	db      *sql.DB
}

type Dealers struct {
	AccountNumber string
}

type inventory struct {
	VendorCode  string
	PartNumber  string
	Description string
	MSRP        float32
	Cost        float32
	Category    string
	Stock       int
	NLA         int
	Superseded  int
	Message     string
	DealerKey   string
}

type POStatus 	struct {
	InternalID  string
	Status		int
	DealerPO    string
	DealerKey   string
	Boxes       []box
	Pending     []penditem
}

type box struct {
	BoxNumber      string
	TrackingNumber string
	VendorInvoice 	string
	DueDate		 	string
	Items          []shipitem
}

type shipitem struct {
	VendorCode string
	PartNumber string
	Qty        int
	Cost       float32
}

type penditem struct {
	VendorCode  string
	PartNumber  string
	Qty         int
	Cost        float32
	EstShipDate string
}

//09.27.2013 naj - structure for handling POStats Updates
type POUpdates struct {
	PurchaseOrders []POStatus
	db             *sql.DB
}

//07.02.2014 naj - structures for handling vehicle data
type measurements struct {
	Units         float32
	UnitOfMeasure string
}

type enginedetails struct {
	EngineModel         string
	EngineSerialNumber  string
	EngineDescription   string
	Displacement        measurements
	CylinderConfig      string
	NumberOfCylinders   int
	FuelType            string
	FuelInductionSystem string
	Emissions           string
}

type optiondetails struct {
	OptionName        string
	OptionDescription string
}

type VehicleData struct {
	VendorCode          string
	ModelNumber         string
	ModelYear           string
	ModelDescription    string
	VehicleNotes        string
	PrimaryColor        string
	SecondaryColor      string
	VehicleType         string
	VehicleTypeCategory string
	Drivetrain          string
	SeatingCapacity     int
	MaxOccupantWeight   measurements
	TransmissionType    string
	VIN                 string
	Weight              measurements
	WheelBase           measurements
	SeatHeight          measurements
	WarrantyEndDate     string
	WarrantyEndDistance measurements
	WarrantyTerms       string
	SaleDate            string
	Engine              []enginedetails
	Options             []optiondetails
	PropellerType       string
}
