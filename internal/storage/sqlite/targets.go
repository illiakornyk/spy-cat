package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/illiakornyk/spy-cat/internal/common"
)

func (s *Storage) UpdateTarget(id int64, target common.Target) error {
	const op = "storage.sqlite.UpdateTarget"

	var complete bool
	err := s.db.QueryRow("SELECT complete FROM targets WHERE id = ?", id).Scan(&complete)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: target not found", op)
		}
		return fmt.Errorf("%s: query target: %w", op, err)
	}
	if complete {
		return fmt.Errorf("%s: cannot update a completed target", op)
	}

	stmt, err := s.db.Prepare("UPDATE targets SET name = ?, country = ?, notes = ?, complete = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(target.Name, target.Country, target.Notes, target.Complete, id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}


func (s *Storage) UpdateCompleteStatus(targetID int64, complete bool) error {
	const op = "storage.sqlite.UpdateCompleteStatus"

	// Update the complete status
	stmt, err := s.db.Prepare("UPDATE targets SET complete = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(complete, targetID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *Storage) UpdateNotes(targetID int64, notes string) error {
	const op = "storage.sqlite.UpdateNotes"

	var targetComplete, missionComplete bool
	err := s.db.QueryRow("SELECT t.complete, m.complete FROM targets t JOIN missions m ON t.mission_id = m.id WHERE t.id = ?", targetID).Scan(&targetComplete, &missionComplete)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: target not found", op)
		}
		return fmt.Errorf("%s: query target and mission: %w", op, err)
	}
	if targetComplete || missionComplete {
		return fmt.Errorf("%s: cannot update notes for a completed target or mission", op)
	}

	stmt, err := s.db.Prepare("UPDATE targets SET notes = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(notes, targetID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}


func (s *Storage) TargetExists(targetID int64) (bool, error) {
	const op = "storage.sqlite.TargetExists"

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM targets WHERE id = ?)", targetID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: query row: %w", op, err)
	}

	return exists, nil
}



func (s *Storage) DeleteTarget(targetID int64) error {
	const op = "storage.sqlite.DeleteTarget"

	var complete bool
	err := s.db.QueryRow("SELECT complete FROM targets WHERE id = ?", targetID).Scan(&complete)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: target not found", op)
		}
		return fmt.Errorf("%s: query target: %w", op, err)
	}
	if complete {
		return fmt.Errorf("%s: cannot delete a completed target", op)
	}

	stmt, err := s.db.Prepare("DELETE FROM targets WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(targetID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}


func (s *Storage) AddTarget(missionID int64, name, country, notes string) (int64, error) {
    const op = "storage.sqlite.AddTarget"

    var complete bool
    err := s.db.QueryRow("SELECT complete FROM missions WHERE id = ?", missionID).Scan(&complete)
    if err != nil {
        if err == sql.ErrNoRows {
            return 0, fmt.Errorf("%s: mission not found", op)
        }
        return 0, fmt.Errorf("%s: query mission: %w", op, err)
    }
    if complete {
        return 0, fmt.Errorf("%s: cannot add target to a completed mission", op)
    }

    // Check the current number of targets in the mission
    count, err := s.getTargetCountForMission(missionID)
    if err != nil {
        return 0, fmt.Errorf("%s: %w", op, err)
    }
    if count >= 3 {
        return 0, fmt.Errorf("%s: mission already has the maximum number of targets (3)", op)
    }

    stmt, err := s.db.Prepare("INSERT INTO targets (mission_id, name, country, notes, complete) VALUES (?, ?, ?, ?, 0)")
    if err != nil {
        return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
    }
    defer stmt.Close()

    res, err := stmt.Exec(missionID, name, country, notes)
    if err != nil {
        return 0, fmt.Errorf("%s: execute statement: %w", op, err)
    }

    targetID, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
    }

    return targetID, nil
}
