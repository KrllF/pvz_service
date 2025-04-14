package cli

import (
	"context"
	"fmt"
	"time"

	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/handler/cli/helper"
)

const (
	requiredNumberOfParamAccOrd         = 7
	requiredNumberOfParamAccOrdWithFilm = 8
	mskTime                             = 3 * 60 * 60
)

func (h *clihandler) AcceptOrder(ctx context.Context, args []string) {
	if len(args) != requiredNumberOfParamAccOrd && len(args) != requiredNumberOfParamAccOrdWithFilm {
		fmt.Println("Неверное количество аргументов.")
		fmt.Println("Используйте: accept_order <order_id> <user_id> <weight> <price> <packType>")
		fmt.Println("						    	<shelflife(2025-05-05 20:20:20)> [extraFilm]")

		return
	}

	req, err := parseAndValidateParams(args)
	if err != nil {
		fmt.Printf("Ошибка при обработке параметров: %v\n", err)

		return
	}

	err = h.acceptServ.AcceptOrder(ctx, req)
	if err != nil {
		fmt.Printf("Ошибка при принятии заказа: %v\n", err)

		return
	}

	fmt.Println("Заказ успешно принят")
}

func parseAndValidateParams(args []string) (entity.OrderEntry, error) {
	extraFilm := extra(args)

	orderID, err := helper.ParseID(args[0])
	if err != nil {
		return entity.OrderEntry{}, fmt.Errorf("ошибка при парсинге order_id: %w", err)
	}

	userID, err := helper.ParseID(args[1])
	if err != nil {
		return entity.OrderEntry{}, fmt.Errorf("ошибка при парсинге user_id: %w", err)
	}

	weight, err := helper.ParseID(args[2])
	if err != nil {
		return entity.OrderEntry{}, fmt.Errorf("ошибка при парсинге weight: %w", err)
	}

	price, err := helper.ParseID(args[3])
	if err != nil {
		return entity.OrderEntry{}, fmt.Errorf("ошибка при парсинге price: %w", err)
	}

	if err = validateWeightPrice(weight, price); err != nil {
		return entity.OrderEntry{}, fmt.Errorf("недопустимые значения веса или цены: %w", err)
	}

	packType := entity.PackType(args[4])
	if !isValidPackType(packType) {
		return entity.OrderEntry{}, fmt.Errorf("некорректный тип упаковки: %s", packType)
	}

	extraPack := entity.PackType(extraFilm)
	if extraFilm != "" && !isValidExtraPack(extraPack) {
		return entity.OrderEntry{}, fmt.Errorf("некорректный тип дополнительной упаковки: %s", extraPack)
	}

	layout := "2006-01-02 15:04:05"

	loc := time.FixedZone("Europe/Moscow", mskTime)

	moscowTime, err := time.ParseInLocation(layout, args[5]+" "+args[6], loc)
	if err != nil {
		return entity.OrderEntry{}, fmt.Errorf("ошибка при парсинге shelflife: %w", err)
	}

	return entity.OrderEntry{
		OrderID:   orderID,
		UserID:    userID,
		Weight:    weight,
		Price:     price,
		PackType:  packType,
		ExtraPack: entity.PackType(extraFilm),
		Shelflife: moscowTime,
	}, nil
}

func validateWeightPrice(weight, price int64) error {
	if weight < 0 {
		return fmt.Errorf("вес заказа меньше 0: %d", weight)
	}
	if price < 0 {
		return fmt.Errorf("цена заказа меньше 0: %d", price)
	}

	return nil
}

func isValidPackType(packType entity.PackType) bool {
	switch packType {
	case entity.PackTypePlastBag, entity.PackTypeBox, entity.PackTypeFilm:
		return true
	default:
		return false
	}
}

func isValidExtraPack(extraPack entity.PackType) bool {
	switch extraPack {
	case entity.PackTypeFilm:
		return true
	default:
		return false
	}
}

func extra(args []string) string {
	if len(args) == requiredNumberOfParamAccOrdWithFilm {
		return args[requiredNumberOfParamAccOrdWithFilm-1]
	}

	return ""
}
