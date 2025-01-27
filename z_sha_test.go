package codekit_test

import (
	"bytes"
	"encoding/base64"
	"testing"

	"github.com/sebarcode/codekit"
)

func TestSha(t *testing.T) {
	data := "A very random string"
	salt := "A very difficult secret"

	shaBytes := codekit.Sha(data, salt)
	shaString := codekit.ShaString(data, salt)

	decodedBytes, _ := base64.StdEncoding.DecodeString(shaString)
	if !bytes.Equal(decodedBytes, shaBytes) {
		t.Errorf("fail compare hash")
	}
}
