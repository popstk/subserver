package helper

import (
	"bufio"
	"encoding/base64"
	"log"
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

func Base64Decode(data string) (string, error) {
	bytes, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		bytes, err = base64.RawStdEncoding.DecodeString(data)
	}

	if err != nil {
		log.Println("Base64Decode: ", data)
	}

	return string(bytes), err
}