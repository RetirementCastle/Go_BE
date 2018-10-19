// main_test.go

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize("root", "mi nombre es 123", "NursingHomes")

	ensureTableExists()

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func ensureTableExists() {
	if _, err := a.DB.Exec(nursinghomeCreationQuery); err != nil {
		log.Fatal(err)
	}
	if _, err := a.DB.Exec(branchesCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM nursinghome")
	a.DB.Exec("ALTER TABLE nursinghome AUTO_INCREMENT = 1")
	a.DB.Exec("DELETE FROM branches")
	a.DB.Exec("ALTER TABLE branches AUTO_INCREMENT = 1")
}

const nursinghomeCreationQuery = `
CREATE TABLE IF NOT EXISTS nursinghome
(
	idnursinghome INT NOT NULL AUTO_INCREMENT,
  	name VARCHAR(45) NOT NULL,
	PRIMARY KEY (idnursinghome)
)`

const branchesCreationQuery = `
CREATE TABLE IF NOT EXISTS branches (
	idbranches INT NOT NULL AUTO_INCREMENT,
	address VARCHAR(45) NOT NULL,
	total_rooms INT NOT NULL,
	available_rooms INT NOT NULL,
	nursinghome_idnursinghome INT NOT NULL,
	PRIMARY KEY (idbranches, nursinghome_idnursinghome),
	INDEX fk_branches_nursinghome_idx (nursinghome_idnursinghome ASC),
	CONSTRAINT fk_branches_nursinghome
	  FOREIGN KEY (nursinghome_idnursinghome)
	  REFERENCES nursinghome (idnursinghome)
	  ON DELETE NO ACTION
	  ON UPDATE NO ACTION
)`

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestEmptyNursinghome(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/nursinghomes", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	if body := response.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentNursinghome(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/nursinghome/45", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Nursinghome not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Nursinghome not found'. Got '%s'", m["error"])
	}
}

func TestCreateNursinghome(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test nh"}`)

	req, _ := http.NewRequest("POST", "/nursinghome", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test nh" {
		t.Errorf("Expected nh name to be 'test nh'. Got '%v'", m["name"])
	}
	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["idnursinghome"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["idnursinghome"])
	}
}

func addNursinghomes(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO nursinghome(name) VALUES('%s')", ("nh " + strconv.Itoa(i+1)))
		a.DB.Exec(statement)
	}
}

func addBranches(count int, nhid int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		statement := fmt.Sprintf("INSERT INTO branches(addres, total_rooms, available_rooms, nursinghome_idnursinghome) VALUES('%s', %d, %d, %d)", ("branch " + strconv.Itoa(i+1)), (i+1)*10, (i+1)*5, nhid)
		a.DB.Exec(statement)
	}
}

func TestGetNursinghome(t *testing.T) {
	clearTable()
	addNursinghomes(1)

	req, _ := http.NewRequest("GET", "/nursinghome/1", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetBranches(t *testing.T) {
	clearTable()
	addNursinghomes(1)
	addBranches(3, 1)

	req, _ := http.NewRequest("GET", "/branches", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetNHBranches(t *testing.T) {
	clearTable()
	addNursinghomes(1)
	addBranches(3, 1)

	req, _ := http.NewRequest("GET", "/nursinghome/1/branches", nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateNursinghome(t *testing.T) {
	clearTable()
	addNursinghomes(1)

	req, _ := http.NewRequest("GET", "/nursinghome/1", nil)
	response := executeRequest(req)
	var originalNursinghome map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &originalNursinghome)

	payload := []byte(`{"name":"test nh - updated name"}`)

	req, _ = http.NewRequest("PUT", "/nursinghome/1", bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["idnursinghome"] != originalNursinghome["idnursinghome"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalNursinghome["idnursinghome"], m["idnursinghome"])
	}

	if m["name"] == originalNursinghome["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalNursinghome["name"], m["name"], m["name"])
	}
}

func TestDeleteNursinghome(t *testing.T) {
	clearTable()
	addNursinghomes(1)

	req, _ := http.NewRequest("GET", "/nursinghome/1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/nursinghome/1", nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/nursinghome/1", nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}
