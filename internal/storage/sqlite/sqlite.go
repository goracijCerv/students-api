package sqlite

import (
	"database/sql"
	"fmt"
	"log/slog"

	"github.com/goracijCerv/students-api/internal/config"
	"github.com/goracijCerv/students-api/internal/types"
	_ "github.com/mattn/go-sqlite3"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite3", cfg.StoragePath)
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(` CREATE TABLE IF NOT EXISTS students (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	name TEXT,
	lastname TEXT,
	email TEXT,
	number INTEGER,
	age INTEGER
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{
		Db: db,
	}, nil
}

func (s *Sqlite) CreateStudent(name string, lastname string, email string, number int, age int) (int64, error) {
	stmt, err := s.Db.Prepare(`INSERT INTO students (name, lastname, email, number, age) VALUES (? ,?,?,?,?)`)
	if err != nil {
		return 0, nil
	}
	defer stmt.Close()

	result, err := stmt.Exec(name, lastname, email, number, age)

	if err != nil {
		return 0, err
	}

	lastid, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastid, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare(`SELECT * FROM students WHERE id = ? LIMIT 1`)
	if err != nil {
		return types.Student{}, err
	}

	defer stmt.Close()
	var student types.Student
	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.LastName, &student.Email, &student.Number, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found")
		}

		return types.Student{}, fmt.Errorf("error de query: %w", err)
	}

	return student, nil

}

func (s *Sqlite) GetAllStudents() ([]types.Student, error) {
	smtp, err := s.Db.Prepare(`SELECT * FROM students`)
	if err != nil {
		return nil, err
	}
	defer smtp.Close()

	rows, err := smtp.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student
		err := rows.Scan(&student.Id, &student.Name, &student.LastName, &student.Email, &student.Number, &student.Age)
		if err != nil {
			slog.Error("Something went wrong in the  scan")
			return nil, err
		}
		students = append(students, student)

	}

	return students, nil

}

func (s *Sqlite) UpdateStudent(id int64, name string, lastname string, email string, number int, age int) error {
	smtp, err := s.Db.Prepare(
		`UPDATE students
		SET name = ?,
	    lastname = ?,
	    email = ?,
		number = ?,
		age = ?
    	WHERE id = ?`)

	if err != nil {
		return err
	}

	defer smtp.Close()

	result, err := smtp.Exec(name, lastname, email, number, age, id)

	if err != nil {
		return err
	}

	_, err = result.RowsAffected()

	if err != nil {
		slog.Error("Some error in the query")
		return err
	}

	return nil

}

// https://stackoverflow.com/questions/43580131/exec-gcc-executable-file-not-found-in-path-when-trying-go-build si ocurre un error con la gcc Y este comando go env -w CGO_ENABLED=1
