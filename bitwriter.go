package main

import "io"

type BitWriter struct {
	writer io.Writer
	buffer byte
	count int
}

func NewBitWriter(w io.Writer) *BitWriter {
	return &BitWriter{writer: w}
}

func (bw *BitWriter) Writebits(bitString string) error {
	for _, char := range bitString {
		bw.buffer <<= 1
		if char == '1' {
			bw.buffer |= 1
		}
		bw.count += 1

		if bw.count == 8 {
			_, err := bw.writer.Write([]byte{bw.buffer})
			if err != nil {
				return err
			}
			bw.buffer = 0
			bw.count = 0
		}
	}
	return nil
}

func (bw *BitWriter) Flush() error {
	if bw.count > 0 {
		bw.buffer <<= (8 - bw.count)
		_, err := bw.writer.Write([]byte{bw.buffer})
		if err != nil {
			return err
		}
		bw.buffer = 0
		bw.count = 0
	}
	return nil
}