-- MySQL dump 10.14  Distrib 5.5.33a-MariaDB, for debian-linux-gnu (x86_64)
--
-- Host: localhost    Database: merx
-- ------------------------------------------------------
-- Server version	5.5.33a-MariaDB-1~raring-log

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `AuthorizedBSVKeys`
--

DROP TABLE IF EXISTS `AuthorizedBSVKeys`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `AuthorizedBSVKeys` (
  `BSVKeyID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique ID for each authorized BSV key',
  `BSVKey` char(64) NOT NULL DEFAULT '' COMMENT 'The BSV authorization key',
  `BSVVendorCode` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'Flag to tell merX if it needs to use the BSV specific vendor codes',
  PRIMARY KEY (`BSVKeyID`),
  UNIQUE KEY `BSVKey` (`BSVKey`),
  KEY `iBSVKey` (`BSVKey`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='Stores the authorized BSV Keys';
/*!40101 SET character_set_client = @saved_cs_client */;


--
-- Table structure for table `DealerCredentials`
--

DROP TABLE IF EXISTS `DealerCredentials`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `DealerCredentials` (
  `DealerID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique key for each dealer',
  `DealerKey` varchar(40) DEFAULT NULL COMMENT 'Stores the dealer key',
  `IPAddress` int(10) unsigned DEFAULT NULL COMMENT 'Stores the client ip as integer',
  `CreatedDateTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT 'When the record was created',
  `UpdatedDateTime` datetime NOT NULL DEFAULT '0000-00-00 00:00:00' COMMENT 'When the record was modified',
  `AccessedDateTime` datetime COMMENT 'When the client last connected to merx',
  `LastIPAddress` int(10) unsigned DEFAULT NULL COMMENT 'Hold the client''s last know IP Address',
  `Active` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'Turns on or off the clients Account',
  `AccountNumber` varchar(35) NOT NULL DEFAULT '' COMMENT 'Holds the dealer number',
  PRIMARY KEY (`DealerID`),
  KEY `iUUIDIPAddr` (`IPAddress`),
  KEY `iLogin` (`DealerKey`,`Active`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='This table stores the client system''s credentials';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `DealerDiscountLink`
--

DROP TABLE IF EXISTS `DealerDiscountLink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `DealerDiscountLink` (
  `DealerID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to DealerCredentials table',
  `ItemID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the Items table',
  `Discount` float DEFAULT '0' COMMENT 'The discount rate',
  UNIQUE KEY `iDealerCode` (`DealerID`,`ItemID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='This table links dealers to specific parts and discount percentages.';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `Items`
--

DROP TABLE IF EXISTS `Items`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Items` (
  `ItemID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique key for each part in merX',
  `VendorID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the VendorCode table',
  `PartNumber` varchar(30) NOT NULL DEFAULT '' COMMENT 'The part number',
  `Description` varchar(75) NOT NULL DEFAULT '' COMMENT 'The part description',
  `ManufPartNumber` varchar(25) comment 'original manufacturer part number',
  `ManufName` varchar(50) comment 'original manufacturer name',
  `SupersessionNumber` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the ItemID of the superseeding part',
  `NLA` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'No longer available flag, 0 = false, 1 = true',
  `CloseOut` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'Will not be available after inventory depleted 0 = false, 1 = true',
  `PriceCode` varchar(3) NOT NULL DEFAULT '' COMMENT 'Holds the price code that applies to this part',
  `Cost` decimal(13,3) NOT NULL DEFAULT '0.000' COMMENT 'This stores the basic cost of the item',
  `List` decimal(13,3) NOT NULL DEFAULT '0.000' COMMENT 'This stores the suggested retail price of the item',
  `MAP` decimal(13,3) NOT NULL DEFAULT '0.000' COMMENT 'This store the minimum advertise price of the item',
  `Category` varchar(50) NOT NULL DEFAULT '' COMMENT 'Hold category info',
  PRIMARY KEY (`ItemID`),
  KEY `iPart` (`VendorID`,`PartNumber`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='This tables holds the manufacturer/suppliers price file';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ItemStock`
--

DROP TABLE IF EXISTS `ItemStock`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ItemStock` (
  `ItemID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique key for each part in merX',
  `WarehouseID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the Warehouses',
  `Qty` decimal(13,3) NOT NULL DEFAULT '0.000' COMMENT 'This store the actual cost of the item',
  PRIMARY KEY (`ItemID`, `WarehouseID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='This tables holds the stock qty for each warehouse';



DROP TABLE IF EXISTS `ItemCost`;
CREATE TABLE `ItemCost` (
  `ItemID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique key for each part in merX',
  `DealerID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the Warehouses',
  `DealerCost` decimal(13,3) NOT NULL DEFAULT '0.000' COMMENT 'This store the actual cost of the item',
  PRIMARY KEY (`ItemID`, `DealerID`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='This tables holds cost per dealer per item';

--
-- Table structure for table `Warehouses`
--

DROP TABLE IF EXISTS `Warehouses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Warehouses` (
  `WarehouseID` int(10) unsigned NOT NULL auto_increment primary key COMMENT 'Warehouse ID for supported warehouses',
  `WarehouseName` varchar(50) NOT NULL COMMENT 'name of different warehouses',
  `WarehouseState` varchar(5) not null comment 'state the warehouse resides in',
  KEY ( `WarehouseName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='This tables holds the different supported warehouses';
--
-- Table structure for table `PriceCodesLink`
--

DROP TABLE IF EXISTS `PriceCodesLink`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PriceCodesLink` (
  `DealerID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to DealerCredentials table',
  `PriceCode` varchar(3) NOT NULL DEFAULT '' COMMENT 'Hold the price code',
  `Discount` float DEFAULT '0' COMMENT 'The discount rate',
  UNIQUE KEY `iDealerCode` (`DealerID`,`PriceCode`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='This table links dealers to specific price codes and discount percentages';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PurchaseOrderBackOrder`
--

DROP TABLE IF EXISTS `PurchaseOrderBackOrder`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PurchaseOrderBackOrder` (
  `BackOrderID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique ID for every pending entry',
  `POItemID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the specific item in PurchaseOrderItems',
  `QtyPending` int(11) NOT NULL DEFAULT '0' COMMENT 'The qty pending shipment',
  `Cost` decimal(13,3) NOT NULL DEFAULT '0.000' COMMENT 'This stores the cost that the dealer will pay for this part',
  `EstShipDate` date NOT NULL DEFAULT '0000-00-00' COMMENT 'Store the estimated ship date for the back ordered item',
  `Note` varchar(100) comment 'holds any special note transferred back from vendor',
  PRIMARY KEY (`BackOrderID`),
  KEY `iPOItemID` (`POItemID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='Stores items that are on backorder';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PurchaseOrderItems`
--

DROP TABLE IF EXISTS `PurchaseOrderItems`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PurchaseOrderItems` (
  `POItemID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique key for each part',
  `POID` int(10) unsigned NOT NULL COMMENT 'Links the items to a specific purchase order',
  `PartNumber` varchar(30) NOT NULL DEFAULT '' COMMENT 'Stores the item''s part number',
  `Quantity` int(11) DEFAULT NULL COMMENT 'Stores the quantity of parts ordered',
  `OrderType` tinyint(3) NOT NULL DEFAULT 2 COMMENT '1=For Customer, 2=For Stock, 3=Seasonal Order',
  `FillStatus` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0=ship what you have, 1=only ship if you can fill completely',
  `ItemID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the Items table',
  `VendorID` varchar(5) NOT NULL DEFAULT '' COMMENT 'The primary VendorID that was submitted on the PO',
  `Status` tinyint(3) NOT NULL DEFAULT 0 COMMENT '1=Superseded, 2=Obsolete, 3=Rejected',
  `Supersession` varchar(30) comment 'holds supersession number if one exists',
  PRIMARY KEY (`POItemID`),
  KEY `iPOID` (`POID`),
  KEY `iPOIDPartNumber` (`POID`,`PartNumber`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='This table stores the client''s purchase order data.';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `PurchaseOrderShipped`
--

DROP TABLE IF EXISTS `PurchaseOrderShipped`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PurchaseOrderShipped` (
  `ShippedID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique ID for every shipping entry.',
  `POItemID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the specific item in PuchaseOrderItems',
  `BoxID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to the specific box in ShippedBoxes',
  `QtyShipped` int(11) NOT NULL DEFAULT '0' COMMENT 'The qty that were put in the box',
  `Cost` decimal(13,3) NOT NULL DEFAULT '0.000' COMMENT 'This stores the cost that the dealer will pay for this part',
  PRIMARY KEY (`ShippedID`),
  KEY `iPOItemID` (`POItemID`),
  KEY `iBoxID` (`BoxID`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;





drop table if exists `PurchaseOrderUnits`;
create table PurchaseOrderUnits( POUnitID int unsigned not null auto_increment primary key, 
POID int unsigned not null comment 'links to PurchaseOrder table',
VendorID varchar(5) comment 'holds specific code for vendor', 
OrderCode varchar(25) comment 'vendor specific part number for this unit', 
ModelNumber varchar(50) comment 'vendor model number', 
Year int unsigned default 0 comment 'year if available', 
Colors text comment 'list of colors for unit', 
Details text comment 'special notes', 
Quantity int unsigned comment 'how many of this model is being ordered', 
ForCustomer tinyint unsigned default 0 comment '0=for stock, 1=for customer' ) 
engine=innodb DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci comment 'holds unit order information';



--
-- Table structure for table `PurchaseOrders`
--

DROP TABLE IF EXISTS `PurchaseOrders`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `PurchaseOrders` (
  `POID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique key for each purchase order',
  `DealerID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to DealerCredentials table',
  `BSVKeyID` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'Links to BSV table',
  `DealerPONumber` varchar(20) NOT NULL DEFAULT '' COMMENT 'Stores the client''s purchase order number',
  `DueDate` Date  COMMENT 'Due Date Returned By Vendor',
  `POReceivedDate` date NOT NULL DEFAULT '0000-00-00' COMMENT 'Stores the purchase order date received',
  `ShipToFirstName` varchar(50) NOT NULL DEFAULT '' COMMENT 'Stores the ship to contact name',
  `ShipToLastName` varchar(50) NOT NULL DEFAULT '' COMMENT 'Stores the ship to contact name',
  `ShipToCompanyName` varchar(50) NOT NULL DEFAULT '' COMMENT 'Stores the ship to contact name',
  `ShipToAddress1` varchar(50) NOT NULL DEFAULT '' COMMENT 'Stores the ship to address 1',
  `ShipToAddress2` varchar(50) NOT NULL DEFAULT '' COMMENT 'Stores the ship to address 2',
  `ShipToCity` varchar(50) NOT NULL DEFAULT '' COMMENT 'Stores the ship to city',
  `ShipToState` varchar(5) NOT NULL DEFAULT '' COMMENT 'Stores the ship to state/province code',
  `ShipToZip` varchar(15) NOT NULL DEFAULT '' COMMENT 'Stores the ship to postal code',
  `ShipToCountry` varchar(3) NOT NULL DEFAULT '' COMMENT 'Stores the ship to country code',
  `ShipToPhone` varchar(15) NOT NULL DEFAULT '' COMMENT 'Stores the ship to country code',
  `ShipToEmail` varchar(50) NOT NULL DEFAULT '' COMMENT 'Stores the billing postal code',
  `Status` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT 'Holds the PO Status 0 = new, 1 = processing, 2 = pulling, 3 = staging, 4 = shipped, 5 = rejected',
  `PaymentMethod` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0 = No method specified, 1 = VISA, 2 = Mastercard, 3 = American Express, 4 = Discover, 5 = NET',
  `LastFour` char(4) NOT NULL DEFAULT '' COMMENT 'Last four of creditcard on file',
  `ShipMethod` tinyint(3) unsigned NOT NULL DEFAULT '0' COMMENT '0 = No method specified, 1 = VISA, 2 = Mastercard, 3 = American Express, 4 = Discover, 5 = NET',
  `Discount` decimal(9,2) default 0 comment 'Holds any discounts sent back from vendor',
  PRIMARY KEY (`POID`),
  KEY `iClientPONumber` (`DealerID`,`DealerPONumber`),
  KEY `iGetOrders` (`POReceivedDate`,`Status`),
  KEY `iGetOrders2` (`Status`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='This table stores the client''s purchase order data.';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `ShippedBoxes`
--

DROP TABLE IF EXISTS `ShippedBoxes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `ShippedBoxes` (
  `BoxID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique ID for each box shipped',
  `BoxNumber` varchar(25) NOT NULL DEFAULT '1' COMMENT 'Stores the box number assigned by the vendor',
  `TrackingNumber` varchar(50) NOT NULL DEFAULT '' COMMENT 'Stores the boxes tracking number',
  `VendorInvoiceNumber` varchar(20) NOT NULL DEFAULT '' COMMENT 'Stores the vendor''s invoice number',
  `DueDate` Date  COMMENT 'Due Date Returned By Vendor',
  PRIMARY KEY (`BoxID`),
  KEY `iTracking` (`TrackingNumber`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8 COMMENT='This table stores the tracking data for each box shipped.';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VendorCodes`
--

DROP TABLE IF EXISTS `Vendors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `Vendors` (
  `VendorID` int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT 'Unique key for each VendorCode',
  `VendorName` varchar(50) NOT NULL COMMENT 'Standard VendorCode',
  PRIMARY KEY (`VendorID`),
  UNIQUE KEY `iVendorName` (`VendorName`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='This table holds the standard vendor codes';
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `VendorCredentials`
--

DROP TABLE IF EXISTS `VendorCredentials`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `VendorCredentials` (
  `VendorKey` char(64) NOT NULL DEFAULT '' COMMENT 'Stores the Vendors auth key',
  `IPAddress` int(10) unsigned DEFAULT NULL COMMENT 'Holds the vendors ip address'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Used for controlling access to Admin methods';
/*!40101 SET character_set_client = @saved_cs_client */;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;


drop table if exists `UnitModel`;
create table UnitModel( ModelID int unsigned not null auto_increment primary key, 
VendorID int unsigned not null comment 'links to Vendor Table',
OrderCode varchar(25) comment 'vendor specific part number for this unit', 
ModelNumber varchar(50) comment 'vendor model number', 
ModelNumberNoFormat varchar(50) comment 'vendor model number stripped of all formatting', 
VehicleType varchar(50) comment 'Street, Dirt, Atv, Car, Truck...', 
Year int unsigned default 0 comment 'year if available', 
Colors text comment 'list of colors for unit', 
Details text comment 'special notes') 
engine=innodb DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci comment 'holds unit model information';

drop table if exists `UnitModelImages`;
create table UnitModelImages( ImageID int unsigned not null auto_increment primary key, 
ModelID int unsigned not null comment 'links to Model Table',
ImageURL varchar(100) comment 'URL to image of unit') 
engine=innodb DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci comment 'holds unit model image urls';

drop table if exists `ItemImages`;
create table ItemImages( ImageID int unsigned not null auto_increment primary key, 
ItemID int unsigned not null comment 'links to Model Table',
ImageURL varchar(100) comment 'URL to image of unit') 
engine=innodb DEFAULT CHARACTER SET utf8 COLLATE utf8_general_ci comment 'holds unit model image urls';

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2013-10-07 11:14:57
insert into DealerCredentials values( null, '0d3d6381-0e02-11e5-9eb5-20c9d0478db9', null, now(), now(), now(), null, 1,'12345');
insert into AuthorizedBSVKeys values( null, '108b6a78-4027-447b-9b2d-a6c9b7da72dc', 0 );
alter table PurchaseOrders auto_increment=1000;
