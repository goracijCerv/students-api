package student_test

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/goracijCerv/students-api/internal/http/handlers/student"
	"github.com/goracijCerv/students-api/internal/types"
	"github.com/goracijCerv/students-api/internal/utils/test"
)

type mockStorage struct {
	allExpected []types.Student
	expected    types.Student
	returnID    int64
	err         error
	called      bool
}

func (m *mockStorage) CreateStudent(name string, lastname string, email string, number int, age int) (int64, error) {
	m.called = true
	if m.err != nil {
		return 0, m.err
	}
	return m.returnID, nil
}

func (m *mockStorage) GetStudentById(id int64) (types.Student, error) {
	m.called = true
	if m.err != nil {
		return types.Student{}, m.err
	}
	return m.expected, nil
}

func (m *mockStorage) GetAllStudents() ([]types.Student, error) {
	m.called = true
	if m.err != nil {
		return []types.Student{}, m.err
	}

	return m.allExpected, nil
}

func (m *mockStorage) DeleteStudentById(id int64) error {
	m.called = true
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockStorage) UpdateStudent(id int64, name string, lastname string, email string, number int, age int) error {
	m.called = true
	if m.err != nil {
		return m.err
	}
	return nil
}

func (m *mockStorage) GetEmailTemplate(name string) (string, error) {
	m.called = true
	if m.err != nil {
		return "", m.err
	}

	return "hola", nil
}

func TestWelcomeHandler(t *testing.T) {
	req := test.NewRequest(t, http.MethodGet, "/", nil)
	rr := test.RunHandler(t, student.Welcome(), req)

	res := rr.Result()
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		t.Errorf("se espera un status 200. se obtubo un %d", res.StatusCode)
	}

	body, _ := io.ReadAll(res.Body)
	expected := "Bienvenido a la students api"
	if !strings.Contains(string(body), expected) {
		t.Errorf("Se esperaba que el body contubiera %q, se obtubo %q", expected, string(body))
	}
}

func TestNewHandler_Success(t *testing.T) {
	studentData := types.Student{
		Id:       0,
		Name:     "John",
		LastName: "Doe",
		Email:    "john@example.com",
		Number:   12345,
		Age:      21,
	}

	mock := &mockStorage{expected: studentData, returnID: 1}
	req := test.NewRequest(t, http.MethodPost, "/api/student", studentData)
	rr := test.RunHandler(t, student.New(mock), req)
	if rr.Code != http.StatusCreated {
		t.Errorf("se esperaba 201 se obtubo %d", rr.Code)
	}

	if !mock.called {
		t.Error("CreateStudent not called")
	}
}

