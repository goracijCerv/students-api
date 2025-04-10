package storage

import "github.com/goracijCerv/students-api/internal/types"

type Storage interface {
	CreateStudent(name string, lastname string, email string, number int, age int) (int64, error)
	GetStudentById(id int64) (types.Student, error)
	GetAllStudents() ([]types.Student, error)
	UpdateStudent(id int64, name string, lastname string, email string, number int, age int) error
	DeleteStudentById(id int64) error
}
