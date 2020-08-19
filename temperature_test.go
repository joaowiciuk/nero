package main

import (
	"testing"
	"time"
)

func Test_watchTemperature(t *testing.T) {
	type args struct {
		command string
		pattern string
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "Test 1",
			args: args{
				command: "sensors",
				pattern: `Core 0:\ +(\+.*?)°C`,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempChan := watchTemperature(tt.args.command, tt.args.pattern)
			ticker := time.Tick(4 * time.Second)
			temperatures := make([]float64, 0)
			select {
			case _ = <-ticker:
				if len(temperatures) == 0 {
					t.Errorf("timeout: watchTemperature returned 0 values\n")
				}
			case temp := <-tempChan:
				temperatures = append(temperatures, temp)
			}
		})
	}
}
