package protocol

import (
	"bytes"
	"crypto/sha1"
	"fmt"
	"strconv"
	"strings"
)

/* RESPONSE FORMAT
*
* error:[error message (nothing in case of no error)]\n
* segment:[segment number of the file being streamed (-1 in case of error)]\n
* content:[little endian content encoded in base64]\n
* checksum:[checksum of all the other fields appended in order]
*
 */
type Response struct {
	Content string
	Error error
	SegmentNumber int
	Body []byte
	checksum [20]byte
}


func (m *Response) Encode() []byte {

	s := ""
	s += m.formatField(ERROR, m.Error) + "\n"
	s += m.formatField(SEGMENT, m.SegmentNumber) + "\n"
	s += m.formatField(CONTENT, m.Content) + "\n"
	
	m.checksum = calculateChecksum([]byte(s))
	// stupid & inefficient hack to make sure formatting works properly
	s += m.formatField(CHECKSUM, string(m.checksum[:]))

	m.Body = []byte(s)
	return m.Body
}

/*
* this assumes the field name is not an empty string
*/
func (m *Response) formatField(fieldName string, value any) string {
	fieldStr := fieldName + ":"
	switch v := value.(type) {
		case int:
			valueString := strconv.Itoa(v)
			fieldStr += valueString

		case string:
			fieldStr += v

		case error:
			fieldStr += v.Error()

	}

	return fieldStr
}

func (m *Response) Decode(data []byte) error {
	m.Body = data
	bodyStr := string(m.Body)
	fieldValuePair := make([]string, 2)
	fields := strings.SplitSeq(bodyStr, "\n")

	for s := range fields {
		fieldValuePair = strings.SplitN(s, ":", 2)
		field := fieldValuePair[0]
		value := fieldValuePair[1]

		switch field {
			case CHECKSUM:
				copy(m.checksum[:], value)

			case CONTENT:
				m.Content = value

			case ERROR:
				m.Error = FileTransferError{s: value}

			case SEGMENT:
				segment, err := strconv.Atoi(value)
				if err != nil {
					return fmt.Errorf("malformed field %s", s)
				}
				m.SegmentNumber = segment
				
			default:
				return fmt.Errorf("invalid field %s", field)
		}
	}
	
	return nil
}

func (m *Response) ValidateChecksum() (valid bool) {
	// index of the beginning of checksum field
	i := + len(m.Body) - (len(CHECKSUM) + 1 + CHECKSUM_SIZE)
	b := m.Body[:i]

	checksum := calculateChecksum(b)
	return bytes.Equal(checksum[:], m.checksum[:])
}

func calculateChecksum(b []byte) [20]byte {
	return sha1.Sum(b)
}
