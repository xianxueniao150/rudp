package rudp

import (
	"log"

)


func LogError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func AppendHead(buf []byte, ele... byte) []byte{
	buf=append(buf, ele...)
	copy(buf[len(ele):],buf[:len(buf)-len(ele)])
	copy(buf[:len(ele)],ele)
	return buf
}


// func RaiseError(err error) (interface{},error) {
// 	if err != nil {
// 		log.Panic(err)
// 		return nil, errors.WithStack(err)
// 	}
// }