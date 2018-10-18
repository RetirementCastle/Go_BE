// model.go

package main

import (
	"database/sql"
	"fmt"
)

type nursinghome struct {
	ID   int    `json:"idnursinghome"`
	Name string `json:"name"`
}

func (n *nursinghome) getNursinghome(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name FROM nursinghome WHERE idnursinghome=%d", n.ID)
	return db.QueryRow(statement).Scan(&n.Name)
}

func (n *nursinghome) updateNursinghome(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE nursinghome SET name='%s' WHERE idnursinghome=%d", n.Name, n.ID)
	_, err := db.Exec(statement)
	return err
}

func (n *nursinghome) deleteNursinghome(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM nursinghome WHERE idnursinghome=%d", n.ID)
	_, err := db.Exec(statement)
	return err
}

func (n *nursinghome) createNursinghome(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO nursinghome(name) VALUES('%s')", n.Name)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&n.ID)

	if err != nil {
		return err
	}

	return nil
}

func getNursinghomes(db *sql.DB, start, count int) ([]nursinghome, error) {
	statement := fmt.Sprintf("SELECT idnursinghome, name FROM nursinghome LIMIT %d OFFSET %d", count, start)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	nursinghomes := []nursinghome{}

	for rows.Next() {
		var n nursinghome
		if err := rows.Scan(&n.ID, &n.Name); err != nil {
			return nil, err
		}
		nursinghomes = append(nursinghomes, n)
	}

	return nursinghomes, nil
}
