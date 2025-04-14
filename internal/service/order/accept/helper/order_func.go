package helper

import (
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/errs"
	"github.com/KrllF/pvz_service/internal/models"
)

const (
	plastBagLim = 10_000 // в граммах
	boxLim      = 30_000 // в граммах

	plastBagPrice = 5
	boxPrice      = 20
	filmPrice     = 1

	extraFilmPrice = 1
)

// Validate проверка веса
func Validate(o models.Order) error {
	switch o.Packaging.PackType {
	case entity.PackTypePlastBag:
		if o.Weight >= plastBagLim {
			return errs.ErrInvalidData
		}
	case entity.PackTypeBox:
		if o.Weight >= boxLim {
			return errs.ErrInvalidData
		}
	case entity.PackTypeFilm:
	default:
		return errs.ErrPackTypeNotFound
	}

	return nil
}

// CalcPrice подсчёт стоимости заказа
func CalcPrice(o models.Order) int64 {
	price := o.Price

	switch o.Packaging.PackType {
	case entity.PackTypePlastBag:
		price += plastBagPrice
	case entity.PackTypeBox:
		price += boxPrice
	case entity.PackTypeFilm:
		price += filmPrice
	}

	if o.Packaging.ExtraPack == entity.PackTypeFilm {
		price += extraFilmPrice
	}

	return price
}
