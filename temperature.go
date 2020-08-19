package main

import (
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

func watchTemperature(command, pattern string) (chan float64, error) {
	output := make(chan float64)
	cmd := exec.Command(command)
	regex := regexp.MustCompile(pattern)
	buffer := make([]byte, 400)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("error: couldn't get command output pipe: %v", err)
	}

	ticker := time.Tick(2 * time.Second)
	go func() {
		for {
			select {
			case _ = <-ticker:
				if err := cmd.Start(); err != nil {
					log.Printf("error starting command: %v\n", err)
					continue
				}

				n, _ := stdout.Read(buffer)

				if err := cmd.Wait(); err != nil {
					log.Printf("error running command: %v\n", err)
					continue
				}

				bufferStr := string(buffer[:n])
				submatch := regex.FindStringSubmatch(bufferStr)
				if len(submatch) != 2 {
					log.Printf("error reading command output: data or pattern is invalid\n")
					continue
				}
				str := submatch[1]
				temp, err := strconv.ParseFloat(str, 64)
				if err != nil {
					log.Printf("error parsing output: %v\n", err)
					continue
				}
				output <- temp
			}
		}
	}()
	return output, nil
}
