package helper

import (
	"bufio"
	"os"
	"strings"
)

// ReadLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}


func FmtStringReplace(format string, m map[string]string) string {
	result := format
	for k,v := range m {
		result = strings.Replace(result, "{"+k+"}", v, -1)
	}

	return result
}