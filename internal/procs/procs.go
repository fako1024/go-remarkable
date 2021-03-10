package procs

import (
	"bufio"
	"fmt"
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

// PIDOf return the PID of a process based on its name (first match is returned)
func PIDOf(processName string) (string, error) {

	// Extract list of running processes
	processes, err := ioutil.ReadDir(procPath)
	if err != nil {
		return "", fmt.Errorf("failed to parse system processes: %s", err)
	}

	// Loop over all processes and look for a match
	for _, process := range processes {
		if process.IsDir() {

			pidStr := process.Name()

			if _, err := strconv.Atoi(pidStr); err == nil {

				path := filepath.Join(procPath, pidStr, statFileName)
				if _, err := os.Stat(path); err == nil {

					statBytes, err := ioutil.ReadFile(path)
					if err != nil {
						return "", fmt.Errorf("failed to read stat path %s: %s", path, err)
					}

					// Parse binary name from stat file
					statData := string(statBytes)
					binStart := strings.IndexRune(statData, '(') + 1
					binEnd := strings.IndexRune(statData[binStart:], ')')
					binary := statData[binStart : binStart+binEnd]

					if strings.Contains(binary, processName) {
						return pidStr, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("process `%s` not found", processName)
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
