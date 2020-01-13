/* Go IPP - IPP core protocol implementation in pure Go
/*
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Package documentation
*/

/*
Package goipp implements IPP core protocol, as defined by RFC 8010

It doesn't implement high-level operations, such as "print a document",
"cancel print job" and so on. It's scope is limited to proper generation
and parsing of IPP requests and responses.

    IPP protocol uses the following simple model:
    1. Send a request
    2. Receive a response

Request and response both has a similar format, represented here
by type Message, with the only difference, that Code field of
that Message is the Operation code in request and Status code
in response. So most of operations are common for request and
response messages

Example:
    package main

    import (
	    "bytes"
	    "net/http"
	    "os"

	    "github.com/alexpevzner/goipp"
    )

    const uri = "http://192.168.1.102:631"

    // Build IPP OpGetPrinterAttributes request
    func makeRequest() ([]byte, error) {
	    m := goipp.NewRequest(goipp.DefaultVersion, goipp.OpGetPrinterAttributes, 1)
	    m.Operation.Add(goipp.MakeAttribute("attributes-charset",
		    goipp.TagCharset, goipp.String("utf-8")))
	    m.Operation.Add(goipp.MakeAttribute("attributes-natural-language",
		    goipp.TagLanguage, goipp.String("en-US")))
	    m.Operation.Add(goipp.MakeAttribute("printer-uri",
		    goipp.TagURI, goipp.String(uri)))
	    m.Operation.Add(goipp.MakeAttribute("requested-attributes",
		    goipp.TagKeyword, goipp.String("all")))

	    return m.EncodeBytes()
    }

    // Check that there is no error
    func check(err error) {
	    if err != nil {
		    panic(err)
	    }
    }

    func main() {
	    request, err := makeRequest()
	    check(err)

	    resp, err := http.Post(uri, goipp.ContentType, bytes.NewBuffer(request))
	    check(err)

	    var respMsg goipp.Message

	    err = respMsg.Decode(resp.Body)
	    check(err)

	    respMsg.Print(os.Stdout, false)
    }
*/
package goipp
