package converter

import (
	"github.com/KrllF/pvz_service/internal/entity"
	"github.com/KrllF/pvz_service/internal/models"
	desc "github.com/KrllF/pvz_service/pkg/order_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// ConvertToProtoOrder из models.Order в proto
func ConvertToProtoOrder(ord *models.Order) *desc.Order {
	if ord == nil {
		return nil
	}

	inStorageFrom := timestamppb.New(ord.InStorageFrom)
	shelfLife := timestamppb.New(ord.ShelfLife)
	twoDaysOfLife := timestamppb.New(ord.TwoDaysOfLife)
	lastUpdate := timestamppb.New(ord.LastUpdate)

	packaging := &desc.Pack{
		Packtype:  *ConvertPackTypeToProto(ord.Packaging.PackType),
		Extrapack: ConvertPackTypeToProto(ord.Packaging.ExtraPack),
	}

	return &desc.Order{
		Orderid:       ord.OrderID,
		Userid:        ord.UserID,
		Weight:        ord.Weight,
		Price:         ord.Price,
		Packaging:     packaging,
		Instoragefrom: inStorageFrom,
		Shelflife:     shelfLife,
		Status:        ord.Status,
		Twodaysoflife: twoDaysOfLife,
		Lastupdate:    lastUpdate,
	}
}

// ConvertPackType из прото в packtype
func ConvertPackType(protoPackType int64) entity.PackType {
	switch protoPackType {
	case int64(desc.PackType_bag):
		return entity.PackTypePlastBag
	case int64(desc.PackType_box):
		return entity.PackTypeBox
	case int64(desc.PackType_film):
		return entity.PackTypeFilm
	default:
		return ""
	}
}

// ConvertPackTypeToProto из packtype в прото
func ConvertPackTypeToProto(packType entity.PackType) *desc.PackType {
	switch packType {
	case entity.PackTypePlastBag:
		pack := desc.PackType_bag

		return &pack

	case entity.PackTypeBox:
		pack := desc.PackType_box

		return &pack

	case entity.PackTypeFilm:
		pack := desc.PackType_film

		return &pack

	default:
		pack := desc.PackType_without

		return &pack
	}
}
