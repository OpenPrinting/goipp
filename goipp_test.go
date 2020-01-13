/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 */

package goipp

import (
	"bytes"
	"reflect"
	"testing"
	"time"
)

// The good message - 1
var good_message_1 = []byte{
	0x01, 0x01, // IPP version
	0x00, 0x02, // Print-Job operation
	0x00, 0x00, 0x00, 0x01, // Request ID

	uint8(TagOperationGroup),

	uint8(TagCharset),
	0x00, 0x12, // Name length + name
	'a', 't', 't', 'r', 'i', 'b', 'u', 't', 'e', 's', '-',
	'c', 'h', 'a', 'r', 's', 'e', 't',
	0x00, 0x05, // Value length + value
	'u', 't', 'f', '-', '8',

	uint8(TagLanguage),
	0x00, 0x1b, // Name length + name
	'a', 't', 't', 'r', 'i', 'b', 'u', 't', 'e', 's', '-',
	'n', 'a', 't', 'u', 'r', 'a', 'l', '-', 'l', 'a', 'n',
	'g', 'u', 'a', 'g', 'e',
	0x00, 0x02, // Value length + value
	'e', 'n',

	uint8(TagURI),
	0x00, 0x0b, // Name length + name
	'p', 'r', 'i', 'n', 't', 'e', 'r', '-', 'u', 'r', 'i',
	0x00, 0x1c, // Value length + value
	'i', 'p', 'p', ':', '/', '/', 'l', 'o', 'c', 'a', 'l',
	'h', 'o', 's', 't', '/', 'p', 'r', 'i', 'n', 't', 'e',
	'r', 's', '/', 'f', 'o', 'o',

	uint8(TagJobGroup),

	uint8(TagBeginCollection),
	0x00, 0x09, // Name length + name
	'm', 'e', 'd', 'i', 'a', '-', 'c', 'o', 'l',
	0x00, 0x00, // No value

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0a, // Value length + value
	'm', 'e', 'd', 'i', 'a', '-', 's', 'i', 'z', 'e',

	uint8(TagBeginCollection),
	0x00, 0x00, // Name length + name
	0x00, 0x00, // No value

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0b, // Value length + value
	'x', '-', 'd', 'i', 'm', 'e', 'n', 's', 'i', 'o', 'n',

	uint8(TagInteger),
	0x00, 0x00, // No name
	0x00, 0x04, // Value length + value
	0x00, 0x00, 0x54, 0x56,

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0b, // Value length + value
	'y', '-', 'd', 'i', 'm', 'e', 'n', 's', 'i', 'o', 'n',

	uint8(TagInteger),
	0x00, 0x00, // No name
	0x00, 0x04, // Value length + value
	0x00, 0x00, 0x6d, 0x24,

	uint8(TagEndCollection),
	0x00, 0x00, // No name
	0x00, 0x00, // No value

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0b, // Value length + value
	'm', 'e', 'd', 'i', 'a', '-', 'c', 'o', 'l', 'o', 'r',

	uint8(TagKeyword),
	0x00, 0x00, // No name
	0x00, 0x04, // Value length + value
	'b', 'l', 'u', 'e',

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0a, // Value length + value
	'm', 'e', 'd', 'i', 'a', '-', 't', 'y', 'p', 'e',

	uint8(TagKeyword),
	0x00, 0x00, // No name
	0x00, 0x05, // Value length + value
	'p', 'l', 'a', 'i', 'n',

	uint8(TagEndCollection),
	0x00, 0x00, // No name
	0x00, 0x00, // No value

	uint8(TagBeginCollection),
	0x00, 0x00, // No name
	0x00, 0x00, // No value

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0a, // Value length + value
	'm', 'e', 'd', 'i', 'a', '-', 's', 'i', 'z', 'e',

	uint8(TagBeginCollection),
	0x00, 0x00, // Name length + name
	0x00, 0x00, // No value

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0b, // Value length + value
	'x', '-', 'd', 'i', 'm', 'e', 'n', 's', 'i', 'o', 'n',

	uint8(TagInteger),
	0x00, 0x00, // No name
	0x00, 0x04, // Value length + value
	0x00, 0x00, 0x52, 0x08,

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0b, // Value length + value
	'y', '-', 'd', 'i', 'm', 'e', 'n', 's', 'i', 'o', 'n',

	uint8(TagInteger),
	0x00, 0x00, // No name
	0x00, 0x04, // Value length + value
	0x00, 0x00, 0x74, 0x04,

	uint8(TagEndCollection),
	0x00, 0x00, // No name
	0x00, 0x00, // No value

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0b, // Value length + value
	'm', 'e', 'd', 'i', 'a', '-', 'c', 'o', 'l', 'o', 'r',

	uint8(TagKeyword),
	0x00, 0x00, // No name
	0x00, 0x05, // Value length + value
	'p', 'l', 'a', 'i', 'd',

	uint8(TagMemberName),
	0x00, 0x00, // No name
	0x00, 0x0a, // Value length + value
	'm', 'e', 'd', 'i', 'a', '-', 't', 'y', 'p', 'e',

	uint8(TagKeyword),
	0x00, 0x00, // No name
	0x00, 0x06, // Value length + value
	'g', 'l', 'o', 's', 's', 'y',

	uint8(TagEndCollection),
	0x00, 0x00, // No name
	0x00, 0x00, // No value

	uint8(TagEnd),
}

