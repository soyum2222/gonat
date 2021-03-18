package proto

import (
	"encoding/binary"
	"fmt"
	"gonat/safe"
	"testing"
)

func Test01(t *testing.T) {

	p1 := Proto{}

	p1.Body = []byte("proto test")
	p1.Kind = TCP_COMM
	p1.ConversationID = 1
	h := safe.GetSafe("aes-128-cbc", "gonat")
	b := p1.Marshal(h)

	p2 := Proto{}

	p2.Unmarshal(b[4:], h)

	if p1.ConversationID == p2.ConversationID && p1.Kind == p2.Kind {
		return
	} else {
		t.Fail()
	}

}

func Test02(t *testing.T) {

	fmt.Println(binary.BigEndian.Uint32([]byte{0, 0, 1}))

}
