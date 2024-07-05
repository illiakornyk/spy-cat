package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

type Storage struct {
    db *sql.DB
}

func New(storagePath string) (*Storage, error) {
    const op = "storage.sqlite.NewStorage"

    log.Printf("Opening SQLite database at path: %s", storagePath)

    dir := filepath.Dir(storagePath)
    if err := os.MkdirAll(dir, 0755); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    db, err := sql.Open("sqlite3", storagePath)
    if err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    if err := runMigrations(db); err != nil {
        return nil, fmt.Errorf("%s: %w", op, err)
    }

    return &Storage{db: db}, nil
}

func runMigrations(db *sql.DB) error {
    driver, err := sqlite3.WithInstance(db, &sqlite3.Config{})
    if err != nil {
        return fmt.Errorf("could not create migration driver: %w", err)
    }

    m, err := migrate.NewWithDatabaseInstance(
        "file://migrations",
        "sqlite3", driver)
    if err != nil {
        return fmt.Errorf("could not start migration: %w", err)
    }

    if err := m.Up(); err != nil && err != migrate.ErrNoChange {
        return fmt.Errorf("could not apply migration: %w", err)
    }

    return nil
}

func (s *Storage) CreateCat(name string, yearsOfExperience int, breed string, salary float64) (int64, error) {
    const op = "storage.sqlite.SaveCat"

    stmt, err := s.db.Prepare("INSERT INTO spy_cats (name, years_of_experience, breed, salary) VALUES (?, ?, ?, ?)")
    if err != nil {
        return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
    }
    defer stmt.Close()

    res, err := stmt.Exec(name, yearsOfExperience, breed, salary)
    if err != nil {
        if isConstraintViolation(err) {
            return 0, fmt.Errorf("%s: %w", op, errors.New("cat already exists"))
        }
        return 0, fmt.Errorf("%s: execute statement: %w", op, err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
    }

    return id, nil
}


func (s *Storage) DeleteCat(id int64) error {
    const op = "storage.sqlite.DeleteCat"

    stmt, err := s.db.Prepare("DELETE FROM spy_cats WHERE id = ?")
    if err != nil {
        return fmt.Errorf("%s: prepare statement: %w", op, err)
    }
    defer stmt.Close()

    res, err := stmt.Exec(id)
    if err != nil {
        return fmt.Errorf("%s: execute statement: %w", op, err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("%s: get rows affected: %w", op, err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("%s: no cat found with id %d", op, id)
    }

    return nil
}

func (s *Storage) UpdateCatSalary(id int64, salary float64) error {
    const op = "storage.sqlite.UpdateCatSalary"

    stmt, err := s.db.Prepare("UPDATE spy_cats SET salary = ? WHERE id = ?")
    if err != nil {
        return fmt.Errorf("%s: prepare statement: %w", op, err)
    }
    defer stmt.Close()

    res, err := stmt.Exec(salary, id)
    if err != nil {
        return fmt.Errorf("%s: execute statement: %w", op, err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        return fmt.Errorf("%s: get rows affected: %w", op, err)
    }

    if rowsAffected == 0 {
        return fmt.Errorf("%s: no cat found with id %d", op, id)
    }

    return nil
}

func (s *Storage) CatExists(id int64) (bool, error) {
    const op = "storage.sqlite.CatExists"

    var exists bool
    err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM spy_cats WHERE id = ?)", id).Scan(&exists)
    if err != nil {
        return false, fmt.Errorf("%s: query row: %w", op, err)
    }

    return exists, nil
}

func isConstraintViolation(err error) bool {
    return strings.Contains(err.Error(), "UNIQUE constraint failed")
}
