/* Go IPP - IPP core protocol implementation in pure Go
 *
 * Copyright (C) 2020 and up by Alexander Pevzner (pzz@apevzner.com)
 * See LICENSE for license terms and conditions
 *
 * IPP Status Codes tests
 */

package goipp

import "testing"

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
