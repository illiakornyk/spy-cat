package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/illiakornyk/spy-cat/internal/common"
)






func (s *Storage) CreateMission(catID int64, targets []common.Target, complete bool) (int64, error) {
	const op = "storage.sqlite.CreateMission"

	tx, err := s.db.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: begin transaction: %w", op, err)
	}

	stmt, err := tx.Prepare("INSERT INTO missions (cat_id, complete) VALUES (?, ?)")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(catID, complete)
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	missionID, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: last insert id: %w", op, err)
	}

	for _, target := range targets {
		targetStmt, err := tx.Prepare("INSERT INTO targets (mission_id, name, country, notes, complete) VALUES (?, ?, ?, ?, ?)")
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: prepare target statement: %w", op, err)
		}
		defer targetStmt.Close()

		_, err = targetStmt.Exec(missionID, target.Name, target.Country, target.Notes, target.Complete)
		if err != nil {
			tx.Rollback()
			return 0, fmt.Errorf("%s: execute target statement: %w", op, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("%s: commit transaction: %w", op, err)
	}

	return missionID, nil
}


func (s *Storage) GetAllMissions() ([]common.Mission, error) {
    const op = "storage.sqlite.GetAllMissions"

    rows, err := s.db.Query("SELECT id, cat_id, complete FROM missions")
    if err != nil {
        return nil, fmt.Errorf("%s: query: %w", op, err)
    }
    defer rows.Close()

    var missions []common.Mission
    for rows.Next() {
        var mission common.Mission
        if err := rows.Scan(&mission.ID, &mission.CatID, &mission.Complete); err != nil {
            return nil, fmt.Errorf("%s: scan: %w", op, err)
        }

        targetRows, err := s.db.Query("SELECT id, mission_id, name, country, notes, complete FROM targets WHERE mission_id = ?", mission.ID)
        if err != nil {
            return nil, fmt.Errorf("%s: query targets: %w", op, err)
        }
        defer targetRows.Close()

        var targets []common.Target
        for targetRows.Next() {
            var target common.Target
            if err := targetRows.Scan(&target.ID, &target.MissionID, &target.Name, &target.Country, &target.Notes, &target.Complete); err != nil {
                return nil, fmt.Errorf("%s: scan target: %w", op, err)
            }
            targets = append(targets, target)
        }
        mission.Targets = targets

        missions = append(missions, mission)
    }

    if err := rows.Err(); err != nil {
        return nil, fmt.Errorf("%s: rows error: %w", op, err)
    }

    return missions, nil
}




func (s *Storage) DeleteMission(id int64) error {
	const op = "storage.sqlite.DeleteMission"

	var catID int64
	err := s.db.QueryRow("SELECT cat_id FROM missions WHERE id = ?", id).Scan(&catID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: mission not found", op)
		}
		return fmt.Errorf("%s: query: %w", op, err)
	}

	if catID != 0 {
		return fmt.Errorf("%s: cannot delete mission assigned to a cat", op)
	}

	stmt, err := s.db.Prepare("DELETE FROM targets WHERE mission_id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare delete targets statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute delete targets statement: %w", op, err)
	}

	stmt, err = s.db.Prepare("DELETE FROM missions WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare delete mission statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("%s: execute delete mission statement: %w", op, err)
	}

	return nil
}



func (s *Storage) UpdateMissionCompleteStatus(id int64, complete bool) error {
	const op = "storage.sqlite.UpdateMissionCompleteStatus"

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM missions WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: query row: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: mission not found", op)
	}

	if complete {
		var incompleteTargets int
		err = s.db.QueryRow("SELECT COUNT(*) FROM targets WHERE mission_id = ? AND complete = 0", id).Scan(&incompleteTargets)
		if err != nil {
			return fmt.Errorf("%s: query incomplete targets: %w", op, err)
		}
		if incompleteTargets > 0 {
			return fmt.Errorf("%s: cannot complete mission with incomplete targets", op)
		}
	}

	stmt, err := s.db.Prepare("UPDATE missions SET complete = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(complete, id)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}


func (s *Storage) MissionExists(missionID int64) (bool, error) {
	const op = "storage.sqlite.MissionExists"

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM missions WHERE id = ?)", missionID).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: query row: %w", op, err)
	}

	return exists, nil
}
