package aes

import (
	"fmt"
	"testing"
)

func Test_Aes(t *testing.T) {

	aes := AesCbc{Key: "gonat", Ken_len: 16}
	b, _ := aes.Encrypt([]byte("abc"))
	b, _ = aes.Decrypt(b)
	fmt.Println(string(b))

}
