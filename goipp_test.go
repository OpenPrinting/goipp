/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 */

package goipp

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"testing"
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

func check(t *testing.T, err error, mustFail bool) {
	if err != nil {
		fmt.Printf("%s\n", err)
	}

	if mustFail {
		if err != nil {
			err = nil
		} else {
			err = errors.New("Non-nil error expected")
		}
	}

	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

func testDecode(t *testing.T, data []byte, mustFail bool) {
	log_dump(data)

	var m Message
	err := m.Decode(bytes.NewBuffer(data))
	check(t, err, mustFail)

	if err != nil {
		return
	}

	m.Print(os.Stdout, true)

	buf := bytes.NewBuffer(nil)
	err = m.Encode(buf)
	check(t, err, false)

	log_dump(buf.Bytes())

	var m2 Message
	err = m2.Decode(bytes.NewBuffer(buf.Bytes()))
	check(t, err, false)
	m2.Print(os.Stdout, true)
}

func TestGoipp(t *testing.T) {
	testDecode(t, good_message_1, false)
	//testDecode(t, good_message_2, false)
	//testDecode(t, bad_message_1, true)

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
