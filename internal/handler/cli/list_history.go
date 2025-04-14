package cli

import (
	"context"
	"fmt"
)

func (h *clihandler) ListHistory(ctx context.Context) {
	history, err := h.showServ.ListHistory(ctx)
	if err != nil {
		fmt.Printf("Ошибка при получении истории заказов: %v\n", err)

		return
	}
	for _, order := range history {
		fmt.Println(order)
	}
}
