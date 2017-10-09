package utils

import (
	"net/url"
	"strconv"
)

// GetUintQueryParam do as defined
func GetUintQueryParam(u url.Values, name string, defalt uint64) (uint64, error) {
	// defalt is not a typo - default is reserved keyword
	var err error
	result := defalt
	paramStr := u.Get(name)
	if len(paramStr) != 0 {
		result, err = strconv.ParseUint(paramStr, 10, 0)
		if err != nil {
			return 0, err
		}
	}
	return result, nil
}
