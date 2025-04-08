/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP Message tests
 */

package goipp

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

// TestVersion tests Version functions
func TestVersion(t *testing.T) {
	type testData struct {
		major, minor uint8   // Major/minor parts
		ver          Version // Resulting version
		str          string  // Its string representation
	}

	tests := []testData{
		{
			major: 2,
			minor: 0,
			ver:   0x0200,
			str:   "2.0",
		},

		{
			major: 2,
			minor: 1,
			ver:   0x0201,
			str:   "2.1",
		},

		{
			major: 1,
			minor: 1,
			ver:   0x0101,
			str:   "1.1",
		},
	}

	for _, test := range tests {
		ver := MakeVersion(test.major, test.minor)
		if ver != test.ver {
			t.Errorf("MakeVersion test failed:\n"+
				"version expected: 0x%4.4x\n"+
				"version present:  0x%4.4x\n",
				uint(test.ver), uint(ver),
			)
			continue
		}

		str := ver.String()
		if str != test.str {
			t.Errorf("Version.String test failed:\n"+
				"expected: %s\n"+
				"present:  %s\n",
				test.str, str,
			)
			continue
		}

		major := ver.Major()
		if major != test.major {
			t.Errorf("Version.Major test failed:\n"+
				"expected: %d\n"+
				"present:  %d\n",
				test.major, major,
			)
			continue
		}

		minor := ver.Minor()
		if minor != test.minor {
			t.Errorf("Version.Minor test failed:\n"+
				"expected: %d\n"+
				"present:  %d\n",
				test.minor, minor,
			)
			continue
		}
	}
}

// TestNewRequestResponse tests NewRequest and NewResponse functions
func TestNewRequestResponse(t *testing.T) {
	msg := &Message{
		Version:   MakeVersion(2, 0),
		Code:      1,
		RequestID: 0x12345,
	}

	rq := NewRequest(msg.Version, Op(msg.Code), msg.RequestID)
	if !reflect.DeepEqual(msg, rq) {
		t.Errorf("NewRequest test failed:\n"+
			"expected: %#v\n"+
			"present:  %#v\n",
			msg, rq,
		)
	}

	rsp := NewResponse(msg.Version, Status(msg.Code), msg.RequestID)
	if !reflect.DeepEqual(msg, rsp) {
		t.Errorf("NewRequest test failed:\n"+
			"expected: %#v\n"+
			"present:  %#v\n",
			msg, rsp,
		)
	}
}

