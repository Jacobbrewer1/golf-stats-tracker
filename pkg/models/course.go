// Package models contains the database interaction model code
//
// GENERATED BY GOSCHEMA. DO NOT EDIT.
package models

import (
	"github.com/prometheus/client_golang/prometheus"
)

// Course represents a row from 'course'.
type Course struct {
	Id      int    `db:"id,autoinc,pk"`
	RoundId int    `db:"round_id"`
	Name    string `db:"name"`
}

// CourseColumns is the sorted column names for the type Course
var CourseColumns = []string{"Id", "Name", "RoundId"}

// Insert inserts the Course to the database.
func (m *Course) Insert(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_Course"))
	defer t.ObserveDuration()

	const sqlstr = "INSERT INTO course (" +
		"`round_id`, `name`" +
		") VALUES (" +
		"?, ?" +
		")"

	DBLog(sqlstr, m.RoundId, m.Name)
	res, err := db.Exec(sqlstr, m.RoundId, m.Name)
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

func InsertManyCourses(db DB, ms ...*Course) error {
	if len(ms) == 0 {
		return nil
	}

	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_many_Course"))
	defer t.ObserveDuration()

	var sqlstr = "INSERT INTO course (" +
		"`round_id`,`name`" +
		") VALUES"

	var args []interface{}
	for _, m := range ms {
		sqlstr += " (" +
			"?,?" +
			"),"
		args = append(args, m.RoundId, m.Name)
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
func (m *Course) IsPrimaryKeySet() bool {
	return IsKeySet(m.Id)
}

// Update updates the Course in the database.
func (m *Course) Update(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("update_Course"))
	defer t.ObserveDuration()

	const sqlstr = "UPDATE course " +
		"SET `round_id` = ?, `name` = ? " +
		"WHERE `id` = ?"

	DBLog(sqlstr, m.RoundId, m.Name, m.Id)
	res, err := db.Exec(sqlstr, m.RoundId, m.Name, m.Id)
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

// InsertWithUpdate inserts the Course to the database, and tries to update
// on unique constraint violations.
func (m *Course) InsertWithUpdate(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_update_Course"))
	defer t.ObserveDuration()

	const sqlstr = "INSERT INTO course (" +
		"`round_id`, `name`" +
		") VALUES (" +
		"?, ?" +
		") ON DUPLICATE KEY UPDATE " +
		"`round_id` = VALUES(`round_id`), `name` = VALUES(`name`)"

	DBLog(sqlstr, m.RoundId, m.Name)
	res, err := db.Exec(sqlstr, m.RoundId, m.Name)
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

// Save saves the Course to the database.
func (m *Course) Save(db DB) error {
	if m.IsPrimaryKeySet() {
		return m.Update(db)
	}
	return m.Insert(db)
}

// SaveOrUpdate saves the Course to the database, but tries to update
// on unique constraint violations.
func (m *Course) SaveOrUpdate(db DB) error {
	if m.IsPrimaryKeySet() {
		return m.Update(db)
	}
	return m.InsertWithUpdate(db)
}

// Delete deletes the Course from the database.
func (m *Course) Delete(db DB) error {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("delete_Course"))
	defer t.ObserveDuration()

	const sqlstr = "DELETE FROM course WHERE `id` = ?"

	DBLog(sqlstr, m.Id)
	_, err := db.Exec(sqlstr, m.Id)

	return err
}

// CourseById retrieves a row from 'course' as a Course.
//
// Generated from primary key.
func CourseById(db DB, id int) (*Course, error) {
	t := prometheus.NewTimer(DatabaseLatency.WithLabelValues("insert_Course"))
	defer t.ObserveDuration()

	const sqlstr = "SELECT `id`, `round_id`, `name` " +
		"FROM course " +
		"WHERE `id` = ?"

	DBLog(sqlstr, id)
	var m Course
	if err := db.Get(&m, sqlstr, id); err != nil {
		return nil, err
	}

	return &m, nil
}
