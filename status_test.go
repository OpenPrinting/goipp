/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP Status Codes tests
 */

package goipp

import (
	"fmt"
	"testing"
)

// TestStatusString tests Status.String method
func TestStatusString(t *testing.T) {
	type testData struct {
		status Status // Input Op code
		s      string // Expected output string
	}

	tests := []testData{
		{StatusOk, "successful-ok"},
		{StatusOkConflicting, "successful-ok-conflicting-attributes"},
		{StatusOkEventsComplete, "successful-ok-events-complete"},
		{StatusRedirectionOtherSite, "redirection-other-site"},
		{StatusErrorBadRequest, "client-error-bad-request"},
		{StatusErrorForbidden, "client-error-forbidden"},
		{StatusErrorNotFetchable, "client-error-not-fetchable"},
		{StatusErrorInternal, "server-error-internal-error"},
		{StatusErrorTooManyDocuments, "server-error-too-many-documents"},
		{0xabcd, "0xabcd"},
	}

	for _, test := range tests {
		s := test.status.String()
		if s != test.s {
			t.Errorf("testing Status.String:\n"+
				"input:    0x%4.4x\n"+
				"expected: %s\n"+
				"present:  %s\n",
				int(test.status), test.s, s,
			)
		}
	}
}

// TestStatusGoString tests Status.GoString method
func TestStatusGoString(t *testing.T) {
	type testData struct {
		status Status // Input Op code
		s      string // Expected output string
	}

	tests := []testData{
		{StatusOk, "goipp.StatusOk"},
		{StatusOkConflicting, "goipp.StatusOkConflicting"},
		{StatusOkEventsComplete, "goipp.StatusOkEventsComplete"},
		{StatusRedirectionOtherSite, "goipp.StatusRedirectionOtherSite"},
		{StatusErrorBadRequest, "goipp.StatusErrorBadRequest"},
		{StatusErrorForbidden, "goipp.StatusErrorForbidden"},
		{StatusErrorNotFetchable, "goipp.StatusErrorNotFetchable"},
		{StatusErrorInternal, "goipp.StatusErrorInternal"},
		{StatusErrorTooManyDocuments, "goipp.StatusErrorTooManyDocuments"},
		{0xabcd, "goipp.Status(0xabcd)"},
	}

	for _, test := range tests {
		s := fmt.Sprintf("%#v", test.status)
		if s != test.s {
			t.Errorf("testing Status.GoString:\n"+
				"input:    0x%4.4x\n"+
				"expected: %s\n"+
				"present:  %s\n",
				int(test.status), test.s, s,
			)
		}
	}
}
