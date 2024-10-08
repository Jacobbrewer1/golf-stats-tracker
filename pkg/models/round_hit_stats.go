// Package models contains the database interaction model code
//
// GENERATED BY GOSCHEMA. DO NOT EDIT.
package models

import (
	usql "github.com/Jacobbrewer1/golf-stats-tracker/pkg/utils/sql"
	"github.com/prometheus/client_golang/prometheus"
)

// RoundHitStats represents a row from 'round_hit_stats'.
type RoundHitStats struct {
	Id           int       `db:"id,autoinc,pk"`
	RoundStatsId int       `db:"round_stats_id"`
	Type         usql.Enum `db:"type"`
	Miss         string    `db:"miss"`
	Count        int       `db:"count"`
}

// RoundHitStatsColumns is the sorted column names for the type RoundHitStats
var RoundHitStatsColumns = []string{"Count", "Id", "Miss", "RoundStatsId", "Type"}

// Insert inserts the RoundHitStats to the database.
func (m *RoundHitStats) Insert(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_RoundHitStats"))
	defer t.ObserveDuration()

	const sqlstr = "INSERT INTO round_hit_stats (" +
		"`round_stats_id`, `type`, `miss`, `count`" +
		") VALUES (" +
		"?, ?, ?, ?" +
		")"

	DBLog(sqlstr, m.RoundStatsId, m.Type, m.Miss, m.Count)
	res, err := db.Exec(sqlstr, m.RoundStatsId, m.Type, m.Miss, m.Count)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.Id = int(id)
	return nil
}

func InsertManyRoundHitStatss(db DB, ms ...*RoundHitStats) error {
	if len(ms) == 0 {
		return nil
	}

	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_many_RoundHitStats"))
	defer t.ObserveDuration()

	var sqlstr = "INSERT INTO round_hit_stats (" +
		"`round_stats_id`,`type`,`miss`,`count`" +
		") VALUES"

	var args []interface{}
	for _, m := range ms {
		sqlstr += " (" +
			"?,?,?,?" +
			"),"
		args = append(args, m.RoundStatsId, m.Type, m.Miss, m.Count)
	}

	DBLog(sqlstr, args...)
	res, err := db.Exec(sqlstr, args...)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	for i, m := range ms {
		m.Id = int(id + int64(i))
	}

	return nil
}

// IsPrimaryKeySet returns true if all primary key fields are set to none zero values
func (m *RoundHitStats) IsPrimaryKeySet() bool {
	return IsKeySet(m.Id)
}

// Update updates the RoundHitStats in the database.
func (m *RoundHitStats) Update(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("update_RoundHitStats"))
	defer t.ObserveDuration()

	const sqlstr = "UPDATE round_hit_stats " +
		"SET `round_stats_id` = ?, `type` = ?, `miss` = ?, `count` = ? " +
		"WHERE `id` = ?"

	DBLog(sqlstr, m.RoundStatsId, m.Type, m.Miss, m.Count, m.Id)
	res, err := db.Exec(sqlstr, m.RoundStatsId, m.Type, m.Miss, m.Count, m.Id)
	if err != nil {
		return err
	}

	// Requires clientFoundRows=true
	if i, err := res.RowsAffected(); err != nil {
		return err
	} else if i <= 0 {
		return ErrNoAffectedRows
	}

	return nil
}

// InsertWithUpdate inserts the RoundHitStats to the database, and tries to update
// on unique constraint violations.
func (m *RoundHitStats) InsertWithUpdate(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_update_RoundHitStats"))
	defer t.ObserveDuration()

	const sqlstr = "INSERT INTO round_hit_stats (" +
		"`round_stats_id`, `type`, `miss`, `count`" +
		") VALUES (" +
		"?, ?, ?, ?" +
		") ON DUPLICATE KEY UPDATE " +
		"`round_stats_id` = VALUES(`round_stats_id`), `type` = VALUES(`type`), `miss` = VALUES(`miss`), `count` = VALUES(`count`)"

	DBLog(sqlstr, m.RoundStatsId, m.Type, m.Miss, m.Count)
	res, err := db.Exec(sqlstr, m.RoundStatsId, m.Type, m.Miss, m.Count)
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	m.Id = int(id)
	return nil
}

// Save saves the RoundHitStats to the database.
func (m *RoundHitStats) Save(db DB) error {
	if m.IsPrimaryKeySet() {
		return m.Update(db)
	}
	return m.Insert(db)
}

// SaveOrUpdate saves the RoundHitStats to the database, but tries to update
// on unique constraint violations.
func (m *RoundHitStats) SaveOrUpdate(db DB) error {
	if m.IsPrimaryKeySet() {
		return m.Update(db)
	}
	return m.InsertWithUpdate(db)
}

// Delete deletes the RoundHitStats from the database.
func (m *RoundHitStats) Delete(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("delete_RoundHitStats"))
	defer t.ObserveDuration()

	const sqlstr = "DELETE FROM round_hit_stats WHERE `id` = ?"

	DBLog(sqlstr, m.Id)
	_, err := db.Exec(sqlstr, m.Id)

	return err
}

// RoundHitStatsById retrieves a row from 'round_hit_stats' as a RoundHitStats.
//
// Generated from primary key.
func RoundHitStatsById(db DB, id int) (*RoundHitStats, error) {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_RoundHitStats"))
	defer t.ObserveDuration()

	const sqlstr = "SELECT `id`, `round_stats_id`, `type`, `miss`, `count` " +
		"FROM round_hit_stats " +
		"WHERE `id` = ?"

	DBLog(sqlstr, id)
	var m RoundHitStats
	if err := db.Get(&m, sqlstr, id); err != nil {
		return nil, err
	}

	return &m, nil
}

// GetRoundStats Gets an instance of RoundStats
//
// Generated from constraint round_hit_stats_round_stats_id_fk
func (m *RoundHitStats) GetRoundStats(db DB) (*RoundStats, error) {
	return RoundStatsById(db, m.RoundStatsId)
}

// Valid values for the 'Type' enum column
var (
	RoundHitStatsTypeGREEN   = "GREEN"
	RoundHitStatsTypeFAIRWAY = "FAIRWAY"
)
