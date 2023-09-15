package md5

import (
	"fmt"
	"testing"
)

func TestMd5(t *testing.T) {
	plaintext := "hello word"

	ciphertext := Encode(plaintext)
	fmt.Println("md5 Encode:", ciphertext)

	ciphertext16 := Encode16(plaintext)
	fmt.Println("md5 Encode 16:", ciphertext16)
}
