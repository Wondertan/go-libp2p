package libp2pquic

import (
	"errors"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/quic-go/quic-go"
	"sync"
)

const (
	reset quic.StreamErrorCode = 0
)

type stream struct {
	quic.Stream
	sync.Once
}

var _ network.MuxedStream = &stream{}

func (s *stream) Read(b []byte) (n int, err error) {
	n, err = s.Stream.Read(b)
	if err != nil && errors.Is(err, &quic.StreamError{}) {
		err = network.ErrReset
	}
	return n, err
}

func (s *stream) Write(b []byte) (n int, err error) {
	n, err = s.Stream.Write(b)
	if err != nil && errors.Is(err, &quic.StreamError{}) {
		err = network.ErrReset
	}
	return n, err
}

func (s *stream) Reset() error {
	s.Stream.CancelRead(reset)
	s.Stream.CancelWrite(reset)
	return nil
}

func (s *stream) Close() (err error) {
	s.Do(func() {
		err = s.Stream.Close()
	})
	return err
}

func (s *stream) CloseRead() error {
	return nil
}

func (s *stream) CloseWrite() error {
	return s.Close()
}
