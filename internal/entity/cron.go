package entity

// ResponceBDLog для сообщения в кафку
type ResponceBDLog struct {
	ID       int64
	LG       []byte
	Status   string
	Attempts int64
}
