package aes

import (
	"fmt"
	"testing"
)

func Test_Aes(t *testing.T) {

	aes := AesCbc{Key: "abcabc", Ken_len: 16}
	b, _ := aes.Encrypt([]byte("gonat_port:"))
	//b, _ = aes.Decrypt(b)
	fmt.Println(len(b))

}
