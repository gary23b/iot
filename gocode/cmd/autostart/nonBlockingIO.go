package main

import (
	"bufio"
	"errors"
	"io"
	"io/fs"
	"log"
)

type NonBlockingReader struct {
	stringChan <-chan string
}

func NewNonBlockingReader(r io.Reader) *NonBlockingReader {
	ch := make(chan string, 1000)

	go func() {
		reader := bufio.NewReader(r)
		for {
			line, err := reader.ReadString('\n')
			if line != "" {
				ch <- line
			}
			if err != nil { // log non io.EOF errors
				if errors.Is(err, io.EOF) {
					break
				}
				if e := (&fs.PathError{}); errors.As(err, &e) {
					if e.Err.Error() == "file already closed" {
						break
					}
				}
				log.Println(err)
				break
			}
		}
		close(ch)
	}()

	ret := &NonBlockingReader{
		stringChan: ch,
	}

	return ret
}

func (s *NonBlockingReader) GetLine() string {
	select {
	case line := <-s.stringChan:
		return line
	default:
		// do nothing
	}
	return ""
}

type NonBlockingWriteCloser struct {
	writeChan chan<- []byte
}

var _ io.WriteCloser = &NonBlockingWriteCloser{}

func NewNonBlockingWriter(w io.WriteCloser) *NonBlockingWriteCloser {
	ch := make(chan []byte, 1000)

	go func() {
		for data := range ch {
			_, err := w.Write(data)
			if err != nil {
				log.Println(err)
			}
		}
		// when the channel closes, close the io.WriteCloser
		w.Close()
	}()

	ret := &NonBlockingWriteCloser{
		writeChan: ch,
	}

	return ret
}

func (s *NonBlockingWriteCloser) WriteString(in string) {
	s.writeChan <- []byte(in)
}

func (s *NonBlockingWriteCloser) Write(in []byte) (int, error) {
	inCopy := append([]byte{}, in...) // copy the byte slice

	s.writeChan <- inCopy

	return len(inCopy), nil
}

func (s *NonBlockingWriteCloser) Close() error {
	close(s.writeChan)
	return nil
}
