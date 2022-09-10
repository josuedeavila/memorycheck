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

type Linux struct {
	stats Stats
}

func (l Linux) GetUsedPercentage() (*float64, error) {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	stats, err := getMemoryStats(file)
	if err != nil {
		return nil, err
	}

	c := float64(stats.Used) / float64(stats.Total) * 100
	return &c, nil
}


func getMemoryStats(out io.Reader) (*Stats, error) {
	var memory Stats
	memStats := map[string]*uint64{
		"MemTotal":     &memory.Total,
		"MemFree":      &memory.Free,
		"MemAvailable": &memory.Available,
		"Buffers":      &memory.Buffers,
		"Cached":       &memory.Cached,
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
			memory.MemAvailableEnabled = true
		}
	}

	err := scanner.Err()
	if err != nil {
		return nil, fmt.Errorf("scan error for /proc/meminfo: %s", err)
	}

	if memory.MemAvailableEnabled {
		memory.Used = memory.Total - memory.Available
	} else {
		memory.Used = memory.Total - memory.Free - memory.Buffers - memory.Cached
	}

	return &memory, nil
}
