package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

func encodeToBase64(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	return encoded
}

func decodeFromString(data string, user interface{}) error {
	decode, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return fmt.Errorf("error decoding string %s", err.Error())
	}
	err = json.Unmarshal(decode, user)
	if err != nil {
		return fmt.Errorf("error unmarshalig string %s", err.Error())
	}

	return nil
}
