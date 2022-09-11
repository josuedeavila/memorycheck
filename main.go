package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"memorycheck/internal/memory"
)

func panicRecover() {
	if r := recover(); r != nil {
		fmt.Printf("panic!")
	}
	fmt.Print("recovered and stop")
}

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	defer panicRecover()

	sys := memory.Linux{}
	monitor := memory.NewMonitor(2, done, sys)
	go monitor.Memory(95)
	s := []string{}
	go func() {
		for {
			s = append(s, "Testing...")
		}
	}()
	<-done
}