// TestNewMessageWithGroups tests the NewMessageWithGroups function.
func TestNewMessageWithGroups(t *testing.T) {
	// Populate groups
	ops := Group{
		TagOperationGroup,
		Attributes{
			MakeAttr("ops", TagInteger, Integer(1)),
		},
	}

	prn1 := Group{
		TagPrinterGroup,
		Attributes{
			MakeAttr("prn1", TagInteger, Integer(2)),
		},
	}

	prn2 := Group{
		TagPrinterGroup,
		Attributes{
			MakeAttr("prn2", TagInteger, Integer(3)),
		},
	}

	prn3 := Group{
		TagPrinterGroup,
		Attributes{
			MakeAttr("prn3", TagInteger, Integer(4)),
		},
	}

	job := Group{
		TagJobGroup,
		Attributes{
			MakeAttr("job", TagInteger, Integer(5)),
		},
	}

	unsupp := Group{
		TagUnsupportedGroup,
		Attributes{
			MakeAttr("unsupp", TagInteger, Integer(6)),
		},
	}

	sub := Group{
		TagSubscriptionGroup,
		Attributes{
			MakeAttr("sub", TagInteger, Integer(7)),
		},
	}

	evnt := Group{
		TagEventNotificationGroup,
		Attributes{
			MakeAttr("evnt", TagInteger, Integer(8)),
		},
	}

	res := Group{
		TagResourceGroup,
		Attributes{
			MakeAttr("res", TagInteger, Integer(9)),
		},
	}

	doc := Group{
		TagDocumentGroup,
		Attributes{
			MakeAttr("doc", TagInteger, Integer(10)),
		},
	}

	sys := Group{
		TagSystemGroup,
		Attributes{
			MakeAttr("sys", TagInteger, Integer(11)),
		},
	}

	future11 := Group{
		TagFuture11Group,
		Attributes{
			MakeAttr("future11", TagInteger, Integer(12)),
		},
	}

	future12 := Group{
		TagFuture12Group,
		Attributes{
			MakeAttr("future12", TagInteger, Integer(13)),
		},
	}

	future13 := Group{
		TagFuture13Group,
		Attributes{
			MakeAttr("future13", TagInteger, Integer(14)),
		},
	}

	future14 := Group{
		TagFuture14Group,
		Attributes{
			MakeAttr("future14", TagInteger, Integer(15)),
		},
	}

	future15 := Group{
		TagFuture15Group,
		Attributes{
			MakeAttr("future15", TagInteger, Integer(16)),
		},
	}

	groups := Groups{
		ops,
		prn1, prn2, prn3,
		job,
		unsupp,
		sub,
		evnt,
		res,
		doc,
		sys,
		future11,
		future12,
		future13,
		future14,
		future15,
	}

	msg := NewMessageWithGroups(DefaultVersion, 1, 123, groups)
	expected := &Message{
		Version:           DefaultVersion,
		Code:              1,
		RequestID:         123,
		Groups:            groups,
		Operation:         ops.Attrs,
		Job:               job.Attrs,
		Unsupported:       unsupp.Attrs,
		Subscription:      sub.Attrs,
		EventNotification: evnt.Attrs,
		Resource:          res.Attrs,
		Document:          doc.Attrs,
		System:            sys.Attrs,
		Future11:          future11.Attrs,
		Future12:          future12.Attrs,
		Future13:          future13.Attrs,
		Future14:          future14.Attrs,
		Future15:          future15.Attrs,
	}
	expected.Printer = prn1.Attrs
	expected.Printer = append(expected.Printer, prn2.Attrs...)
	expected.Printer = append(expected.Printer, prn3.Attrs...)

	if !reflect.DeepEqual(msg, expected) {
		t.Errorf("NewMessageWithGroups test failed:\n"+
			"expected: %#v\n"+
			"present:  %#v\n",
			expected,
			msg,
		)
	}
}

// TestNewMessageWithGroups tests the Message.AttrGroups function.
func TestMessageAttrGroups(t *testing.T) {
	// Create a message for testing
	uri := "ipp://192/168.0.1/ipp/print"

	m := NewRequest(DefaultVersion, OpCreateJob, 1)

	m.Operation.Add(MakeAttr("attributes-charset",
		TagCharset, String("utf-8")))
	m.Operation.Add(MakeAttr("attributes-natural-language",
		TagLanguage, String("en-US")))
	m.Operation.Add(MakeAttr("printer-uri",
		TagURI, String(uri)))

	m.Job.Add(MakeAttr("copies", TagInteger, Integer(1)))

	// Compare m.AttrGroups() with expectations
	groups := m.AttrGroups()
	expected := Groups{
		Group{
			Tag: TagOperationGroup,
			Attrs: Attributes{
				MakeAttr("attributes-charset",
					TagCharset, String("utf-8")),
				MakeAttr("attributes-natural-language",
					TagLanguage, String("en-US")),
				MakeAttr("printer-uri",
					TagURI, String(uri)),
			},
		},
		Group{
			Tag: TagJobGroup,
			Attrs: Attributes{
				MakeAttr("copies", TagInteger, Integer(1)),
			},
		},
	}

	if !reflect.DeepEqual(groups, expected) {
		t.Errorf("Message.AttrGroups test failed:\n"+
			"expected: %#v\n"+
			"present:  %#v\n",
			expected, groups,
		)
	}

	// Set m.Groups. Check that it takes precedence.
	expected = Groups{
		Group{
			Tag: TagOperationGroup,
			Attrs: Attributes{
				MakeAttr("attributes-charset",
					TagCharset, String("utf-8")),
			},
		},
	}

	m.Groups = expected
	groups = m.AttrGroups()

	if !reflect.DeepEqual(groups, expected) {
		t.Errorf("Message.AttrGroups test failed:\n"+
			"expected: %#v\n"+
			"present:  %#v\n",
			expected, groups,
		)
	}
}

