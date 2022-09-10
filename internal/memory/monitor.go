package memory

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

type Monitor struct {
	interval  time.Duration
	done      chan os.Signal
	sys OSMonitor
}

func NewMonitor(t int, done chan os.Signal, sys OSMonitor) *Monitor {
	return &Monitor{
		interval:  time.Duration(t) * time.Second,
		done:      done,
		sys: sys,
	}
}

func (m *Monitor) Memory() {
	for {
		<-time.After(m.interval)
		memory, err := m.sys.GetUsedPercentage()
		if err != nil {
			return
		}

		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			return
		}

		fmt.Println(*memory)
		if *memory > 90 {
			m.done <- syscall.SIGTERM
		}
	}
}
