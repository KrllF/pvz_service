package grpc

import (
	"errors"
)

type reqInter interface {
	GetFinder() bool
	GetOrderFinder() string
}

func parseFinderOrder(req reqInter) (bool, string, error) {
	var orderFinder string

	finderStr := req.GetFinder()

	if finderStr {
		orderFinder = req.GetOrderFinder()
		if orderFinder == "" {
			return false, "", errors.New("некорректный параметр 'order'")
		}
	}

	return finderStr, orderFinder, nil
}
