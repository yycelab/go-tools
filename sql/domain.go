package sql

type Pagable interface {
	Current() int
	Size() int
	SetCurrent(current int)
	SetTotal(total int)
}
