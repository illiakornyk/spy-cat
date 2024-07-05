package common

type Mission struct {
	ID       int64
	CatID    int64
	Complete bool
	Targets  []Target
}
