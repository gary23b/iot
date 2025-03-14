package gocode

import (
	"fmt"
	"strings"
)

/*
https://learn.sparkfun.com/tutorials/avr-based-serial-enabled-lcds-hookup-guide/all
*/
type SparkfunSerialLcd struct {
	i2c     *I2C
	address uint16
}

func NewSparkfunSerialLcd(i2c *I2C, deviceAddress uint16) *SparkfunSerialLcd {
	ret := &SparkfunSerialLcd{
		i2c:     i2c,
		address: deviceAddress,
	}

	return ret
}

func (s *SparkfunSerialLcd) Write(in string) error {
	filteredIn := strings.Builder{}
	filteredIn.Grow(len(in) + 10)

	for _, r := range in {
		if r == '|' {
			filteredIn.WriteString("||")
			continue
		}
		filteredIn.WriteRune(r)
	}
	err := s.i2c.Write(s.address, []byte(filteredIn.String()))
	if err != nil {
		return fmt.Errorf("SparkfunSerialLcd.Write had error, in=%s: %w", in, err)
	}
	return nil
}

func (s *SparkfunSerialLcd) ClearDisplay() error {
	err := s.i2c.Write(s.address, []byte("|-"))
	if err != nil {
		return fmt.Errorf("SparkfunSerialLcd.ClearDisplay had error: %w", err)
	}
	return nil
}

func (s *SparkfunSerialLcd) MoveCursorTo(row, column int) error {
	if row < 0 || row >= 4 || column < 0 || column >= 20 {
		return fmt.Errorf("SparkfunSerialLcd.MoveCursorTo invalid input row=%d, column=%d", row, column)
	}

	rowsLookup := [4]int{0, 64, 20, 84}

	cmd := 128 + rowsLookup[row] + column

	err := s.i2c.Write(s.address, []byte{254, byte(cmd)})
	if err != nil {
		return fmt.Errorf("SparkfunSerialLcd.MoveCursorTo had error, row=%d, column=%d: %w", row, column, err)
	}
	return nil
}

func (s *SparkfunSerialLcd) SetBacklightPercent(r, g, b float64) error {
	if r < 0 || r > 100 || g < 0 || g > 100 || b < 0 || b > 100 {
		return fmt.Errorf("SparkfunSerialLcd.SetBacklightPercent invalid input rgb= %f %f %f", r, g, b)
	}

	steps := 29.0
	rInt := 128 + int(steps*r)
	gInt := 158 + int(steps*g)
	bInt := 188 + int(steps*b)

	cmd := []byte{'|', byte(rInt), '|', byte(gInt), '|', byte(bInt)}
	err := s.i2c.Write(s.address, cmd)
	if err != nil {
		return fmt.Errorf("SparkfunSerialLcd.SetBacklightPercent had error, rgb= %f %f %f: %w", r, g, b, err)
	}
	return nil
}
