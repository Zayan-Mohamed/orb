package protocol

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
)

// Protocol constants
const (
	MaxFrameSize = 1 << 20 // 1 MB max frame size
	HeaderSize   = 8       // 4 bytes length + 4 bytes type
)

// Frame types
const (
	FrameTypeHandshake     = 0x01
	FrameTypeHandshakeResp = 0x02
	FrameTypeList          = 0x10
	FrameTypeStat          = 0x11
	FrameTypeRead          = 0x12
	FrameTypeWrite         = 0x13
	FrameTypeDelete        = 0x14
	FrameTypeRename        = 0x15
	FrameTypeMkdir         = 0x16
	FrameTypeResponse      = 0x20
	FrameTypeError         = 0x21
	FrameTypePing          = 0x30
	FrameTypePong          = 0x31
)

var (
	ErrFrameTooLarge    = errors.New("frame exceeds maximum size")
	ErrInvalidFrame     = errors.New("invalid frame format")
	ErrUnknownFrameType = errors.New("unknown frame type")
)

// Frame represents a protocol frame
type Frame struct {
	Type    uint32
	Payload []byte
}

// WriteFrame writes a frame to the writer
// Format: [4-byte length][4-byte type][encrypted payload]
func WriteFrame(w io.Writer, frame *Frame) error {
	if len(frame.Payload) > MaxFrameSize {
		return ErrFrameTooLarge
	}

	// Write length with overflow check
	payloadLen := len(frame.Payload)
	if payloadLen < 0 || payloadLen > MaxFrameSize {
		return fmt.Errorf("invalid payload length: %d", payloadLen)
	}
	length := uint32(payloadLen) // #nosec G115 -- length is validated above
	if err := binary.Write(w, binary.BigEndian, length); err != nil {
		return fmt.Errorf("failed to write length: %w", err)
	}

	// Write type
	if err := binary.Write(w, binary.BigEndian, frame.Type); err != nil {
		return fmt.Errorf("failed to write type: %w", err)
	}

	// Write payload
	if _, err := w.Write(frame.Payload); err != nil {
		return fmt.Errorf("failed to write payload: %w", err)
	}

	return nil
}

// ReadFrame reads a frame from the reader
func ReadFrame(r io.Reader) (*Frame, error) {
	// Read length
	var length uint32
	if err := binary.Read(r, binary.BigEndian, &length); err != nil {
		return nil, fmt.Errorf("failed to read length: %w", err)
	}

	// Validate length
	if length > MaxFrameSize {
		return nil, ErrFrameTooLarge
	}

	// Read type
	var frameType uint32
	if err := binary.Read(r, binary.BigEndian, &frameType); err != nil {
		return nil, fmt.Errorf("failed to read type: %w", err)
	}

	// Read payload
	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return nil, fmt.Errorf("failed to read payload: %w", err)
	}

	return &Frame{
		Type:    frameType,
		Payload: payload,
	}, nil
}

// ValidateFrameType checks if a frame type is valid
func ValidateFrameType(frameType uint32) bool {
	validTypes := map[uint32]bool{
		FrameTypeHandshake:     true,
		FrameTypeHandshakeResp: true,
		FrameTypeList:          true,
		FrameTypeStat:          true,
		FrameTypeRead:          true,
		FrameTypeWrite:         true,
		FrameTypeDelete:        true,
		FrameTypeRename:        true,
		FrameTypeMkdir:         true,
		FrameTypeResponse:      true,
		FrameTypeError:         true,
		FrameTypePing:          true,
		FrameTypePong:          true,
	}
	return validTypes[frameType]
}

// Request types for filesystem operations
type ListRequest struct {
	Path string
}

type StatRequest struct {
	Path string
}

type ReadRequest struct {
	Path   string
	Offset int64
	Length int64
}

type WriteRequest struct {
	Path   string
	Offset int64
	Data   []byte
}

type DeleteRequest struct {
	Path string
}

type RenameRequest struct {
	OldPath string
	NewPath string
}

type MkdirRequest struct {
	Path string
	Perm uint32
}

// Response types
type FileInfo struct {
	Name    string
	Size    int64
	Mode    uint32
	ModTime int64
	IsDir   bool
}

type ListResponse struct {
	Files []FileInfo
}

type StatResponse struct {
	Info FileInfo
}

type ReadResponse struct {
	Data []byte
}

type WriteResponse struct {
	BytesWritten int64
}

type ErrorResponse struct {
	Code    uint32
	Message string
}

// Error codes
const (
	ErrCodeNotFound      = 1
	ErrCodePermission    = 2
	ErrCodeExists        = 3
	ErrCodeIsDirectory   = 4
	ErrCodeNotDirectory  = 5
	ErrCodeInvalidPath   = 6
	ErrCodeQuotaExceeded = 7
	ErrCodeIO            = 8
	ErrCodeUnknown       = 99
)
