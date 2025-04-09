package storage

import "github.com/goracijCerv/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, lastname string, email string, number int, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
}
