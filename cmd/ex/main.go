package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type ver struct {
	FirstName string
	LastName  string
}

func encode(va interface{}) string {
	b, _ := json.Marshal(va)
	encoder := base64.StdEncoding.EncodeToString(b)

	return encoder
}

func decode(v string, user interface{}) interface{} {
	b, _ := base64.StdEncoding.DecodeString(v)
	c := json.Unmarshal(b, user)

	return c
}

func main() {
	user := ver{
		FirstName: "miloud",
		LastName:  "bmr",
	}

  b := encode(&user)
  fmt.Printf("encode: %s\n", b)

  var f ver
	s, _ := base64.StdEncoding.DecodeString(b)
	json.Unmarshal(s, &f)
  fmt.Printf("decode: %+v", f)
}
