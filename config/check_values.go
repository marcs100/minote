package config

import (
	"errors"
	"fmt"
	"strconv"
)

func CheckRGBString(rgb string) error {
	var err error = nil
	if len(rgb) != 7 {
		return errors.New(fmt.Sprint("length error in ", rgb))
	}

	if _, err := strconv.ParseUint(rgb[1:], 16, 64); err != nil {
		return errors.New(fmt.Sprint("format error in ", rgb))
	}

	return err
}
