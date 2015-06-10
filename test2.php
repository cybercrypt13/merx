<?
$purchaseorders = array(
	array(
		"MerxPO"			=> 'MERX-00001',
		"Status"			=> 3,
		"EstShipDate"	=> '2013-10-31',
		"Boxes"			=> array(
			array(
				"BoxNumber"			=> "1",
				"TrackingNumber"	=> "1Z2134234242341",
				"Items"				=> array(
					array(
						"VendorCode"	=> "PRTUN",
						"PartNumber"	=> "123456",
						"Qty"				=> 1,
						"Cost"			=> 5.55
					),
					array(
						"VendorCode"	=> "PRTUN",
						"PartNumber"	=> "12358",
						"Qty"				=> 3,
						"Cost"			=> 1.00
					)
				)
			)
		),
		"Pending"		=> array(
			array(
				"VendorCode"	=> "PRTUN",
				"PartNumber"	=> "12357",
				"Qty"				=> 2,
				"Cost"			=> 10.25
			)
		)
	),
	array(
		"MerxPO"			=> 'MERX-00002',
		"Status"			=> 4,
		"EstShipDate"	=> '2013-10-31',
		"Boxes"			=> array(
			array(
				"BoxNumber"			=> "1",
				"TrackingNumber"	=> "1Z2134234242341",
				"Items"				=> array(
					array(
						"VendorCode"	=> "PRTUN",
						"PartNumber"	=> "123456",
						"Qty"				=> 1,
						"Cost"			=> 5.55
					),
					array(
						"VendorCode"	=> "PRTUN",
						"PartNumber"	=> "12357",
						"Qty"				=> 2,
						"Cost"			=> 10.25
					),
					array(
						"VendorCode"	=> "PRTUN",
						"PartNumber"	=> "12358",
						"Qty"				=> 3,
						"Cost"			=> 1.00
					)
				)
			)
		)
	)
);

$json = json_encode($purchaseorders);

echo "$json\n";
print_r(json_decode($json));

$purchaseorders = array(
	array(
		"MerxPO"			=> 'MERX-00004',
		"Status"			=> 3,
		"EstShipDate"	=> '2013-10-31',
		"Boxes"			=> array(
			array(
				"BoxNumber"			=> "1",
				"TrackingNumber"	=> "912134234242341",
				"Items"				=> array(
					array(
						"VendorCode"	=> "PRTUN",
						"PartNumber"	=> "12357",
						"Qty"				=> 2,
						"Cost"			=> 10.25
					)
				)
			)
		)
	)
);
$json = json_encode($purchaseorders);
echo "$json\n";
?>
