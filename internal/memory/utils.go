package memory

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

func HostProc(path string) string {
	value := os.Getenv("HOST_PROC")
	if value == "" {
		value = fmt.Sprintf("/proc/%s", path)
	}
	return value
}

// ReadLines reads contents from a file and splits them by new lines.
func ReadLines(filename string) (content []string, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return content, err
	}
	defer f.Close()

	r := bufio.NewReader(f)
	for {
		line, err := r.ReadString('\n')
		if err != nil {
			if err == io.EOF && len(line) > 0 {
				content = append(content, strings.Trim(line, "\n"))
			}
			break
		}
		content = append(content, strings.Trim(line, "\n"))
	}

	return content, nil
}

// ParseMemStats parse mem stats to bytes
func ParseMemStats(value string) (*uint64, error) {
	mem, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return &mem, err
	}
	memKb := mem * 1024
	return &memKb, nil
}
