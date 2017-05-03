package utils

import (
	"errors"
	"net/url"
	"strconv"
)

// GetUintParamFromURL do as defined
func GetUintParamFromURL(u url.Values, name string, defalt uint64) (uint64, error) {
	// defalt is not a typo - default is reserved keyword
	var err error
	result := defalt
	paramStr := u.Get(name)
	if len(paramStr) != 0 {
		result, err = strconv.ParseUint(paramStr, 10, 0)
		if err != nil {
			return 0, errors.New(name + " not recognized")
		}
	}
	return result, nil
}
