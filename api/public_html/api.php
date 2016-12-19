<?php
/*
 * The MIT License (MIT)
 * 
 *  Copyright (c) 2014 Stephen Parker (withaspark.com)
 * 
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *  
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */



try {
	// Get json post data
	$sData = file_get_contents('php://input');
	$Data = json_decode($sData);

	if (! isset($Data->from)) {
		throw new Exception('');
	}

	// Obscure from email addresses
	$Data->from = preg_replace('/^(.).*@/', '\\1*****@', $Data->from);

	$response = array(
		'status'=>'good',
		'response'=>0,
	);

	// Push record to firebase
	$curl = curl_init();
	curl_setopt($curl, CURLOPT_URL, '<YOUR FIREBASE>');
	curl_setopt($curl, CURLOPT_POST, true);
	curl_setopt($curl, CURLOPT_POSTFIELDS, json_encode($Data));
	curl_exec($curl);
	curl_close($curl);
} catch (Exception $e) {
	$response = [
		'status'=>'error',
		'response'=>0,
	];
}

echo json_encode($response);
