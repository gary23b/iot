package main

import (
	"fmt"
	"time"

	"github.com/gary23b/iot/gocode"
	"periph.io/x/conn/v3/gpio"
	"periph.io/x/host/v3"
	"periph.io/x/host/v3/rpi"
)

func NoError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Received unexpected error:\n%+v", err))
	}
}

func main() {
	fmt.Println("Hello World")
	host.Init()

	i2c, err := gocode.OpenI2c()
	NoError(err)

	lcd := gocode.NewSparkfunSerialLcd(i2c, 0x72)

	err = lcd.SetBacklightPercent(0, 0, 0)
	NoError(err)

	err = lcd.ClearDisplay()
	NoError(err)
	err = lcd.Write("Hello World!\nSecond Line")
	NoError(err)
	err = lcd.MoveCursorTo(3, 10)
	NoError(err)
	err = lcd.Write("| ! ||")
	NoError(err)

	t := time.NewTicker(500 * time.Millisecond)
	for l := gpio.Low; ; l = !l {
		rpi.P1_33.Out(l)
		<-t.C
	}

}
