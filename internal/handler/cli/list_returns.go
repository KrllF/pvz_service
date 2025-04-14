package cli

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/handler/cli/helper"
)

func (h *clihandler) ListReturns(ctx context.Context, args []string) {
	limit, page := int64(1), int64(1)

	if len(args) > 0 {
		var err error
		limit, err = helper.ParseOptionalInt(args[0], "limit")
		if err != nil {
			fmt.Printf("Ошибка при разборе параметра 'limit': %v\n", err)

			return
		}
	}

	if len(args) > 1 {
		var err error
		page, err = helper.ParseOptionalInt(args[1], "page")
		if err != nil {
			fmt.Printf("Ошибка при разборе параметра 'page': %v\n", err)

			return
		}
	}

	returns, count, err := h.showServ.ListReturns(ctx, limit, page)
	if err != nil {
		fmt.Printf("Ошибка при получении списка возвратов: %v\n", err)

		return
	}

	fmt.Println("Количество возвратов: ", count)
	fmt.Printf("Возвраты: \n\n")

	for _, val := range returns {
		fmt.Print(val, "\n\n")
	}
}
