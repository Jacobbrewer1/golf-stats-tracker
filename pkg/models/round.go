// Package models contains the database interaction model code
//
// GENERATED BY GOSCHEMA. DO NOT EDIT.
package models

import (
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// Round represents a row from 'round'.
type Round struct {
	Id      int       `db:"id,autoinc,pk"`
	UserId  int       `db:"user_id"`
	TeeTime time.Time `db:"tee_time"`
}

// RoundColumns is the sorted column names for the type Round
var RoundColumns = []string{"Id", "TeeTime", "UserId"}

// Insert inserts the Round to the database.
func (m *Round) Insert(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_Round"))
	defer t.ObserveDuration()

	const sqlstr = "INSERT INTO round (" +
		"`user_id`, `tee_time`" +
		") VALUES (" +
		"?, ?" +
		")"

	DBLog(sqlstr, m.UserId, m.TeeTime)
	res, err := db.Exec(sqlstr, m.UserId, m.TeeTime)
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

func InsertManyRounds(db DB, ms ...*Round) error {
	if len(ms) == 0 {
		return nil
	}

	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_many_Round"))
	defer t.ObserveDuration()

	var sqlstr = "INSERT INTO round (" +
		"`user_id`,`tee_time`" +
		") VALUES"

	var args []interface{}
	for _, m := range ms {
		sqlstr += " (" +
			"?,?" +
			"),"
		args = append(args, m.UserId, m.TeeTime)
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
func (m *Round) IsPrimaryKeySet() bool {
	return IsKeySet(m.Id)
}

// Update updates the Round in the database.
func (m *Round) Update(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("update_Round"))
	defer t.ObserveDuration()

	const sqlstr = "UPDATE round " +
		"SET `user_id` = ?, `tee_time` = ? " +
		"WHERE `id` = ?"

	DBLog(sqlstr, m.UserId, m.TeeTime, m.Id)
	res, err := db.Exec(sqlstr, m.UserId, m.TeeTime, m.Id)
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

// InsertWithUpdate inserts the Round to the database, and tries to update
// on unique constraint violations.
func (m *Round) InsertWithUpdate(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_update_Round"))
	defer t.ObserveDuration()

	const sqlstr = "INSERT INTO round (" +
		"`user_id`, `tee_time`" +
		") VALUES (" +
		"?, ?" +
		") ON DUPLICATE KEY UPDATE " +
		"`user_id` = VALUES(`user_id`), `tee_time` = VALUES(`tee_time`)"

	DBLog(sqlstr, m.UserId, m.TeeTime)
	res, err := db.Exec(sqlstr, m.UserId, m.TeeTime)
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

// Save saves the Round to the database.
func (m *Round) Save(db DB) error {
	if m.IsPrimaryKeySet() {
		return m.Update(db)
	}
	return m.Insert(db)
}

// SaveOrUpdate saves the Round to the database, but tries to update
// on unique constraint violations.
func (m *Round) SaveOrUpdate(db DB) error {
	if m.IsPrimaryKeySet() {
		return m.Update(db)
	}
	return m.InsertWithUpdate(db)
}

// Delete deletes the Round from the database.
func (m *Round) Delete(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("delete_Round"))
	defer t.ObserveDuration()

	const sqlstr = "DELETE FROM round WHERE `id` = ?"

	DBLog(sqlstr, m.Id)
	_, err := db.Exec(sqlstr, m.Id)

	return err
}

// RoundById retrieves a row from 'round' as a Round.
//
// Generated from primary key.
func RoundById(db DB, id int) (*Round, error) {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_Round"))
	defer t.ObserveDuration()

	const sqlstr = "SELECT `id`, `user_id`, `tee_time` " +
		"FROM round " +
		"WHERE `id` = ?"

	DBLog(sqlstr, id)
	var m Round
	if err := db.Get(&m, sqlstr, id); err != nil {
		return nil, err
	}

	return &m, nil
}

// GetUser Gets an instance of User
//
// Generated from constraint round_user_id_fk
func (m *Round) GetUser(db DB) (*User, error) {
	return UserById(db, m.UserId)
}
