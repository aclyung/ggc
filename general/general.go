package general

import "log"

func ErrCheck(err error) {
	if err != nil {
		log.Fatal(err.Error())
	}
}
