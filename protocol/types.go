package protocol

const (
	CHECKSUM = "checksum"
	CONTENT = "content"
	ERROR = "error"
	SEGMENT = "segment"

	CHECKSUM_SIZE = 20
)

type FileTransferError struct {
	s string
}

func (e FileTransferError) Error() string {
	return e.s
}
