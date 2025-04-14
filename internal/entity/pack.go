package entity

// PackType - тип упаковки
type PackType string

const (
	// PackTypePlastBag пластиковый пакет
	PackTypePlastBag PackType = "bag"
	// PackTypeBox коробка
	PackTypeBox PackType = "box"
	// PackTypeFilm пленка
	PackTypeFilm PackType = "film"
)

// Pack содержит информацию об упаковки
// Она содержит 2 поля:
// - PackType: тип упаковки
// - ExtraPack: тип дополнительной упаковки
type Pack struct {
	PackType  PackType `json:"packtype"`
	ExtraPack PackType `json:"extrapack,omitempty"`
}

// NewPack новый экземляр упаковки
func NewPack(packType PackType, extraPack PackType) Pack {
	return Pack{PackType: packType, ExtraPack: extraPack}
}
