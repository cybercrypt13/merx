<?php
$parts = array(
	array(
		"VendorID"	=> "PRTUN",
		"PartNumber"	=> "123456",
		"Qty"				=> 1
	),
	array(
		"VendorID"	=> "PRTUN",
		"PartNumber"	=> "12357",
		"Qty"				=> 2
	),
	array(
		"VendorID"	=> "PRTUN",
		"PartNumber"	=> "12358",
		"Qty"				=> 3
	)
);

$po = array(
	"DealerPONumber"	=> "PO-1234",
	"PODate"				=> "2013-05-29",
	"BillToName"		=> "Nizex Inc.",
	"BillToAddress1"	=> "1735 Pennsylvania Ave.",
	"BillToAddress2"	=> "None",
	"BillToCity"		=> "McDonough",
	"BillToState"		=> "GA",
	"BillToZip"			=> "30253",
	"BillToCountry"	=> "US",
	"ShipToName"		=> "John Smith",
	"ShipToAddress1"	=> "123 Peachtree St.",
	"ShipToAddress2"	=> "",
	"ShipToCity"		=> "Atlanta",
	"ShipToState"		=> "GA",
	"ShipToZip"			=> "30313",
	"ShipToCountry"	=> "US",
	"Items"				=> $parts
);

$url = "http://localhost:8000";
$fields = array(
	'accountnumber'	=> '12345',
	'bsvkey'				=> '108b6a78-4027-447b-9b2d-a6c9b7da72dc',
	'dealerkey'			=> '0d3d6381-0e02-11e5-9eb5-20c9d0478db9',
	//'dealerkey'			=> '',
	'data'				=> json_encode(array($po))
);

print_r($fields);
echo "\n\n";

$field_data = "";
foreach ($fields as $key => $value)
	$field_data .= "&$key=$value";


$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, "$url/sendorder");
curl_setopt($ch, CURLOPT_POST, true);
curl_setopt($ch, CURLOPT_POSTFIELDS, $field_data);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$result = curl_exec($ch);
$info = curl_getinfo($ch);
curl_close($ch);

if ($info['http_code'] == 200) {
	$pos = json_decode($result);
	echo "$result\n";
	
	foreach ($pos as $po) {
		echo "Dealer's PO# $po->DealerPO\n";
		echo "Merx PO# $po->MerxPO\n";
	}
} else {
	echo "Request Failed\n$result\n";
	print_r($info);
}
exit;
/*
*/

/*
//06.03.2013 naj - get po status
$getstr = "?UUID=123&DealerKey=abc&BSVKey=abc&MerxPO=MERX-00072";
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, "$url/postatus$getstr");
curl_setopt($ch, CURLOPT_POST, false);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$result = curl_exec($ch);
$info = curl_getinfo($ch);
curl_close($ch);

if ($info['http_code'] == 200) {
	$pos = json_decode($result);
	echo "$result\n";
	print_r($pos);	
} else {
	echo "Request Failed\n$result\n";
	print_r($info);
}
exit;
*/

//06.03.2013 naj - inventoryverify
$getstr = "?UUID=123&DealerKey=abc&BSVKey=abc&VendorID=PRTUN&PartNumber=12356";
$ch = curl_init();
curl_setopt($ch, CURLOPT_URL, "$url/inventoryverify$getstr");
curl_setopt($ch, CURLOPT_POST, false);
curl_setopt($ch, CURLOPT_RETURNTRANSFER, true);
$result = curl_exec($ch);
$info = curl_getinfo($ch);
curl_close($ch);

if ($info['http_code'] == 200) {
	$pos = json_decode($result);
	echo "$result\n";
	print_r($pos);	
} else {
	echo "Request Failed\n$result\n";
	print_r($info);
}
?>
