// Go IPP - IPP core protocol implementation in pure Go
//
// Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
// See LICENSE for license terms and conditions
//
// Package documentation

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
*/
package goipp
