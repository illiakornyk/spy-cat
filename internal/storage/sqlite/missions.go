package sqlite

import (
	"database/sql"
	"fmt"

	"github.com/illiakornyk/spy-cat/internal/common"
)






func (s *Storage) CreateMission(catID sql.NullInt64, targets []common.Target, complete bool) (int64, error) {
	const op = "storage.sqlite.CreateMission"

	stmt, err := s.db.Prepare("INSERT INTO missions (cat_id, complete) VALUES (?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	res, err := stmt.Exec(catID, complete)
	if err != nil {
		return 0, fmt.Errorf("%s: execute statement: %w", op, err)
	}

	missionID, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	for _, target := range targets {
		_, err := s.AddTarget(missionID, target.Name, target.Country, target.Notes)
		if err != nil {
			return 0, fmt.Errorf("%s: failed to add target: %w", op, err)
		}
	}

	return missionID, nil
}



func (s *Storage) UpdateMissionCompleteStatus(id int64, complete bool) error {
	const op = "storage.sqlite.UpdateMissionCompleteStatus"

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

func (s *Storage) AssignCatToMission(missionID, catID int64) error {
	const op = "storage.sqlite.AssignCatToMission"

	// Check if the mission exists
	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM missions WHERE id = ?)", missionID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: query mission: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: mission does not exist", op)
	}

	// Check if the cat exists
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM spy_cats WHERE id = ?)", catID).Scan(&exists)
	if err != nil {
		return fmt.Errorf("%s: query cat: %w", op, err)
	}
	if !exists {
		return fmt.Errorf("%s: cat does not exist", op)
	}

	// Assign the cat to the mission
	stmt, err := s.db.Prepare("UPDATE missions SET cat_id = ? WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(catID, missionID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
}

func (s *Storage) MissionExists(id int64) (bool, error) {
	const op = "storage.sqlite.MissionExists"

	var exists bool
	err := s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM missions WHERE id = ?)", id).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: query mission: %w", op, err)
	}

	return exists, nil
}



func (s *Storage) DeleteMission(missionID int64) error {
	const op = "storage.sqlite.DeleteMission"

	var catID sql.NullInt64
	err := s.db.QueryRow("SELECT cat_id FROM missions WHERE id = ?", missionID).Scan(&catID)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s: mission not found", op)
		}
		return fmt.Errorf("%s: query mission: %w", op, err)
	}
	if catID.Valid {
		return fmt.Errorf("%s: cannot delete a mission assigned to a cat", op)
	}

	stmt, err := s.db.Prepare("DELETE FROM missions WHERE id = ?")
	if err != nil {
		return fmt.Errorf("%s: prepare statement: %w", op, err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(missionID)
	if err != nil {
		return fmt.Errorf("%s: execute statement: %w", op, err)
	}

	return nil
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

func (s *Storage) GetMission(id int64) (*common.Mission, error) {
const op = "storage.sqlite.GetMissionWithTargets"

	var mission common.Mission
	err := s.db.QueryRow("SELECT id, cat_id, complete FROM missions WHERE id = ?", id).
		Scan(&mission.ID, &mission.CatID, &mission.Complete)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("%s: query row: %w", op, err)
	}

	rows, err := s.db.Query("SELECT id, name, country, notes, complete FROM targets WHERE mission_id = ?", id)
	if err != nil {
		return nil, fmt.Errorf("%s: query targets: %w", op, err)
	}
	defer rows.Close()

	for rows.Next() {
		var target common.Target
		err := rows.Scan(&target.ID, &target.Name, &target.Country, &target.Notes, &target.Complete)
		if err != nil {
			return nil, fmt.Errorf("%s: scan target: %w", op, err)
		}
		mission.Targets = append(mission.Targets, target)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("%s: rows error: %w", op, err)
	}

	return &mission, nil
}
