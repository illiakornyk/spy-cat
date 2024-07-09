package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/illiakornyk/spy-cat/internal/http-server/handlers/spycat"
)



func isConstraintViolation(err error) bool {
    return strings.Contains(err.Error(), "UNIQUE constraint failed")
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

    // Begin transaction
    tx, err := s.db.Begin()
    if err != nil {
        return fmt.Errorf("%s: begin transaction: %w", op, err)
    }

    // Get all missions associated with the cat
    missions, err := s.getMissionsByCatID(tx, id)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("%s: get missions by cat ID: %w", op, err)
    }

    // Delete each mission within the same transaction
    for _, missionID := range missions {
        err = s.deleteMissionTx(tx, missionID, true)
        if err != nil {
            tx.Rollback()
            return fmt.Errorf("%s: delete mission: %w", op, err)
        }
    }

    // Delete the cat
    stmt, err := tx.Prepare("DELETE FROM spy_cats WHERE id = ?")
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("%s: prepare statement: %w", op, err)
    }
    defer stmt.Close()

    res, err := stmt.Exec(id)
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("%s: execute statement: %w", op, err)
    }

    rowsAffected, err := res.RowsAffected()
    if err != nil {
        tx.Rollback()
        return fmt.Errorf("%s: get rows affected: %w", op, err)
    }

    if rowsAffected == 0 {
        tx.Rollback()
        return fmt.Errorf("%s: no cat found with id %d", op, id)
    }

    // Commit transaction
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("%s: commit transaction: %w", op, err)
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

func (s *Storage) GetAllCats() ([]spycat.SpyCat, error) {
	const op = "storage.sqlite.GetAllCats"

	rows, err := s.db.Query("SELECT id, name, years_of_experience, breed, salary FROM spy_cats")
	if err != nil {
		return nil, fmt.Errorf("%s: query: %w", op, err)
	}
	defer rows.Close()

	var cats []spycat.SpyCat
	for rows.Next() {
		var cat spycat.SpyCat
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary); err != nil {
			return nil, fmt.Errorf("%s: scan: %w", op, err)
		}
		cats = append(cats, cat)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return cats, nil
}

func (s *Storage) GetCatByID(id int64) (*spycat.SpyCat, error) {
	const op = "storage.sqlite.GetCatByID"

	var cat spycat.SpyCat
	err := s.db.QueryRow("SELECT id, name, years_of_experience, breed, salary FROM spy_cats WHERE id = ?", id).
		Scan(&cat.ID, &cat.Name, &cat.YearsOfExperience, &cat.Breed, &cat.Salary)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query row: %w", op, err)
	}

	return &cat, nil
}


func (s *Storage) getMissionsByCatID(tx *sql.Tx, catID int64) ([]int64, error) {
	const op = "storage.sqlite.getMissionsByCatID"

	rows, err := tx.Query("SELECT id FROM missions WHERE cat_id = ?", catID)
	if err != nil {
		return nil, fmt.Errorf("%s: query missions: %w", op, err)
	}
	defer rows.Close()

	var missions []int64
	for rows.Next() {
		var missionID int64
		if err := rows.Scan(&missionID); err != nil {
			return nil, fmt.Errorf("%s: scan mission ID: %w", op, err)
		}
		missions = append(missions, missionID)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return missions, nil
}
