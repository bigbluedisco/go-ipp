package ipp

import (
	"bytes"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

var requestTestCases = []struct {
	Request Request
	Bytes   []byte
	// some test cases are only for encoding, skip them in decoding tests
	SkipDecoding bool
}{
	{
		Request: Request{
			ProtocolVersionMajor: ProtocolVersionMajor,
			ProtocolVersionMinor: ProtocolVersionMinor,
			Operation:            OperationPrintJob,
			RequestId:            12345,
		},
		Bytes:        []byte{2, 0, 0, 2, 0, 0, 48, 57, 1, 71, 0, 18, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 99, 104, 97, 114, 115, 101, 116, 0, 5, 117, 116, 102, 45, 56, 72, 0, 27, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 110, 97, 116, 117, 114, 97, 108, 45, 108, 97, 110, 103, 117, 97, 103, 101, 0, 5, 101, 110, 45, 85, 83, 3},
		SkipDecoding: true,
	},
	{
		Request: Request{
			ProtocolVersionMajor: ProtocolVersionMajor,
			ProtocolVersionMinor: ProtocolVersionMinor,
			Operation:            OperationPrintJob,
			RequestId:            12345,
			OperationAttributes: map[string]interface{}{
				AttributeCharset:         Charset,
				AttributeNaturalLanguage: CharsetLanguage,
			},
		},
		Bytes:        []byte{2, 0, 0, 2, 0, 0, 48, 57, 1, 71, 0, 18, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 99, 104, 97, 114, 115, 101, 116, 0, 5, 117, 116, 102, 45, 56, 72, 0, 27, 97, 116, 116, 114, 105, 98, 117, 116, 101, 115, 45, 110, 97, 116, 117, 114, 97, 108, 45, 108, 97, 110, 103, 117, 97, 103, 101, 0, 5, 101, 110, 45, 85, 83, 3},
		SkipDecoding: true,
	},
}

func TestRequest_Encode(t *testing.T) {
	for _, c := range requestTestCases {
		data, err := c.Request.Encode()
		assert.Nil(t, err)
		assert.Equal(t, c.Bytes, data, "encoded request is not correct")
	}
}

func TestRequestDecoder_Decode(t *testing.T) {
	for _, c := range requestTestCases {
		if c.SkipDecoding {
			continue
		}

		request, err := NewRequestDecoder(bytes.NewReader(c.Bytes)).Decode(nil)
		assert.Nil(t, err)
		assert.Equal(t, &c.Request, request, "decoded request is not correct")
	}
}

func TestRequest_EncodeCollationJobAttributes(t *testing.T) {
	req := NewRequest(OperationPrintJob, 12345)
	req.JobAttributes[AttributeMultipleDocumentHandling] = "separate-documents-collated-copies"
	req.JobAttributes[AttributeSheetCollate] = "collated"

	data, err := req.Encode()
	require.NoError(t, err)

	// job attributes are encoded in map order, so check each attribute on its own
	// instead of the whole request. 0x44 is the keyword value tag.
	for _, attr := range []string{
		"\x44\x00\x1amultiple-document-handling\x00\x22separate-documents-collated-copies",
		"\x44\x00\x0dsheet-collate\x00\x08collated",
	} {
		assert.Contains(t, string(data), attr, "encoded request is not correct")
	}

	request, err := NewRequestDecoder(bytes.NewReader(data)).Decode(nil)
	require.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		AttributeMultipleDocumentHandling: "separate-documents-collated-copies",
		AttributeSheetCollate:             "collated",
	}, request.JobAttributes, "decoded request is not correct")
}
