package cli

import (
	"context"
	"fmt"

	"github.com/KrllF/pvz_service/internal/consts"
	"github.com/KrllF/pvz_service/internal/handler/cli/helper"
)

const (
	lessThenError = 3
)

func (h *clihandler) ProcessClient(ctx context.Context, args []string) {
	if len(args) < lessThenError {
		fmt.Println("Недостаточно аргументов. Используйте: process_client <user_id> <action> <order_ids...>")

		return
	}

	userID, err := helper.ParseID(args[0])
	if err != nil {
		fmt.Printf("Ошибка при парсинге user_id: %v\n", err)

		return
	}

	action := args[1]
	if action != consts.IssueAction && action != consts.ReturnAction {
		fmt.Printf("Некорректное действие. Используйте %s или %s", consts.IssueAction, consts.ReturnAction)

		return
	}

	orderIDs, err := helper.ParseIDs(args[2:])
	if err != nil {
		fmt.Printf("Ошибка при парсинге order_ids: %v\n", err)

		return
	}

	var n int64

	switch action {
	case "issue":
		n, err = h.acceptServ.ProcessClientIssue(ctx, userID, orderIDs)
	case "return":
		n, err = h.returnsServ.ProcessClientReturn(ctx, userID, orderIDs)
	}
	if err != nil {
		fmt.Printf("Ошибка при обработке некоторых клиентов: %v\n", err)
	}
	fmt.Printf("Обработано заказов %d\n", n)
	fmt.Println("Операция успешно выполнена")
}
