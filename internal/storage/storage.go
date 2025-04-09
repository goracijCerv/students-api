package storage

type Storage interface {
	CreateStudent(name string, lastname string, email string, number int, age int) (int64, error)
}
