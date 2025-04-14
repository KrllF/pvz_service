package main

import (
	"context"
	"log"

	"github.com/KrllF/pvz_service/internal/app"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	a, err := app.NewApp(ctx)
	if err != nil {
		log.Println("не получилось запустить приложение: ", err)

		return
	}

	if err := a.Run(); err != nil {
		log.Panic("ошибка при выполнении: ", err.Error())
	}
}
