package safe

import (
	"gonat/interface"
	"gonat/safe/aes"
)

//var handler_map map[string]_interface.Safe

//func init() {
//	register("aes-128-cbc", &aes.AesCbc{Ken_len: 16})
//}

//func register(name string, _type _interface.Safe) {
//	if handler_map == nil {
//		handler_map = map[string]_interface.Safe{}
//	}
//	handler_map[name] = _type
//}

func GetSafe(crypt string, key string) _interface.Safe {
	switch crypt {
	case "aes-128-cbc":
		return &aes.AesCbc{Ken_len: 16, Key: key}
	}
	return nil
}
