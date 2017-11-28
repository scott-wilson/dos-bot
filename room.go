package dosbot

type Room interface {
	Name() string
	ID() interface{}
}
