package audite

// Close закрыть все каналы
func (a *Audit) Close() {
	close(a.DB)
	close(a.HTTP)
	close(a.Status)
	close(a.Stdout)
}
