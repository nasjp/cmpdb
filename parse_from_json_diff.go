package cmpdb

import (
	"bytes"
	"encoding/json"
)

func ParseFromJSONDiff(data []byte) (*DBDiff, error) {
	beforeData, afterData, err := SeparateDiff(data)
	if err != nil {
		return nil, err
	}

	beforeDB, err := NewDBDecoder(json.NewDecoder(bytes.NewBuffer(beforeData))).Decode()
	if err != nil {
		return nil, err
	}

	afterDB, err := NewDBDecoder(json.NewDecoder(bytes.NewBuffer(afterData))).Decode()
	if err != nil {
		return nil, err
	}

	return &DBDiff{BeforeDB: beforeDB, AfterDB: afterDB}, nil
}
