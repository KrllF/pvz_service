package cli

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/handler/cli/helper"
)

func (h *clihandler) ReturnOrderToCourier(ctx context.Context, args []string) {
	if len(args) != 1 {
		fmt.Println("Недостаточно аргументов. Используйте: return_order_to_courier <order_id>")

		return
	}

	orderID, err := helper.ParseID(args[0])
	if err != nil {
		fmt.Printf("Ошибка при парсинге order_id: %v\n", err)

		return
	}

	err = h.returnsServ.ReturnOrder(ctx, orderID)
	if err != nil {
		fmt.Printf("Ошибка при возврате заказа: %v\n", err)
	}
	fmt.Println("Заказ успешно возвращён курьеру")
}
