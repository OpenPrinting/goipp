/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * Tests for enumeration of value types
 */

package goipp

import "testing"

// TestTypeString tests Type.String function
func TestTypeString(t *testing.T) {
	type testData struct {
		ty Type   // Input Type
		s  string // Expected output string
	}

	tests := []testData{
		{TypeInvalid, "Invalid"},
		{TypeVoid, "Void"},
		{TypeBoolean, "Boolean"},
		{TypeString, "String"},
		{TypeDateTime, "DateTime"},
		{TypeResolution, "Resolution"},
		{TypeRange, "Range"},
		{TypeTextWithLang, "TextWithLang"},
		{TypeBinary, "Binary"},
		{TypeCollection, "Collection"},
		{0x1234, "0x1234"},
	}

	for _, test := range tests {
		s := test.ty.String()
		if s != test.s {
			t.Errorf("testing Type.String:\n"+
				"input:    %d\n"+
				"expected: %s\n"+
				"present:  %s\n",
				int(test.ty), test.s, s,
			)
		}
	}
}
