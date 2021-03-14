package procs

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

const (
	procPath     = "/proc"
	statFileName = "stat"
	mapsFileName = "maps"
	memFileName  = "mem"
	fb0DevName   = "/dev/fb0"
)

// PIDOf return the PID of a process based on its PID file
func PIDOf(pidFile string) (string, error) {

	file, err := ioutil.ReadFile(pidFile)
	if err != nil {
		return "", err
	}

	return string(bytes.TrimSpace(file)), nil
}

func MemoryOffset(pid string) (int64, error) {

	file, err := os.Open(filepath.Join(procPath, pid, mapsFileName))
	if err != nil {
		return 0, err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		line := scanner.Text()
		if strings.HasSuffix(line, fb0DevName) {

			scanner.Scan()
			hexAddr := strings.Split(strings.Fields(scanner.Text())[0], "-")[0]

			addr, err := strconv.ParseInt(hexAddr, 16, 64)
			if err != nil {
				return 0, err
			}

			return addr + 8, nil
		}
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	return 0, file.Close()
}
