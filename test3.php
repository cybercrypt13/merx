<?

$data ='[{"MerxPO":"MERX-00004","AccountNumber":"988655","PODate":"2013-10-04","PONumber":"PO-1234","BillToName":"Nizex Inc.","BillToAddress1":"1735 Pennsylvania Ave.","BillToAddress2":"None","BillToCity":"McDonough","BillToState":"GA","BillToZip":"30253","BillToCountry":"US","ShipToName":"John Smith","ShipToAddress1":"123 Peachtree St.","ShipToAddress2":"","ShipToCity":"Atlanta","ShipToState":"GA","ShipToZip":"30313","ShipToCountry":"US","PaymentMethod":0,"LastFour":"","Items":[{"VendorCode":"PRTUN","PartNumber":"123456","Qty":1},{"VendorCode":"PRTUN","PartNumber":"12357","Qty":2},{"VendorCode":"PRTUN","PartNumber":"12358","Qty":3}]},{"MerxPO":"MERX-00005","AccountNumber":"988655","PODate":"2013-10-04","PONumber":"PO-1234","BillToName":"Nizex Inc.","BillToAddress1":"1735 Pennsylvania Ave.","BillToAddress2":"None","BillToCity":"McDonough","BillToState":"GA","BillToZip":"30253","BillToCountry":"US","ShipToName":"John Smith","ShipToAddress1":"123 Peachtree St.","ShipToAddress2":"","ShipToCity":"Atlanta","ShipToState":"GA","ShipToZip":"30313","ShipToCountry":"US","PaymentMethod":0,"LastFour":"","Items":[{"VendorCode":"PRTUN","PartNumber":"123456","Qty":1},{"VendorCode":"PRTUN","PartNumber":"12357","Qty":2},{"VendorCode":"PRTUN","PartNumber":"12358","Qty":3}]}]';

echo $data."\n";
print_r(json_decode($data));
?>
