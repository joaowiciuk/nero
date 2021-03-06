package main

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"time"
)

type reading struct {
	when  time.Time
	value float64
}

func watchTemperature(command, pattern string) chan reading {
	output := make(chan reading)
	regex := regexp.MustCompile(pattern)
	buffer := make([]byte, 400)

	ticker := time.Tick(2 * time.Second)
	go func() {
		for {
			select {
			case when := <-ticker:
				cmd := exec.Command(command)
				stdout, err := cmd.StdoutPipe()
				if err != nil {
					log.Printf("error getting stdout: %v", err)
					continue
				}

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
				output <- reading{
					when:  when,
					value: temp,
				}
			}
		}
	}()
	return output
}
