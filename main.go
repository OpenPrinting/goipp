/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 */

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	ipp "github.com/phin1x/go-ipp"
)

func check(err error) {
	if err != nil && err != io.EOF {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}
}

type buffer struct {
	data []byte
	off  int
}

func (b *buffer) Read(out []byte) (int, error) {
	av := len(b.data) - b.off
	if av == 0 {
		return 0, io.EOF
	}

	av = copy(out, b.data[b.off:])
	log_debug("0x%x: %d bytes of %d, next: 0x%x: %x",
		b.off, av, cap(out), b.off+av, out[:av])

	b.off += av

	return av, nil
}

func main() {
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
	log_dump(data)

	var m Message
	err = m.Decode(bytes.NewBuffer(data))
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
}
