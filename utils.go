package rudp

import (
	"log"

	"github.com/pkg/errors"
)


func LogError(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func RaiseError(err error) (interface{},error) {
	if err != nil {
		log.Panic(err)
		return nil, errors.WithStack(err)
	}
}