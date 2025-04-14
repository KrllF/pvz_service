package cli

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

type commandHandler func(context.Context, []string)

func (h *clihandler) Run(ctx context.Context) error {
	fmt.Println("Добро пожаловать в систему управления ПВЗ!")
	fmt.Println("Введите 'help' для просмотра доступных команд.")
	reader := bufio.NewReader(os.Stdin)
	commandHandlers := map[string]commandHandler{
		"help":                    func(context.Context, []string) { h.Help() },
		"accept_order":            h.AcceptOrder,
		"return_order_to_courier": h.ReturnOrderToCourier,
		"process_client":          h.ProcessClient,
		"list_orders_user":        h.ListOrdersUser,
		"list_returns":            h.ListReturns,
		"list_history":            func(context.Context, []string) { h.ListHistory(ctx) },
		"bulk_accept_orders":      h.BulkAcceptOrders,
	}
	for {
		fmt.Print("> ")

		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		if input == "" {
			fmt.Println("Некорректный запрос")

			continue
		}

		args := strings.Fields(input)
		command := args[0]
		handler, ok := commandHandlers[command]
		if ok {
			handler(ctx, args[1:])

			continue
		}

		if command == "exit" {
			fmt.Println("Завершение работы...")

			return nil
		}
		fmt.Println("Неизвестная команда. Введите 'help' для просмотра доступных команд.")
	}
}
