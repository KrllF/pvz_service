package cli

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/handler/cli/helper"
	"github.com/KrllF/pvz_service/internal/models"
)

const (
	all = 1
	seg = 2
)

func (h *clihandler) ListOrdersUser(ctx context.Context, args []string) {
	if len(args) < 1 {
		fmt.Println("Недостаточно аргументов.")
		fmt.Println("Используйте: list_orders_user <user_id> [опционально введите число n, чтобы получить n заказов,")
		fmt.Println("или любое слово, чтобы получить заказы в ПВЗ")

		return
	}

	userID, err := helper.ParseID(args[0])
	if err != nil {
		fmt.Printf("Ошибка при парсинге user_id: %v\n", err)

		return
	}

	var ret []models.Order

	switch len(args) {
	case all:
		ret, err = h.showServ.ListOrdersUserAllN(ctx, userID, -1)
	case seg:
		num, errN := helper.ParseID(args[1])
		if errN != nil {
			ret, err = h.showServ.ListOrdersUserPVZ(ctx, userID)
		} else {
			ret, err = h.showServ.ListOrdersUserAllN(ctx, userID, num)
		}

	default:
		fmt.Println("Много аргументов.")
		fmt.Println("Используйте: list_orders_user <user_id> [опционально введите число n, чтобы получить n заказов,")
		fmt.Println("или любое слово, чтобы получить заказы в ПВЗ")

		return
	}

	if err != nil {
		fmt.Println("Ошибка при выполнении:", err)

		return
	}

	fmt.Println("Количество заказов по запросу:", len(ret))
	displayOrders(ret)
}

func displayOrders(orders []models.Order) {
	cursor := 0
	limit := 2

	for {
		tmp := 0
		for cursor < len(orders) && tmp < limit {
			fmt.Println(orders[cursor])
			cursor++
			tmp++
		}

		if cursor >= len(orders) {
			break
		}

		var str string
		fmt.Print("\nПродолжить просмотр? (next/stop(любая строка)): ")
		_, err := fmt.Scan(&str)
		if err != nil {
			fmt.Println("Ошибка при вводе")

			return
		}
		if str != "next" {
			break
		}
	}
}
