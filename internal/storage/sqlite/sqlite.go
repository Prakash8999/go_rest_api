package sqlite

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/prakash8999/go_rest_apis/internal/types"

	"github.com/prakash8999/go_rest_apis/internal/config"

	//beacuse we are using the database/sql package, we need to import the driver for the database we are using however it used under the hood not in the project directly that's why we used _ to ignore the unused import, but we just need it for initialization
	_ "modernc.org/sqlite"
)

type Sqlite struct {
	Db *sql.DB
}

func New(cfg *config.Config) (*Sqlite, error) {
	db, err := sql.Open("sqlite", cfg.StoragePath)

	if err != nil {
		return nil, err
	}

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS students (
	 id INTEGER PRIMARY KEY AUTOINCREMENT,
	 name TEXT NOT NULL,
	 email TEXT NOT NULL,
	 age INTEGER NOT NULL
	 
	
	)`)

	if err != nil {
		return nil, err
	}

	return &Sqlite{Db: db}, nil
}

func (s *Sqlite) CreateStudent(name string, email string, age int) (int64, error) {
	stmt, err := s.Db.Prepare("INSERT INTO students (name,email,age) VALUES (?,?,?)")
	if err != nil {
		return 0, err
	}

	defer stmt.Close()

	result, err := stmt.Exec(name, email, age)
	if err != nil {
		return 0, err
	}

	lastId, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return lastId, nil
}

func (s *Sqlite) GetStudentById(id int64) (types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id,name,email,age FROM students WHERE id = ? LIMIT 1")

	if err != nil {
		return types.Student{}, err
	}
	defer stmt.Close()

	var student types.Student

	err = stmt.QueryRow(id).Scan(&student.Id, &student.Name, &student.Email, &student.Age)
	if err != nil {
		if err == sql.ErrNoRows {
			return types.Student{}, fmt.Errorf("no student found with id %s", fmt.Sprint(id))
		}

		return types.Student{}, fmt.Errorf("query error %w", err)

	}
	// fmt.Println(fmt.Sprint(student))
	return student, nil

}

func (s *Sqlite) GetStudents() ([]types.Student, error) {
	stmt, err := s.Db.Prepare("SELECT id, name, email, age FROM students")
	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var students []types.Student

	for rows.Next() {
		var student types.Student

		err := rows.Scan(&student.Id, &student.Name, &student.Email, &student.Age)
		if err != nil {
			return nil, err
		}

		students = append(students, student)
	}

	return students, nil
}

func (s *Sqlite) UpdateStudent(update types.UpdateStudentRequest) (string, error) {
	query := "UPDATE students SET "
	args := []interface{}{}
	setParts := []string{}

	if update.Name != nil {
		setParts = append(setParts, "name = ?")
		args = append(args, *update.Name)
	}
	if update.Email != nil {
		setParts = append(setParts, "email = ?")
		args = append(args, *update.Email)
	}
	if update.Age != nil {
		setParts = append(setParts, "age = ?")
		args = append(args, *update.Age)
	}

	if len(setParts) == 0 {
		return "", fmt.Errorf("no fields to update")
	}

	query += strings.Join(setParts, ", ") + " WHERE id = ?"
	args = append(args, update.Id)

	stmt, err := s.Db.Prepare(query)
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	_, err = stmt.Exec(args...)
	if err != nil {
		return "", err
	}

	return "student updated successfully", nil
}

func (s *Sqlite) DeleteStudentById(id int64) (string, error) {
	// First, get the student to return after deletion

	// Prepare delete statement
	stmt, err := s.Db.Prepare("DELETE FROM students WHERE id = ?")
	if err != nil {
		return "Fail to delete data", fmt.Errorf("failed to prepare delete statement: %w", err)
	}
	defer stmt.Close()

	// Execute deletion
	_, err = stmt.Exec(id)
	if err != nil {
		return "Fail to delete data", fmt.Errorf("failed to delete student: %w", err)
	}

	return "User deleted successfully", nil
}
