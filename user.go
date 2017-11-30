package dosbot

type User interface {
	Name() string
	ID() interface{}
}
