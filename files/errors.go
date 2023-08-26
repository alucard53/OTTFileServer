package files

import "errors"

var FileNotFound = errors.New("File not Found")
var RangeError = errors.New("Invalid Range Header")
