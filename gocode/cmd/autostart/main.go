package main

import (
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"time"
)

// https://raspberrypi-guide.github.io/programming/run-script-on-boot
// https://github.com/thagrol/Guides/blob/main/boot.pdf
func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	// workingDir := os.Args[0]
	// gitRepo := os.Args[1]
	// goProgramDir := os.Args[2]
	workingDir := "/home/vm/iot_program"
	gitRepo := "https://github.com/gary23b/iot.git"
	goProgramDir := "/home/vm/iot_program/iot/gocode/cmd/sprinklers"

	err := os.MkdirAll(workingDir, os.ModePerm)
	if err != nil {
		panic(err)
	}

	os.Chdir(workingDir)
	if !DoesFolderExist(goProgramDir) {
		c, err := NewChildProgram("git", "clone", gitRepo)
		if err != nil {
			panic(err)
		}
		err = c.cmd.Wait()
		if err != nil {
			fmt.Print(c.ReadAllStdOut())
			log.Print(c.ReadAllStdErr())
			panic(err)
		}
	}

	os.Chdir(goProgramDir)

	restartDelayTime := 5 * time.Second

	for {
		c, err := NewChildProgram("git", "pull")
		if err != nil {
			panic(err)
		}
		err = c.Wait()
		log.Println(c.cmd.Path, c.cmd.Args)
		fmt.Print(c.ReadAllStdOut())
		log.Print(c.ReadAllStdErr())
		if err != nil {
			panic(err)
		}

		c, err = NewChildProgram("go", "run", goProgramDir)
		if err != nil {
			panic(err)
		}

		err = c.cmd.Wait()
		if err != nil {
			log.Println(err)
			log.Println(c.cmd.Path, c.cmd.Args)
			fmt.Print(c.ReadAllStdOut())
			log.Print(c.ReadAllStdErr())
		}

		fmt.Println("Waiting ", restartDelayTime, " before trying again.")
		time.Sleep(restartDelayTime)
		restartDelayTime = time.Millisecond * time.Duration(math.Round(restartDelayTime.Seconds()*1000*1.05))
	}
}

func DoesFolderExist(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}
