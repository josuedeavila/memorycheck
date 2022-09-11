package memory

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// Stats represents memory statistics for linux
type Stats struct {
	Total, Used, Buffers, Cached, Free, Available uint64
	MemAvailableEnabled                           bool
}


// Linux represents linux as a entity on Monitor context
type Linux struct {
	stats Stats
}


// GetUsedPercentage return the percentage of memory used on OS
func (l Linux) GetUsedPercentage() (*float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, err := l.getMemoryStats(file)
	if err != nil {
		return nil, err
	}

	c := float64(stats.Used) / float64(stats.Total) * 100
	return &c, nil
}


func (l Linux)getMemoryStats(out io.Reader) (*Stats, error) {
	memStats := map[string]*uint64{
		"MemTotal":     &l.stats.Total,
		"MemFree":      &l.stats.Free,
		"MemAvailable": &l.stats.Available,
		"Buffers":      &l.stats.Buffers,
		"Cached":       &l.stats.Cached,
	}

	scanner := bufio.NewScanner(out)
	for scanner.Scan() {
		line := scanner.Text()
		i := strings.IndexRune(line, ':')
		if i < 0 {
			continue
		}
		statKey := line[:i]
		if statValue := memStats[statKey]; statValue == nil {
			continue
		}
		v := strings.TrimSpace(strings.TrimRight(line[i+1:], "kB"))
		if v, err := strconv.ParseUint(v, 10, 64); err == nil {
			*memStats[statKey] = v * 1024
		}
		if statKey == "MemAvailable" {
			l.stats.MemAvailableEnabled = true
		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("scan error for /proc/meminfo: %s", err)
	}

	if l.stats.MemAvailableEnabled {
		l.stats.Used = l.stats.Total - l.stats.Available
		return &l.stats, nil
	}
	
	l.stats.Used = l.stats.Total - l.stats.Free - l.stats.Buffers - l.stats.Cached
	return &l.stats, nil
}