// TestMessageEqualSimilar tests Message.Equal and Message.Similar functions.
func TestMessageEqualSimilar(t *testing.T) {
	type testData struct {
		m1, m2  Message // Input messages
		equal   bool    // Expected Message.Equal output
		similar bool    // Expected Message.Similar output
	}

	uri := "ipp://192/168.0.1/ipp/print"

	tests := []testData{
		// Empty messages are equal and similar
		{
			m1:      Message{},
			m2:      Message{},
			equal:   true,
			similar: true,
		},

		// Messages with different Version/Code/RequestID are
		// neither equal or similar
		{
			m1:      Message{},
			m2:      Message{Version: 1},
			equal:   false,
			similar: false,
		},

		{
			m1:      Message{},
			m2:      Message{Code: 1},
			equal:   false,
			similar: false,
		},

		{
			m1:      Message{},
			m2:      Message{RequestID: 1},
			equal:   false,
			similar: false,
		},

		// If the same attributes represented as Message.Groups in one
		// message and via Message.Operation/Job/Printer etc in the
		// another message, these messages are equal and similar
		{
			m1: Message{
				Groups: Groups{
					Group{
						Tag: TagOperationGroup,
						Attrs: Attributes{
							MakeAttr("attributes-charset",
								TagCharset, String("utf-8")),
							MakeAttr("attributes-natural-language",
								TagLanguage, String("en-US")),
							MakeAttr("printer-uri",
								TagURI, String(uri)),
						},
					},
					Group{
						Tag: TagJobGroup,
						Attrs: Attributes{
							MakeAttr("copies", TagInteger, Integer(1)),
						},
					},
				},
			},

			m2: Message{
				Operation: Attributes{
					MakeAttr("attributes-charset",
						TagCharset, String("utf-8")),
					MakeAttr("attributes-natural-language",
						TagLanguage, String("en-US")),
					MakeAttr("printer-uri",
						TagURI, String(uri)),
				},

				Job: Attributes{
					MakeAttr("copies", TagInteger, Integer(1)),
				},
			},

			equal:   true,
			similar: true,
		},

		// Messages with the different order of the same set of attributes
		// are similar but not equal.
		{
			m1: Message{
				Operation: Attributes{
					MakeAttr("attributes-charset",
						TagCharset, String("utf-8")),
					MakeAttr("attributes-natural-language",
						TagLanguage, String("en-US")),
					MakeAttr("printer-uri",
						TagURI, String(uri)),
				},
			},

			m2: Message{
				Operation: Attributes{
					MakeAttr("attributes-charset",
						TagCharset, String("utf-8")),
					MakeAttr("printer-uri",
						TagURI, String(uri)),
					MakeAttr("attributes-natural-language",
						TagLanguage, String("en-US")),
				},
			},

			equal:   false,
			similar: true,
		},
	}

	for _, test := range tests {
		equal := test.m1.Equal(test.m2)
		if equal != test.equal {
			var buf1, buf2 bytes.Buffer
			test.m1.Print(&buf1, true)
			test.m2.Print(&buf2, true)

			t.Errorf("testing Message.Equal:\n"+
				"message 1: %s\n"+
				"message 2: %s\n"+
				"expected:  %v\n"+
				"present:   %v\n",
				&buf1, &buf2,
				test.equal, equal,
			)
		}

		similar := test.m1.Similar(test.m2)
		if similar != test.similar {
			var buf1, buf2 bytes.Buffer
			test.m1.Print(&buf1, true)
			test.m2.Print(&buf2, true)

			t.Errorf("testing Message.Similar:\n"+
				"message 1: %s\n"+
				"message 2: %s\n"+
				"expected:  %v\n"+
				"present:   %v\n",
				&buf1, &buf2,
				test.similar, similar,
			)
		}
	}
}

