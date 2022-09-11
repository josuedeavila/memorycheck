package memory

import (
	"fmt"
	"os"
	"syscall"
	"time"
)

// Monitor represents a mrmoty monitor
type Monitor struct {
	interval time.Duration
	done     chan os.Signal
	sys      OSMonitor
}

// NewMonitor creates a new monitor instance
func NewMonitor(t int, done chan os.Signal, sys OSMonitor) *Monitor {
	return &Monitor{
		interval: time.Duration(t) * time.Second,
		done:     done,
		sys:      sys,
	}
}


// Memory is responsable to monitor memory usage and send signal 
func (m *Monitor) Memory(threshold float64) {
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
		if *memory > threshold {
			m.done <- syscall.SIGTERM
		}
	}
}
