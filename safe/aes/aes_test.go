package aes

import (
	"fmt"
	"testing"
)

func Test_Aes(t *testing.T) {

	aes := AesCbc{Key: "abcabc", KenLen: 16}
	b, _ := aes.Encrypt([]byte("gonat_port:"))
	//b, _ = aes.Decrypt(b)
	fmt.Println(len(b))

}
