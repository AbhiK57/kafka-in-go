package main

import (
	"encoding/binary"
	"io"
)

// --- Binary Protocol Helpers ---

func readInt16(r io.Reader) (int16, error) {
	var v int16
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func readInt32(r io.Reader) (int32, error) {
	var v int32
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func readInt64(r io.Reader) (int64, error) {
	var v int64
	err := binary.Read(r, binary.BigEndian, &v)
	return v, err
}

func readString(r io.Reader) (string, error) {
	len, err := readInt16(r)
	if err != nil || len == -1 {
		return "", err
	}
	buf := make([]byte, len)
	_, err = io.ReadFull(r, buf)
	return string(buf), err
}

func readBytes(r io.Reader) ([]byte, error) {
	len, err := readInt32(r)
	if err != nil || len == -1 {
		return nil, err
	}
	buf := make([]byte, len)
	_, err = io.ReadFull(r, buf)
	return buf, err
}

func writeInt16(w io.Writer, v int16) {
	binary.Write(w, binary.BigEndian, v)
}

func writeInt32(w io.Writer, v int32) {
	binary.Write(w, binary.BigEndian, v)
}

func writeInt64(w io.Writer, v int64) {
	binary.Write(w, binary.BigEndian, v)
}

func writeString(w io.Writer, s string) {
	writeInt16(w, int16(len(s)))
	w.Write([]byte(s))
}

func writeBytes(w io.Writer, b []byte) {
	writeInt32(w, int32(len(b)))
	if len(b) > 0 {
		w.Write(b)
	}
}
