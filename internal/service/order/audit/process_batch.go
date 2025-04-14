package audite

// ProcessBatch запись батча в каналы
func (a *Audit) ProcessBatch(batch []interface{}) {
	for _, val := range batch {
		a.DB <- val
		a.Stdout <- val
	}
}
