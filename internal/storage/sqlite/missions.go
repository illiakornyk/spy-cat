package sqlite

import (
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
