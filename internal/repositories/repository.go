package repositories

type Repository[T any] interface {
	GetById(id string) (T, error)
	Save(*T) error
}