// The good message - 2
var good_message_2 = []byte{
	0x01, 0x01, // IPP version
	0x00, 0x02, // Print-Job operation
	0x00, 0x00, 0x00, 0x01, // Request ID

	uint8(TagOperationGroup),

	uint8(TagInteger),
	0x00, 0x1f, // Name length + name
	'n', 'o', 't', 'i', 'f', 'y', '-', 'l', 'e', 'a', 's', 'e',
	'-', 'd', 'u', 'r', 'a', 't', 'i', 'o', 'n', '-', 's', 'u',
	'p', 'p', 'o', 'r', 't', 'e', 'd',
	0x00, 0x04, // Value length + value
	0x00, 0x00, 0x00, 0x01,

	uint8(TagRange),
	0x00, 0x00, // No name
	0x00, 0x08, // Value length + value
	0x00, 0x00, 0x00, 0x10,
	0x00, 0x00, 0x00, 0x20,

	uint8(TagEnd),
}

// The bad message - 1
var bad_message_1 = []byte{
	0x01, 0x01, // IPP version */
	0x00, 0x02, // Print-Job operation */
	0x00, 0x00, 0x00, 0x01, // Request ID */

	uint8(TagOperationGroup),

	uint8(TagCharset),

	0x00, 0x12, // Name length + name
	'a', 't', 't', 'r', 'i', 'b', 'u', 't', 'e', 's', '-',
	'c', 'h', 'a', 'r', 's', 'e', 't',
	0x00, 0x05, // Value length + value
	'u', 't', 'f', '-', '8',

	uint8(TagLanguage),
	0x00, 0x1b, // Name length + name
	'a', 't', 't', 'r', 'i', 'b', 'u', 't', 'e', 's', '-',
	'n', 'a', 't', 'u', 'r', 'a', 'l', '-', 'l', 'a', 'n',
	'g', 'u', 'a', 'g', 'e',
	0x00, 0x02, // Value length + value
	'e', 'n',

	uint8(TagURI),
	0x00, 0x0b, // Name length + name
	'p', 'r', 'i', 'n', 't', 'e', 'r', '-', 'u', 'r', 'i',
	0x00, 0x1c, // Value length + value
	'i', 'p', 'p', ':', '/', '/', 'l', 'o', 'c', 'a', 'l',
	'h', 'o', 's', 't', '/', 'p', 'r', 'i', 'n', 't', 'e',
	'r', 's', '/', 'f', 'o', 'o',

	uint8(TagJobGroup),

	uint8(TagBeginCollection),
	0x00, 0x09, // Name length + name
	'm', 'e', 'd', 'i', 'a', '-', 'c', 'o', 'l',
	0x00, 0x00, // No value

	uint8(TagBeginCollection),
	0x00, 0x0a, // Name length + name
	'm', 'e', 'd', 'i', 'a', '-', 's', 'i', 'z', 'e',
	0x00, 0x00, // No value

	uint8(TagInteger),
	0x00, 0x0b, // Name length + name
	'x', '-', 'd', 'i', 'm', 'e', 'n', 's', 'i', 'o', 'n',
	0x00, 0x04, // Value length + value
	0x00, 0x00, 0x54, 0x56,

	uint8(TagInteger),
	0x00, 0x0b, // Name length + name
	'y', '-', 'd', 'i', 'm', 'e', 'n', 's', 'i', 'o', 'n',
	0x00, 0x04, // Value length + value
	0x00, 0x00, 0x6d, 0x24,

	uint8(TagEndCollection),
	0x00, 0x00, // No name
	0x00, 0x00, // No value

	uint8(TagEndCollection),
	0x00, 0x00, // No name
	0x00, 0x00, // No value

	uint8(TagEnd),
}

// Check that err == nil
func assertNoError(t *testing.T, err error) {
	if err != nil {
		t.Errorf("%s", err)
	}
}

// Check that err != nil
func assertWithError(t *testing.T, err error) {
	if err == nil {
		t.Errorf("Error expected")
	}
}

// Check that value type is as specified
func assertValueType(t *testing.T, val Value, typ Type) {
	if val.Type() != typ {
		t.Errorf("%s: type is %s, must be %s", reflect.TypeOf(val).Name(), val.Type(), typ)
	}
}

func assertDataSize(t *testing.T, data []byte, size int) {
	if len(data) != size {
		t.Errorf("data size must be %d, present %d", size, len(data))
	}
}

// Check that encode() works without error and returns expected size
func assertEncodeSize(t *testing.T, encode func() ([]byte, error), size int) {
	data, err := encode()
	assertNoError(t, err)
	assertDataSize(t, data, size)
}

