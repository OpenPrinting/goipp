/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP formatter test
 */

package goipp

import (
	"strings"
	"testing"
)

// TestFmtAttribute runs Formatter.FmtAttribute tests
func TestFmtAttribute(t *testing.T) {
	type testData struct {
		attr   Attribute // Inpur attribute
		out    []string  // Expected output
		indent int       // Indentation
	}

	tests := []testData{
		// Simple test
		{
			attr: MakeAttr(
				"attributes-charset",
				TagCharset,
				String("utf-8")),
			out: []string{
				`ATTR "attributes-charset" charset: utf-8`,
			},
		},

		// Simple test with indentation
		{
			attr: MakeAttr(
				"attributes-charset",
				TagCharset,
				String("utf-8")),
			indent: 2,
			out: []string{
				`  ATTR "attributes-charset" charset: utf-8`,
			},
		},

		// Collection
		{
			attr: MakeAttrCollection("media-col",
				MakeAttrCollection("media-size",
					MakeAttribute("x-dimension",
						TagInteger, Integer(10160)),
					MakeAttribute("y-dimension",
						TagInteger, Integer(15240)),
				),
				MakeAttribute("media-left-margin",
					TagInteger, Integer(0)),
				MakeAttribute("media-right-margin",
					TagInteger, Integer(0)),
				MakeAttribute("media-top-margin",
					TagInteger, Integer(0)),
				MakeAttribute("media-bottom-margin",
					TagInteger, Integer(0)),
			),

			out: []string{
				`ATTR "media-col" collection: {`,
				`    MEMBER "media-size" collection: {`,
				`        MEMBER "x-dimension" integer: 10160`,
				`        MEMBER "y-dimension" integer: 15240`,
				`    }`,
				`    MEMBER "media-left-margin" integer: 0`,
				`    MEMBER "media-right-margin" integer: 0`,
				`    MEMBER "media-top-margin" integer: 0`,
				`    MEMBER "media-bottom-margin" integer: 0`,
				`}`,
			},
		},

		// 1SetOf Collection
		{
			attr: MakeAttr("media-size-supported",
				TagBeginCollection,
				Collection{
					MakeAttribute("x-dimension",
						TagInteger, Integer(20990)),
					MakeAttribute("y-dimension",
						TagInteger, Integer(29704)),
				},
				Collection{
					MakeAttribute("x-dimension",
						TagInteger, Integer(14852)),
					MakeAttribute("y-dimension",
						TagInteger, Integer(20990)),
				},
			),
			indent: 2,
			out: []string{
				`  ATTR "media-size-supported" collection: {`,
				`      MEMBER "x-dimension" integer: 20990`,
				`      MEMBER "y-dimension" integer: 29704`,
				`  }`,
				`  {`,
				`      MEMBER "x-dimension" integer: 14852`,
				`      MEMBER "y-dimension" integer: 20990`,
				`  }`,
			},
		},

		// Multiple values
		{
			attr: MakeAttr("page-delivery-supported",
				TagKeyword,
				String("reverse-order"),
				String("same-order")),

			out: []string{
				`ATTR "page-delivery-supported" keyword: reverse-order same-order`,
			},
		},

		// Values of mixed type
		{
			attr: Attribute{
				Name: "page-ranges",
				Values: Values{
					{TagInteger, Integer(1)},
					{TagInteger, Integer(2)},
					{TagInteger, Integer(3)},
					{TagRange, Range{5, 7}},
				},
			},

			out: []string{
				`ATTR "page-ranges" integer: 1 2 3 rangeOfInteger: 5-7`,
			},
		},
	}

	f := NewFormatter()
	for _, test := range tests {
		f.Reset()
		f.SetIndent(test.indent)

		expected := strings.Join(test.out, "\n") + "\n"
		f.FmtAttribute(test.attr)
		out := f.String()

		if out != expected {
			t.Errorf("output mismatch\n"+
				"expected:\n%s"+
				"present:\n%s",
				expected, out)
		}
	}
}

// TestFmtRequestResponse runs Formatter.FmtRequest and
// Formatter.FmtResponse tests
func TestFmtRequestResponse(t *testing.T) {
	type testData struct {
		msg *Message // Input message
		rq  bool     // This is request
		out []string // Expected output
	}

	tests := []testData{
		{
			msg: &Message{
				Version:   MakeVersion(2, 0),
				Code:      Code(OpGetPrinterAttributes),
				RequestID: 1,

				Operation: []Attribute{
					MakeAttribute(
						"attributes-charset",
						TagCharset,
						String("utf-8")),
					MakeAttribute(
						"attributes-natural-language",
						TagLanguage,
						String("en-us")),
					MakeAttribute(
						"requested-attributes",
						TagKeyword,
						String("printer-name")),
				},
			},
			rq: true,
			out: []string{
				`{`,
				`    REQUEST-ID 1`,
				`    VERSION 2.0`,
				`    OPERATION Get-Printer-Attributes`,
				``,
				`    GROUP operation-attributes-tag`,
				`    ATTR "attributes-charset" charset: utf-8`,
				`    ATTR "attributes-natural-language" naturalLanguage: en-us`,
				`    ATTR "requested-attributes" keyword: printer-name`,
				`}`,
			},
		},

		{
			msg: &Message{
				Version:   MakeVersion(2, 0),
				Code:      Code(StatusOk),
				RequestID: 1,

				Operation: []Attribute{
					MakeAttribute(
						"attributes-charset",
						TagCharset,
						String("utf-8")),
					MakeAttribute(
						"attributes-natural-language",
						TagLanguage,
						String("en-us")),
				},

				Printer: []Attribute{
					MakeAttribute(
						"printer-name",
						TagName,
						String("Kyocera_ECOSYS_M2040dn")),
				},
			},
			rq: false,
			out: []string{
				`{`,
				`    REQUEST-ID 1`,
				`    VERSION 2.0`,
				`    STATUS successful-ok`,
				``,
				`    GROUP operation-attributes-tag`,
				`    ATTR "attributes-charset" charset: utf-8`,
				`    ATTR "attributes-natural-language" naturalLanguage: en-us`,
				``,
				`    GROUP printer-attributes-tag`,
				`    ATTR "printer-name" nameWithoutLanguage: Kyocera_ECOSYS_M2040dn`,
				`}`,
			},
		},
	}

	f := NewFormatter()
	for _, test := range tests {
		f.Reset()
		if test.rq {
			f.FmtRequest(test.msg)
		} else {
			f.FmtResponse(test.msg)
		}

		out := f.String()
		expected := strings.Join(test.out, "\n") + "\n"

		if out != expected {
			t.Errorf("output mismatch\n"+
				"expected:\n%s"+
				"present:\n%s",
				expected, out)
		}
	}
}
