package packet

import "errors"

const HeaderSize = PublicHeaderSize + PrivateHeaderSize

var ErrHeaderSizeMismatch = errors.New("decoded offset don't match header size")

var ErrSmallBufferSize = errors.New("buffer size to small")

type Header struct {
	PublicHeader
	PrivateHeader
}

func (h *Header) Encode(buf []byte) (int, error) {
	var totalOffset int
	n, err := h.PublicHeader.Encode(buf)
	if err != nil {
		return 0, err
	}
	totalOffset += n
	n, err = h.PrivateHeader.Encode(buf[n:])
	if err != nil {
		return 0, err
	}
	totalOffset += n
	return totalOffset, nil
}

func (h *Header) Decode(data []byte) (int, error) {
	if len(data) < HeaderSize {
		return 0, ErrSmallBufferSize
	}
	var totalOffset int

	n, err := h.PublicHeader.Decode(data)
	if err != nil {
		return 0, err
	}
	totalOffset += n
	n, err = h.PrivateHeader.Decode(data[n:])
	if err != nil {
		return 0, err
	}
	totalOffset += n
	if totalOffset != HeaderSize {
		return 0, ErrHeaderSizeMismatch
	}
	return totalOffset, nil
}