// Check that decode() works without error and returns expected value
func assertDecode(t *testing.T, data []byte, expected Value) {
	val, err := expected.decode(data)
	assertNoError(t, err)

	if !ValueEqual(val, expected) {
		t.Errorf("decode: expected %s, present %s", val, expected)
	}
}

// Check that decode returns error
func assertDecodeErr(t *testing.T, data []byte, val Value) {
	_, err := val.decode(data)
	if err == nil {
		t.Errorf("decode: expected error")
	}
}

// Test Void Value
func TestVoidValue(t *testing.T) {
	var v Void

	assertValueType(t, v, TypeVoid)
	assertEncodeSize(t, v.encode, 0)

	assertDecode(t, []byte{}, Void{})
	assertDecode(t, []byte{1, 2, 3, 4}, Void{})
}

// Test Integer Value
func TestIntegerValue(t *testing.T) {
	var v Integer

	assertValueType(t, v, TypeInteger)
	assertEncodeSize(t, v.encode, 4)

	assertDecode(t, []byte{1, 2, 3, 4}, Integer(0x01020304))
	assertDecodeErr(t, []byte{1, 2, 3}, Integer(0))
}

// Test Boolean Value
func TestBooleanValue(t *testing.T) {
	var v Boolean

	assertValueType(t, v, TypeBoolean)
	assertEncodeSize(t, v.encode, 1)

	assertDecode(t, []byte{0}, Boolean(false))
	assertDecode(t, []byte{1}, Boolean(true))
	assertDecodeErr(t, []byte{1, 2, 3}, Integer(0))
}

// Test String Value
func TestStringValue(t *testing.T) {
	var v String

	assertValueType(t, v, TypeString)
	assertEncodeSize(t, v.encode, 0)

	v = "12345"
	assertEncodeSize(t, v.encode, 5)

	assertDecode(t, []byte{}, String(""))
	assertDecode(t, []byte("hello"), String("hello"))
}

// Test Time Value
func TestDateTimeValue(t *testing.T) {
	var v Time

	assertValueType(t, v, TypeDateTime)
	assertEncodeSize(t, v.encode, 11)

	tm := time.Date(2020, 1, 13, 15, 35, 12, 300000000, time.UTC)

	v = Time{tm}
	data, _ := v.encode()

	assertDecode(t, data, Time{tm})
	assertDecodeErr(t, []byte{1, 2, 3}, Time{})
}

// Test Resolution value
func TestResolutionValue(t *testing.T) {
	v := Resolution{100, 100, UnitsDpi}

	assertValueType(t, v, TypeResolution)
	assertEncodeSize(t, v.encode, 9)

	data, _ := v.encode()
	assertDecode(t, data, v)
	assertDecodeErr(t, []byte{1, 2, 3}, Resolution{})
}

// Test Range value
func TestRangeValue(t *testing.T) {
	v := Range{100, 200}

	assertValueType(t, v, TypeRange)
	assertEncodeSize(t, v.encode, 8)

	data, _ := v.encode()
	assertDecode(t, data, v)
	assertDecodeErr(t, []byte{1, 2, 3}, Range{})
}

// Test Binary value
func TestBinaryValue(t *testing.T) {
	v := Binary([]byte("12345"))

	assertValueType(t, v, TypeBinary)
	assertEncodeSize(t, v.encode, 5)

	data, _ := v.encode()
	assertDecode(t, data, v)
}

// Test message decoding
func testDecode(t *testing.T, data []byte, mustFail bool) {
	var m Message
	err := m.Decode(bytes.NewBuffer(data))

	if mustFail {
		assertWithError(t, err)
	} else {
		assertNoError(t, err)
	}
}

func TestMessageDecode(t *testing.T) {
	testDecode(t, good_message_1, false)
	testDecode(t, good_message_2, false)
	testDecode(t, bad_message_1, true)

	/*
		//client := ipp.NewIPPClient("192.168.1.102", 631, "", "", false)
		//_, err := client.GetPrinterAttributes("printer", nil)
		//check(err)

		//url := "http://192.168.1.102:631"
		url := "http://localhost:631"

		rq := ipp.NewRequest(ipp.OperationGetPrinterAttributes, 1)
		rq.OperationAttributes[ipp.OperationAttributePrinterURI] = url
		rq.OperationAttributes[ipp.OperationAttributeRequestedAttributes] = ipp.DefaultPrinterAttributes

		data, err := rq.Encode()
		check(err)
		log_dump(good_message_1)

		var m Message
		err = m.Decode(bytes.NewBuffer(test_message))
		check(err)

		for _, a := range m.Operation {
			log_debug("%s: %v", a.Name, a.Values)
		}

		m.Print(os.Stdout, true)

		return

		rsp, err := http.Post(url, "application/ipp", bytes.NewBuffer(data))
		check(err)

		log_debug("status %s", rsp.Status)
		data, err = ioutil.ReadAll(rsp.Body)
		check(err)
		rsp.Body.Close()

		log_dump(data)

		dec := ipp.NewResponseDecoder(&buffer{data, 0})
		ipprsp, err := dec.Decode(nil)
		check(err)
		_ = ipprsp

		log_debug("%v", ipprsp)
	*/
}
