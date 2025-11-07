package courier

import "errors"

var (
	ErrCourierNotFound     = errors.New("courier wasn't found")
	ErrCourierExistPhone   = errors.New("courier with such a phone already exists")
	ErrCourierEmptyData    = errors.New("required courier fields aren't filled")
	ErrCourierInvalidPhone = errors.New("courier's phone number is incorrect")
	ErrCourierInvalidData  = errors.New("courier's information is incorrect")
	ErrCourierInvalidID    = errors.New("courier's ID is incorrect")
)
