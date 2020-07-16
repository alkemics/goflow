package nodes

import (
	"errors"
)

func ErrorGenerator() (err error) {
	return errors.New("error")
}

func Failer(value string, raises bool) (transmitted string, err error) {
	if raises {
		return "", errors.New("error")
	}

	return value, nil
}

func Invalid2Errors() (err1, err2 error) {
	return nil, nil
}
