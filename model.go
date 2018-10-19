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

type branch struct {
	ID             int    `json:"idbranches"`
	Address        string `json:"address"`
	Totalrooms     int    `json:"total_rooms"`
	Availablerooms int    `json:"available_rooms"`
	NHID           int    `json:"nursinghome_idnursinghome"`
}

func (n *nursinghome) getNursinghome(db *sql.DB) error {
	statement := fmt.Sprintf("SELECT name FROM nursinghome WHERE idnursinghome=%d", n.ID)
	return db.QueryRow(statement).Scan(&n.Name)
}

func getBranches(db *sql.DB, start, count int) ([]branch, error) {
	statement := fmt.Sprintf("SELECT idbranches, address, total_rooms, available_rooms, nursinghome_idnursinghome FROM branches LIMIT %d OFFSET %d", count, start)
	//	statement := fmt.Sprintf("SELECT idbranches, addres, total_rooms, available_rooms, nursinghome_idnursinghome FROM branches")
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	branches := []branch{}

	for rows.Next() {
		var b branch
		if err := rows.Scan(&b.ID, &b.Address, &b.Totalrooms, &b.Availablerooms, &b.NHID); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}
	return branches, nil
}

func getNHBranches(db *sql.DB, idnh int) ([]branch, error) {
	statement := fmt.Sprintf("SELECT address, total_rooms, available_rooms FROM branches WHERE nursinghome_idnursinghome = %d;", idnh)
	rows, err := db.Query(statement)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	branches := []branch{}

	for rows.Next() {
		var b branch
		if err := rows.Scan(&b.Address, &b.Totalrooms, &b.Availablerooms); err != nil {
			return nil, err
		}
		branches = append(branches, b)
	}
	return branches, nil
}

func (n *nursinghome) updateNursinghome(db *sql.DB) error {
	statement := fmt.Sprintf("UPDATE nursinghome SET name='%s' WHERE idnursinghome=%d", n.Name, n.ID)
	_, err := db.Exec(statement)
	return err
}

func (b *branch) updateBranch(db *sql.DB, idb int) error {
	statement := fmt.Sprintf("UPDATE branches SET address='%s', total_rooms=%d, available_rooms=%d WHERE idbranches=%d", b.Address, b.Totalrooms, b.Availablerooms, idb)
	_, err := db.Exec(statement)

	return err
}

func (n *nursinghome) deleteNursinghome(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM nursinghome WHERE idnursinghome=%d", n.ID)
	_, err := db.Exec(statement)
	db.Exec("DELETE FROM branches WHERE nursinghome_idnnursinghome=%d", n.ID)
	return err
}

func (b *branch) deleteBranch(db *sql.DB) error {
	statement := fmt.Sprintf("DELETE FROM branches WHERE idbranches=%d", b.ID)
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

func (b *branch) createBranch(db *sql.DB) error {
	statement := fmt.Sprintf("INSERT INTO branches(address, total_rooms, available_rooms, nursinghome_idnursinghome) VALUES('%s', %d, %d, %d)", b.Address, b.Totalrooms, b.Availablerooms, b.NHID)
	_, err := db.Exec(statement)

	if err != nil {
		return err
	}

	err = db.QueryRow("SELECT LAST_INSERT_ID()").Scan(&b.ID)

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
