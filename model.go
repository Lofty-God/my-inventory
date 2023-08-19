package main

import (
	"database/sql"
	"errors"
	"fmt"
)

type product struct {
	Id       int     `json:"id"`
	Name     string  `json:"name"`
	Quantity int     `json:"quantity"`
	Price    float64 `json:"price"`
}

func getProducts(db *sql.DB) ([]product, error) {
	query := "select id, name,Quantity,price from Products"
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}

	Products := []product{}
	for rows.Next() {
		var p product
		err := rows.Scan(&p.Id, &p.Name, &p.Quantity, &p.Price)
		if err != nil {
			return nil, err
		}
		Products = append(Products, p)
	}

	return Products, nil
}
func (p *product) getProduct(db *sql.DB) error {
	query := fmt.Sprintf("SELECT name, Quantity,price from Products where id=%v", p.Id)
	row := db.QueryRow(query)
	err := row.Scan(&p.Name, &p.Quantity, &p.Price)
	if err != nil {
		return err
	}
	return nil

}
func (p *product) createProduct(db *sql.DB) error {
	query := fmt.Sprintf("INSERT into Products(name, Quantity, price ) values('%v', %v,  %v)", p.Name, p.Quantity, p.Price)
	result, err := db.Exec(query)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	p.Id = int(id)
	return nil

}
func (p *product) updateproduct(db *sql.DB) error {
	query := fmt.Sprintf("update Products set name='%v', Quantity=%v, price=%v where id=%v", p.Name, p.Quantity, p.Price, p.Id)
	result, err := db.Exec(query)
	rowsAffected, err := result.RowsAffected()
	if rowsAffected == 0 {
		return errors.New("no such rows exit")
	}
	return err
}
func (p *product) deleteProduct(db *sql.DB) error {
	query := fmt.Sprintf("delete from Products where id=%v", p.Id)
	_, err := db.Exec(query)
	return err

}
