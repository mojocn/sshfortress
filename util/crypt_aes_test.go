package util

import (
	"encoding/base64"
	"testing"
)

func TestAesDecryptEn(t *testing.T) {
	key := "NlVaUmFnTFkxMjM0NTZrdWQ3R3o"

	//msg := "123456"

	crytedByte, err := base64.RawStdEncoding.DecodeString(key)
	if err != nil {
		t.Log(err)
	}

	t.Log(string(crytedByte[8 : len(crytedByte)-6]))
}
