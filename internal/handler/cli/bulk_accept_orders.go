package cli

import (
	"context"
	"fmt"
)

func (h *clihandler) BulkAcceptOrders(ctx context.Context, args []string) {
	if len(args) != 1 {
		fmt.Println("Недостаточно аргументов. Используйте: bulk_accept_orders <path>.")

		return
	}
	n, err := h.acceptServ.BulkAcceptOrders(ctx, args[0])
	if err != nil {
		fmt.Printf("Ошибка при добавлении заказов: %v\n", err)
	}
	fmt.Printf("Добавлено %d новых заказов\n", n)
}
