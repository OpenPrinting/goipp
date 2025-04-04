/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP Message tests
 */

package goipp

import (
	"reflect"
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