func TestNewHandler_EmptyBody(t *testing.T) {
	req := test.NewRequest(t, http.MethodPost, "/api/student", nil)
	rr := test.RunHandler(t, student.New(&mockStorage{}), req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestNewHandler_InvalidJSON(t *testing.T) {
	req := test.NewRequest(t, http.MethodPost, "/api/student", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	rr := test.RunHandler(t, student.New(&mockStorage{}), req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestNewHandler_ValidateError(t *testing.T) {
	studentData := types.Student{
		Id:       0,
		Name:     "John",
		LastName: "Doe",
		Email:    "",
		Number:   12345,
		Age:      21,
	}
	req := test.NewRequest(t, http.MethodPost, "/api/student", studentData)
	rr := test.RunHandler(t, student.New(&mockStorage{}), req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}
}

func TestNewHandler_StorageError(t *testing.T) {
	studentData := types.Student{
		Id:       0,
		Name:     "Jane",
		LastName: "Smith",
		Email:    "jane@example.com",
		Number:   67890,
		Age:      22,
	}
	mock := &mockStorage{expected: studentData, err: errors.New("db fail")}
	req := test.NewRequest(t, http.MethodPost, "/api/student", studentData)
	rr := test.RunHandler(t, student.New(mock), req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}
}

func TestGetByIdHandler_Succes(t *testing.T) {
	studentData := types.Student{
		Id:       1,
		Name:     "Jane",
		LastName: "Smith",
		Email:    "jane@example.com",
		Number:   67890,
		Age:      22,
	}
	mock := &mockStorage{expected: studentData}
	req := test.NewRequest(t, http.MethodGet, "/api/student/1", nil)
	req.SetPathValue("id", "1")
	rr := test.RunHandler(t, student.GetById(mock), req)

	if !mock.called {
		t.Errorf("not called")
	}

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 , got %d", rr.Code)
	}
}

func TestGetByIdHandler_NotValidPath(t *testing.T) {
	req := test.NewRequest(t, http.MethodGet, "/api/student/1", nil)
	req.SetPathValue("id", "ac")
	rr := test.RunHandler(t, student.GetById(&mockStorage{}), req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400 , got %d", rr.Code)
	}
}

func TestGetByIdHanlder_NotFound(t *testing.T) {
	studentData := types.Student{
		Id:       1,
		Name:     "Jane",
		LastName: "Smith",
		Email:    "jane@example.com",
		Number:   67890,
		Age:      22,
	}
	mock := &mockStorage{expected: studentData, err: errors.New("no student found")}
	req := test.NewRequest(t, http.MethodGet, "/api/student/1", nil)
	req.SetPathValue("id", "1")
	rr := test.RunHandler(t, student.GetById(mock), req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected 404 , got %d", rr.Code)
	}
}

func TestGetByIdHanlder_StrorageError(t *testing.T) {
	studentData := types.Student{
		Id:       1,
		Name:     "Jane",
		LastName: "Smith",
		Email:    "jane@example.com",
		Number:   67890,
		Age:      22,
	}
	mock := &mockStorage{expected: studentData, err: errors.New("storage error")}
	req := test.NewRequest(t, http.MethodGet, "/api/student/1", nil)
	req.SetPathValue("id", "1")
	rr := test.RunHandler(t, student.GetById(mock), req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 , got %d", rr.Code)
	}
}

func TestGetAll_Success(t *testing.T) {
	studentData := []types.Student{
		{Id: 1,
			Name:     "Jane",
			LastName: "Smith",
			Email:    "jane@example.com",
			Number:   67890,
			Age:      22},
		{Id: 2,
			Name:     "Julia",
			LastName: "Sanders",
			Email:    "julia@example.com",
			Number:   67890,
			Age:      24},
	}

	mock := &mockStorage{allExpected: studentData}
	req := test.NewRequest(t, http.MethodGet, "/api/student", nil)
	rr := test.RunHandler(t, student.GetListStudents(mock), req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200 , got %d", rr.Code)
	}

	if !mock.called {
		t.Errorf("not called")
	}

}

func TestGetAll_StorageError(t *testing.T) {
	studentData := []types.Student{
		{Id: 1,
			Name:     "Jane",
			LastName: "Smith",
			Email:    "jane@example.com",
			Number:   67890,
			Age:      22},
		{Id: 2,
			Name:     "Julia",
			LastName: "Sanders",
			Email:    "julia@example.com",
			Number:   67890,
			Age:      24},
	}

	mock := &mockStorage{allExpected: studentData, err: errors.New("Error in db")}
	req := test.NewRequest(t, http.MethodGet, "/api/student", nil)
	rr := test.RunHandler(t, student.GetListStudents(mock), req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500 , got %d", rr.Code)
	}

}

func TestUpdateById_Succes(t *testing.T) {
	studentData := types.Student{
		Id:       1,
		Name:     "Janew",
		LastName: "Smith",
		Email:    "jane@example.com",
		Number:   67890,
		Age:      22,
	}
	mock := &mockStorage{expected: studentData}
	req := test.NewRequest(t, http.MethodPut, "/api/student/1", studentData)
	req.SetPathValue("id", "1")
	rr := test.RunHandler(t, student.UpdateById(mock), req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rr.Code)
	}

	if !mock.called {
		t.Errorf("not called")
	}

}

func TestUpdateById_NotValidPath(t *testing.T) {
	studentData := types.Student{
		Id:       1,
		Name:     "Janew",
		LastName: "Smith",
		Email:    "jane@example.com",
		Number:   67890,
		Age:      22,
	}
	mock := &mockStorage{expected: studentData}
	req := test.NewRequest(t, http.MethodPut, "/api/student/", studentData)
	req.SetPathValue("id", "a")
	rr := test.RunHandler(t, student.UpdateById(mock), req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

}

func TestUpdateById_EmptyBody(t *testing.T) {

	req := test.NewRequest(t, http.MethodPut, "/api/student/", nil)
	req.SetPathValue("id", "1")
	rr := test.RunHandler(t, student.UpdateById(&mockStorage{}), req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

}

func TestUpdateById_ValidateErros(t *testing.T) {

	studentData := types.Student{
		Id:       1,
		Name:     "Janew",
		LastName: "Smith",
		Email:    "",
		Number:   67890,
		Age:      22,
	}
	mock := &mockStorage{expected: studentData}
	req := test.NewRequest(t, http.MethodPut, "/api/student/", studentData)
	req.SetPathValue("id", "1")
	rr := test.RunHandler(t, student.UpdateById(mock), req)

	if rr.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d", rr.Code)
	}

}

func TestUpdateById_StorageErros(t *testing.T) {

	studentData := types.Student{
		Id:       1,
		Name:     "Janew",
		LastName: "Smith",
		Email:    "janew@email.com",
		Number:   67890,
		Age:      22,
	}
	mock := &mockStorage{expected: studentData, err: errors.New("db error")}
	req := test.NewRequest(t, http.MethodPut, "/api/student/", studentData)
	req.SetPathValue("id", "1")
	rr := test.RunHandler(t, student.UpdateById(mock), req)

	if rr.Code != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", rr.Code)
	}

}
