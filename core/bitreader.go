package core

import "io"

type BitReader struct {
	reader io.Reader
	buffer byte
	count  uint8
}

func NewBitReader(r io.Reader) *BitReader {
	return &BitReader{reader: r}
}

func (br *BitReader) ReadBit() (rune, error) {
	if br.count == 0 {
		var buf [1]byte
		_, err := br.reader.Read(buf[:])
		if err != nil {
			return 0, err
		}
		br.buffer = buf[0]
		br.count = 8
	}

	bit := (br.buffer >> 7) & 1
	br.buffer <<= 1
	br.count--

	if bit == 1 {
		return '1', nil
	}
	return '0', nil
}