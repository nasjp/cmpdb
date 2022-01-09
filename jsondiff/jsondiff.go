package jsondiff

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func Unmarshal(data []byte, beforeDist interface{}, afterDist interface{}) error {
	beforeJSON, err := getJSON(data, []byte("- "), []byte("+ "))
	if err != nil {
		return err
	}

	afterJSON, err := getJSON(data, []byte("+ "), []byte("- "))
	if err != nil {
		return err
	}

	if err := json.Unmarshal(beforeJSON, beforeDist); err != nil {
		return err
	}

	if err := json.Unmarshal(afterJSON, afterDist); err != nil {
		return err
	}

	return nil
}

func NewDecoder(data []byte) (*json.Decoder, *json.Decoder, error) {
	beforeJSON, err := getJSON(data, []byte("- "), []byte("+ "))
	if err != nil {
		return nil, nil, err
	}

	afterJSON, err := getJSON(data, []byte("+ "), []byte("- "))
	if err != nil {
		return nil, nil, err
	}

	return json.NewDecoder(bytes.NewReader(beforeJSON)), json.NewDecoder(bytes.NewReader(afterJSON)), nil
}

func getJSON(data []byte, targetDiffType, noTargetDiffType []byte) ([]byte, error) {
	lineNumber := 1

	beforeData := make([]byte, 0, len(data))

	offset := 0

	for {
		switch {
		case
			bytesEqual(data, offset, targetDiffType),
			bytesEqual(data, offset, []byte("  ")):
			offset += 2

			index := bytes.IndexByte(data[offset:], '\n')
			if index == -1 {
				beforeData = append(beforeData, data[offset:]...)

				return beforeData, nil
			}

			beforeData = append(beforeData, data[offset:offset+index]...)
			offset = offset + index + 1
			lineNumber++
		case
			bytesEqual(data, offset, noTargetDiffType):
			offset += 2

			index := bytes.IndexByte(data[offset:], '\n')
			if index == -1 {
				return beforeData, nil
			}

			offset = offset + index + 1
			lineNumber++
		default:
			return nil, fmt.Errorf("not diff format, line: %d", lineNumber)
		}

		if offset >= len(data) {
			return beforeData, nil
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
