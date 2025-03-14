package gocode

import (
	"fmt"
	"sync"
	"time"

	"periph.io/x/conn/v3"
	"periph.io/x/conn/v3/driver/driverreg"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

func pinControl() {
	fmt.Println("Hello World")
	host.Init()

	sendMessageToScreen()

	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		rpi.P1_33.Out(l)
		<-t.C
	}
}

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
	if err != nil {
		return fmt.Errorf("error while WriteThenRead I2C, Device=0x%X, WriteData=%v, readLength=%d: %w", address, writeData, readLength, err)
	}
	return nil
}

func sendMessageToScreen() {
	// Make sure periph is initialized.
	// TODO: Use host.Init(). It is not used in this example to prevent circular
	// go package import.
	if _, err := driverreg.Init(); err != nil {
		panic(err)
	}

	// Use i2creg I2C bus registry to find the first available I2C bus.
	b, err := i2creg.Open("")
	if err != nil {
		panic(err)
	}
	defer b.Close()

	// Dev is a valid conn.Conn.
	d := &i2c.Dev{Addr: 0x72, Bus: b}

	// Send a command 0x10 and expect a 5 bytes reply.
	write := []byte{'h', 'e', 'l', 'l', 'o'}
	read := make([]byte, len(write))
	if err := d.Tx(write, read); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", read)
}