// TestMessageReset tests Message.Reset function
func TestMessageReset(t *testing.T) {
	uri := "ipp://192/168.0.1/ipp/print"
	m := Message{
		Groups: Groups{
			Group{
				Tag: TagOperationGroup,
				Attrs: Attributes{
					MakeAttr("attributes-charset",
						TagCharset, String("utf-8")),
					MakeAttr("attributes-natural-language",
						TagLanguage, String("en-US")),
					MakeAttr("printer-uri",
						TagURI, String(uri)),
				},
			},
			Group{
				Tag: TagJobGroup,
				Attrs: Attributes{
					MakeAttr("copies", TagInteger, Integer(1)),
				},
			},
		},
	}

	m.Reset()

	if !reflect.ValueOf(m).IsZero() {
		t.Errorf("Message.Reset test failed")
	}
}

// TestMessagePrint tests Message.Print function
func TestMessagePrint(t *testing.T) {
	uri := "ipp://192/168.0.1/ipp/print"
	m := Message{
		Code:      2,
		Version:   MakeVersion(2, 0),
		RequestID: 1,
		Groups: Groups{
			Group{
				Tag: TagOperationGroup,
				Attrs: Attributes{
					MakeAttr("attributes-charset",
						TagCharset, String("utf-8")),
					MakeAttr("attributes-natural-language",
						TagLanguage, String("en-US")),
					MakeAttr("printer-uri",
						TagURI, String(uri)),
				},
			},
			Group{
				Tag: TagJobGroup,
				Attrs: Attributes{
					MakeAttr("copies", TagInteger, Integer(1)),
				},
			},
		},
	}

	// Check request formatting
	reqExpected := []string{
		`{`,
		`    REQUEST-ID 1`,
		`    VERSION 2.0`,
		`    OPERATION Print-Job`,
		``,
		`    GROUP operation-attributes-tag`,
		`    ATTR "attributes-charset" charset: utf-8`,
		`    ATTR "attributes-natural-language" naturalLanguage: en-US`,
		`    ATTR "printer-uri" uri: ipp://192/168.0.1/ipp/print`,
		``,
		`    GROUP job-attributes-tag`,
		`    ATTR "copies" integer: 1`,
		`}`,
	}

	var buf bytes.Buffer
	m.Print(&buf, true)
	exp := strings.Join(reqExpected, "\n") + "\n"

	if buf.String() != exp {
		t.Errorf("Message.Print test failed for request:\n"+
			"expected: %s\n"+
			"present:  %s\n",
			exp, &buf,
		)
	}

	// Check response formatting
	rspExpected := []string{
		`{`,
		`    REQUEST-ID 1`,
		`    VERSION 2.0`,
		`    STATUS successful-ok`,
		``,
		`    GROUP operation-attributes-tag`,
		`    ATTR "attributes-charset" charset: utf-8`,
		`    ATTR "attributes-natural-language" naturalLanguage: en-US`,
		`    ATTR "printer-uri" uri: ipp://192/168.0.1/ipp/print`,
		``,
		`    GROUP job-attributes-tag`,
		`    ATTR "copies" integer: 1`,
		`}`,
	}

	buf.Reset()
	m.Code = 0
	m.Print(&buf, false)
	exp = strings.Join(rspExpected, "\n") + "\n"

	if buf.String() != exp {
		t.Errorf("Message.Print test failed for response:\n"+
			"expected: %s\n"+
			"present:  %s\n",
			exp, &buf,
		)
	}
}
