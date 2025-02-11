package codekit

import (
	"bytes"
	"encoding/json"
)

func Jsonify(o interface{}) []byte {
	bs, e := json.Marshal(o)
	if e != nil {
		bs, _ = json.Marshal(struct{}{})
	}
	return bs
}

func JsonString(o interface{}) string {
	bs := Jsonify(o)
	return string(bs)
}

var _indentChar string

func SetIndentChar(c string) {
	_indentChar = c
}

func IndentChar() string {
	if _indentChar == "" {
		_indentChar = " "
	}
	return _indentChar
}

func JsonStringIndent(o interface{}, indentChar string) string {
	bs, e := json.MarshalIndent(o, "", IndentChar())
	if e != nil {
		return "{}"
	}
	return string(bs)
}

func Unjson(b []byte, result interface{}) error {
	buff := bytes.NewBuffer(b)
	decoder := json.NewDecoder(buff)
	decoder.UseNumber()
	e := decoder.Decode(&result)
	return e
}

func UnjsonFromString(s string, result interface{}) error {
	b := bytes.NewBufferString(s)
	decoder := json.NewDecoder(b)
	decoder.UseNumber()
	e := decoder.Decode(&result)
	return e
}
