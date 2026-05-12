package cache

import "errors"

var CannotSetKeyError = errors.New("Error on set redis key")
var CannotGetKeyError = errors.New("Error on get redis key")
