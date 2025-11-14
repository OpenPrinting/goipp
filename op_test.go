/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP Operation Codes tests
 */

package goipp

import (
	"fmt"
	"testing"
)

// TestOpString tests Op.String method
func TestOpString(t *testing.T) {
	type testData struct {
		op Op     // Input Op code
		s  string // Expected output string
	}

	tests := []testData{
		{OpPrintJob, "Print-Job"},
		{OpPrintURI, "Print-URI"},
		{OpPausePrinter, "Pause-Printer"},
		{OpRestartSystem, "Restart-System"},
		{OpCupsGetDefault, "CUPS-Get-Default"},
		{OpCupsGetPpd, "CUPS-Get-PPD"},
		{OpCupsCreateLocalPrinter, "CUPS-Create-Local-Printer"},
		{0xabcd, "0xabcd"},
	}

	for _, test := range tests {
		s := test.op.String()
		if s != test.s {
			t.Errorf("testing Op.String:\n"+
				"input:    0x%4.4x\n"+
				"expected: %s\n"+
				"present:  %s\n",
				int(test.op), test.s, s,
			)
		}
	}
}

// TestOpGoString tests Op.GoString method
func TestOpGoString(t *testing.T) {
	type testData struct {
		op Op     // Input Op code
		s  string // Expected output string
	}

	tests := []testData{
		{OpPrintJob, "goipp.OpPrintJob"},
		{OpPrintURI, "goipp.OpPrintURI"},
		{OpPausePrinter, "goipp.OpPausePrinter"},
		{OpRestartSystem, "goipp.OpRestartSystem"},
		{OpCupsGetDefault, "goipp.OpCupsGetDefault"},
		{OpCupsGetPpd, "goipp.OpCupsGetPpd"},
		{OpCupsCreateLocalPrinter, "goipp.OpCupsCreateLocalPrinter"},
		{0xabcd, "goipp.Op(0xabcd)"},
	}

	for _, test := range tests {
		s := fmt.Sprintf("%#v", test.op)
		if s != test.s {
			t.Errorf("testing Op.String:\n"+
				"input:    0x%4.4x\n"+
				"expected: %s\n"+
				"present:  %s\n",
				int(test.op), test.s, s,
			)
		}
	}
}
