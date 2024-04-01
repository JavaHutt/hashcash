package mocks

import (
	"bytes"
	"net"
	"time"
)

type MockConn struct {
	ReadBuffer  bytes.Buffer
	WriteBuffer bytes.Buffer
}

func (mc *MockConn) Read(b []byte) (n int, err error) {
	return mc.ReadBuffer.Read(b)
}

func (mc *MockConn) Write(b []byte) (n int, err error) {
	return mc.WriteBuffer.Write(b)
}

func (mc *MockConn) Close() error {
	return nil
}

func (mc *MockConn) LocalAddr() net.Addr {
	return nil
}

func (mc *MockConn) RemoteAddr() net.Addr {
	return nil
}

func (mc *MockConn) SetDeadline(_ time.Time) error {
	return nil
}

func (mc *MockConn) SetReadDeadline(_ time.Time) error {
	return nil
}

func (mc *MockConn) SetWriteDeadline(_ time.Time) error {
	return nil
}
