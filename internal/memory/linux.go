package memory

import (
	"encoding/json"
	"strings"
)

type MemoryStat struct {
	Total        uint64  `json:"MemTotal"`
	Available    uint64  `json:"MemAvailable"`
	Used         uint64  `json:"Used"`
	UsedPercent  float64 `json:"UsedPercent"`
	Free         uint64  `json:"MemFree"`
	Buffers      uint64  `json:"Buffers"`
	Cached       uint64  `json:"Cached"`
	Sreclaimable uint64  `json:"Sreclaimable"`
}

// Linux represents linux as a entity on Monitor context
type Linux struct {
	memoryStats MemoryStat
}

// GetUsedPercentage return the percentage of memory used on OS
func (l Linux) GetUsedPercentage() (*float64, error) {
	stats, err := l.getMemoryStats()
	if err != nil {
		return nil, err
	}

	return &stats.UsedPercent, nil
}

func (l Linux) getMemoryStats() (*MemoryStat, error) {
	filename := HostProc("meminfo")
	lines, _ := ReadLines(filename)

	memavailable := false
	m := map[string]*uint64{
		"MemTotal":     nil,
		"MemAvailable": nil,
		"Used":         nil,
		"UsedPercent":  nil,
		"MemFree":      nil,
		"Buffers":      nil,
		"Cached":       nil,
		"SReclaimable": nil,
	}
	for _, line := range lines {
		fields := strings.Split(line, ":")
		if len(fields) != 2 {
			continue
		}
		key := strings.TrimSpace(fields[0])
		value := strings.TrimSpace(fields[1])
		value = strings.Replace(value, " kB", "", -1)

		if _, ok := m[key]; !ok {
			continue
		}

		if key == "MemAvailable" {
			memavailable = true
		}

		stat, err := ParseMemStats(value)
		if err != nil {
			return nil, err
		}
		m[key] = stat
	}

	b, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	if err = json.Unmarshal(b, &l.memoryStats); err != nil {
		return nil, err
	}

	l.memoryStats.Cached += l.memoryStats.Sreclaimable
	l.memoryStats.Used = l.memoryStats.Total - l.memoryStats.Free - l.memoryStats.Buffers - l.memoryStats.Cached
	l.memoryStats.UsedPercent = float64(l.memoryStats.Used) / float64(l.memoryStats.Total) * 100.0

	if !memavailable {
		l.memoryStats.Available = l.memoryStats.Cached + l.memoryStats.Free
		return &l.memoryStats, nil
	}

	return &l.memoryStats, nil
}
