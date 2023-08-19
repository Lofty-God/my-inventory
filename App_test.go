package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialise(DbUser, DbPassword, "Assessment")
	if err != nil {
		log.Fatal("error occured while initializing the database")
	}
	createTable()

	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS Products(
		id int NOT NULL AUTO_INCREMENT,
		name varchar(255) NOT NULL,
		Quantity int,
		price float(10,7),
		PRIMARY KEY(id)
    );`
	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

}

func clearTable() {
	a.DB.Exec("DELETE from Products")
	a.DB.Exec("alter table Products AUTO_INCREMENT=1")
	log.Println("clearTable")
}
func addProduct(name string, Quantity int, price float64) {
	query := fmt.Sprintf("INSERT into Products(name, Quantity, price) VALUES('%v', %v, %v)", name, Quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Println(err)
	}

}
func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 500)
	request, _ := http.NewRequest("GET", "/product/1", nil)
	Response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, Response.Code)

}

func checkStatusCode(t *testing.T, expectedStatusCode int, actualStatuscode int) {
	if expectedStatusCode != actualStatuscode {
		t.Errorf("expected Status: %v, received: %v", expectedStatusCode, actualStatuscode)
	}
}
func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder

}
func TestCreateProduct(t *testing.T) {
	var product = []byte(`{"name":"chair", "Quantity":20, "price":150}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("content_type", "Application/json")
	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)
	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["name"] != "chair" {
		t.Errorf("expected value:%v, Got:%v", "chair", m["name"])
	}
	if m["quantity"] != 20.0 {
		t.Errorf("expected value:%v, Got:%v", 20.0, m["quantity"])
	}
}
func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("Bicycle", 20, 50)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)

}
func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("microwave", 45, 234)
	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)

	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	var product = []byte(`{"microwave", 450, 234}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("content_type", "Application/json")
	response = sendRequest(req)

	var NewValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &NewValue)

	if oldValue["id"] != NewValue["id"] {
		t.Errorf("expected id: %v, Got id: %v", NewValue["id"], oldValue["id"])
	}

	if oldValue["name"] != NewValue["name"] {
		t.Errorf("expected name: %v, Got name: %v", NewValue["name"], oldValue["name"])
	}

	if oldValue["Quantity"] != NewValue["Quantity"] {
		t.Errorf("expected Quantity: %v, Got Quantity: %v", NewValue["Quantity"], oldValue["Quantity"])
	}

	if oldValue["price"] != NewValue["price"] {
		t.Errorf("expected price: %v, Got price: %v", NewValue["price"], oldValue["price"])
	}

}
