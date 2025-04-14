package httph

import (
	"errors"
	"net/url"
	"strconv"
)

func parseFinderOrder(queryParams url.Values) (bool, string, error) {
	var orderFinder string

	finderStr := queryParams.Get("finder")
	finder := false
	if finderStr != "" {
		var err error
		finder, err = strconv.ParseBool(finderStr)
		if err != nil {
			return false, "", errors.New("некорректный параметр 'finder'")
		}
	}

	if finder {
		orderFinder = queryParams.Get("order")
		if orderFinder == "" {
			return false, "", errors.New("некорректный параметр 'order'")
		}
		_, err := strconv.ParseInt(orderFinder, 10, 64)
		if err != nil {
			return false, "", errors.New("некорректный параметр 'order'")
		}
	}

	return finder, orderFinder, nil
}
