package cmpdb

import (
	"bytes"
	"fmt"
)

func SeparateDiff(data []byte) (before []byte, after []byte, err error) {
	beforeJSON, err := getJSON(data, []byte("- "), []byte("+ "))
	if err != nil {
		return nil, nil, err
	}

	afterJSON, err := getJSON(data, []byte("+ "), []byte("- "))
	if err != nil {
		return nil, nil, err
	}

	return beforeJSON, afterJSON, nil
}

func getJSON(data []byte, targetDiffType, noTargetDiffType []byte) ([]byte, error) {
	lineNumber := 1

	jsonData := make([]byte, 0, len(data))

	offset := 0

	for {
		switch {
		case
			bytesEqual(data, offset, targetDiffType),
			bytesEqual(data, offset, []byte("  ")):
			offset += 2

			index := bytes.IndexByte(data[offset:], '\n')
			if index == -1 {
				jsonData = append(jsonData, data[offset:]...)

				return jsonData, nil
			}

			jsonData = append(jsonData, data[offset:offset+index+1]...)
			offset = offset + index + 1
			lineNumber++
		case
			bytesEqual(data, offset, noTargetDiffType):
			offset += 2

			index := bytes.IndexByte(data[offset:], '\n')
			if index == -1 {
				return jsonData, nil
			}

			offset = offset + index + 1
			lineNumber++
		default:
			return nil, fmt.Errorf("not diff format, line: %d", lineNumber)
		}

		if offset >= len(data) {
			return jsonData, nil
		}

	}

}

func bytesEqual(data []byte, dataOffSet int, tar []byte) bool {
	src := data[dataOffSet : dataOffSet+len(tar)]

	if len(src) != len(tar) {
		return false
	}

	for i := 0; i < len(src); i++ {
		if src[i] != tar[i] {
			return false
		}
	}

	return true
}
