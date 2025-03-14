package gocode

import (
	"fmt"
	"sync"
	"time"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
)

type I2C struct {
	bus   i2c.BusCloser
	mutex sync.Mutex
}

func OpenI2c() (*I2C, error) {
	x, err := i2creg.Open("")
	if err != nil {
		return nil, fmt.Errorf("error while opening I2C: %w", err)
	}

	ret := &I2C{
		bus: x,
	}
	return ret, nil
}

func (s *I2C) Close() {
	s.bus.Close()
}

func (s *I2C) GetConnection(address uint16) conn.Conn {
	device := &i2c.Dev{Addr: address, Bus: s.bus}
	return device
}

func (s *I2C) Write(address uint16, writeData []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	device := &i2c.Dev{Addr: address, Bus: s.bus}
	err := device.Tx(writeData, nil)
	time.Sleep(time.Millisecond)
	if err != nil {
		return fmt.Errorf("error while writing I2C, Device=0x%X, WriteData=%v: %w", address, writeData, err)
	}
	return nil
}

func (s *I2C) Read(address uint16, readLength int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	device := &i2c.Dev{Addr: address, Bus: s.bus}
	readData := make([]byte, readLength)
	err := device.Tx(nil, readData)
	time.Sleep(time.Millisecond)
	if err != nil {
		return fmt.Errorf("error while reading I2C, Device=0x%X, readLength=%d: %w", address, readLength, err)
	}
	return nil
}

func (s *I2C) WriteThenRead(address uint16, writeData []byte, readLength int) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	device := &i2c.Dev{Addr: address, Bus: s.bus}
	readData := make([]byte, readLength)
	err := device.Tx(nil, readData)
	time.Sleep(time.Millisecond)
	if err != nil {
		return fmt.Errorf("error while WriteThenRead I2C, Device=0x%X, WriteData=%v, readLength=%d: %w", address, writeData, readLength, err)
	}
	return nil
}
