package ipp

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestResponseDecoder_UnsupportedAttributes(t *testing.T) {
	// Exact printer response captured on 2026-09-14: successful response
	// with an unsupported requested attribute after the printer group.
	const encoded = "AgAAAQAAAAEBRwASYXR0cmlidXRlcy1jaGFyc2V0AAV1dGYtOEgAG2F0dHJpYnV0ZXMtbmF0dXJhbC1sYW5ndWFnZQAFZW4tdXMEIwANcHJpbnRlci1zdGF0ZQAEAAAAA0kAGWRvY3VtZW50LWZvcm1hdC1zdXBwb3J0ZWQAGGFwcGxpY2F0aW9uL29jdGV0LXN0cmVhbUkAAAAPYXBwbGljYXRpb24vcGRmSQAAAAlpbWFnZS91cmZJAAAACmltYWdlL2pwZWczABBjb3BpZXMtc3VwcG9ydGVkAAgAAAABAAAD50QADW1lZGlhLWRlZmF1bHQAEGlzb19hNF8yMTB4Mjk3bW0FRAAUcmVxdWVzdGVkLWF0dHJpYnV0ZXMAKHB3Zy1yYXN0ZXItZG9jdW1lbnQtcmVzb2x1dGlvbi1zdXBwb3J0ZWQD"
	body, err := base64.StdEncoding.DecodeString(encoded)
	require.NoError(t, err)
	// Also exercise the standard ordering, with unsupported attributes
	// between the operation and printer groups.
	printerOffset := bytes.Index(body, []byte{byte(TagPrinter), byte(TagEnum), 0, 13})
	require.NotEqual(t, -1, printerOffset)
	beforePrinter := append([]byte{}, body[:printerOffset]...)
	beforePrinter = append(beforePrinter, body[263:329]...)
	beforePrinter = append(beforePrinter, body[printerOffset:263]...)
	beforePrinter = append(beforePrinter, byte(TagEnd))
	for name, payload := range map[string][]byte{"captured": body, "before printer": beforePrinter} {
		t.Run(name, func(t *testing.T) {
			r := bytes.NewReader(payload)
			resp, err := NewResponseDecoder(r).Decode(nil)
			require.NoError(t, err)
			assert.Zero(t, r.Len())
			assert.Equal(t, int16(1), resp.StatusCode)
			assert.Equal(t, int32(1), resp.RequestId)
			assert.Equal(t, "utf-8", resp.OperationAttributes[AttributeCharset][0].Value)
			require.Len(t, resp.PrinterAttributes, 1)
			attrs := resp.PrinterAttributes[0]
			assert.Equal(t, 3, attrs["printer-state"][0].Value)
			assert.Equal(t, []int32{1, 999}, attrs["copies-supported"][0].Value)
			assert.Equal(t, "iso_a4_210x297mm", attrs["media-default"][0].Value)
			require.Len(t, resp.UnsupportedAttributes["requested-attributes"], 1)
			assert.Equal(t, "pwg-raster-document-resolution-supported", resp.UnsupportedAttributes["requested-attributes"][0].Value)
			assert.NotContains(t, attrs, "requested-attributes")
		})
	}
}
