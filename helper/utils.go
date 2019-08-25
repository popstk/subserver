package helper

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"log"
	"os"
	"strconv"
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

// FmtStringReplace replace {key} ,key in map
func FmtStringReplace(format string, m map[string]string) string {
	result := format
	for k, v := range m {
		result = strings.Replace(result, "{"+k+"}", v, -1)
	}

	return result
}

// Base64Decode try base64 decode
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

// Number support number 1 and "1"
type Number int

// MarshalJSON -
func (n Number) MarshalJSON() ([]byte, error) {
	return json.Marshal(int(n))
}

// UnmarshalJSON -
func (n *Number) UnmarshalJSON(b []byte) error {
	value := string(b)
	value = strings.Trim(value, "\"")
	val, err := strconv.Atoi(value)
	if err != nil {
		return err
	}

	*n = Number(val)
	return nil
}
