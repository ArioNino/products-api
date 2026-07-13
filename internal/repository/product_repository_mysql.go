package repository

import (
	"database/sql"
	"product-api/internal/model"
	"fmt"
)

type productRepositoryMySQL struct {
	db *sql.DB
}

func NewProductRepositoryMySQL(db *sql.DB) *productRepositoryMySQL {
	return &productRepositoryMySQL{db: db}
}

func (r *productRepositoryMySQL) GetAll() ([]model.Product, error) {
	rows, err := r.db.Query("SELECT id, name, price, stock FROM products")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(&p.ID, &p.Name, &p.Price, &p.Stock); err != nil {
			return nil, err
		}
		products = append(products, p)
	}
	return products, nil
}

func (r *productRepositoryMySQL) Create(p model.Product) (model.Product, error) {
	result, err := r.db.Exec("INSERT INTO products (name, price, stock) VALUES (?, ?, ?)", p.Name, p.Price, p.Stock)
	if err != nil {
		return model.Product{}, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return model.Product{}, err
	}

	p.ID = int(id)
	return p, nil
}

func (r *productRepositoryMySQL) GetByID(id int) (model.Product, error) {
	var p model.Product
	err := r.db.QueryRow(
		"SELECT id, name, price, stock FROM products WHERE id = ?", 
		id).Scan(&p.ID, &p.Name, &p.Price, &p.Stock)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
		}
		return model.Product{}, err
	}
	return p, nil
}

func (r *productRepositoryMySQL) Update(id int, updated model.Product) (model.Product, error) {
	_, err := r.db.Exec(
		"UPDATE products SET name = ?, price = ?, stock = ? WHERE id = ?", 
		updated.Name, updated.Price, updated.Stock, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return model.Product{}, fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
		}
		return model.Product{}, err
	}

	updated.ID = id
	return updated, nil
}

func (r *productRepositoryMySQL) Delete(id int) error {
	_, err := r.db.Exec("DELETE FROM products WHERE id = ?", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("produk dengan ID %d tidak ditemukan", id)
		}
		return err
	}
	return nil
}