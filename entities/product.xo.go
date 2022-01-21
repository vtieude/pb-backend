package entities

// Code generated by xo. DO NOT EDIT.

import (
	"context"

	"github.com/elgris/sqrl"
)

// Product represents a row from 'product'.
type Product struct {
	ID         int     `json:"ID" db:"id"`                  // id
	Name       string  `json:"Name" db:"name"`              // name
	ProductKey string  `json:"ProductKey" db:"product_key"` // product_key
	Active     bool    `json:"Active" db:"active"`          // active
	Price      float64 `json:"Price" db:"price"`            // price
	// xo fields
	_exists, _deleted bool
}

type FilterProduct struct {
	ID         *int     // id
	Name       *string  // name
	ProductKey *string  // product_key
	Active     *bool    // active
	Price      *float64 // price

}

// Apply filter to sqrl Product .
func (p *Product) ApplyFilterSale(sqrlBuilder *sqrl.SelectBuilder, filter FilterProduct) bool {
	if filter.ID != nil {
		sqrlBuilder.Where(sqrl.Eq{"id": filter.ID})
	}
	if filter.Name != nil {
		sqrlBuilder.Where(sqrl.Eq{"name": filter.Name})
	}
	if filter.ProductKey != nil {
		sqrlBuilder.Where(sqrl.Eq{"product_key": filter.ProductKey})
	}
	if filter.Active != nil {
		sqrlBuilder.Where(sqrl.Eq{"active": filter.Active})
	}
	if filter.Price != nil {
		sqrlBuilder.Where(sqrl.Eq{"price": filter.Price})
	}

	return true
}

// Exists returns true when the Product exists in the database.
func (p *Product) Exists() bool {
	return p._exists
}

// Deleted returns true when the Product has been marked for deletion from
// the database.
func (p *Product) Deleted() bool {
	return p._deleted
}

// Insert inserts the Product to the database.
func (p *Product) Insert(ctx context.Context, db DB) error {
	switch {
	case p._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case p._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO product (` +
		`name, product_key, active, price` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`
	// run
	logf(sqlstr, p.Name, p.ProductKey, p.Active, p.Price)
	res, err := db.ExecContext(ctx, sqlstr, p.Name, p.ProductKey, p.Active, p.Price)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	p.ID = int(id)
	// set exists
	p._exists = true
	return nil
}

// Update updates a Product in the database.
func (p *Product) Update(ctx context.Context, db DB) error {
	switch {
	case !p._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case p._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE product SET ` +
		`name = ?, product_key = ?, active = ?, price = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, p.Name, p.ProductKey, p.Active, p.Price, p.ID)
	if _, err := db.ExecContext(ctx, sqlstr, p.Name, p.ProductKey, p.Active, p.Price, p.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the Product to the database.
func (p *Product) Save(ctx context.Context, db DB) error {
	if p.Exists() {
		return p.Update(ctx, db)
	}
	return p.Insert(ctx, db)
}

// Upsert performs an upsert for Product.
func (p *Product) Upsert(ctx context.Context, db DB) error {
	switch {
	case p._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO product (` +
		`id, name, product_key, active, price` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`name = VALUES(name), product_key = VALUES(product_key), active = VALUES(active), price = VALUES(price)`
	// run
	logf(sqlstr, p.ID, p.Name, p.ProductKey, p.Active, p.Price)
	if _, err := db.ExecContext(ctx, sqlstr, p.ID, p.Name, p.ProductKey, p.Active, p.Price); err != nil {
		return logerror(err)
	}
	// set exists
	p._exists = true
	return nil
}

// Delete deletes the Product from the database.
func (p *Product) Delete(ctx context.Context, db DB) error {
	switch {
	case !p._exists: // doesn't exist
		return nil
	case p._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM product ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, p.ID)
	if _, err := db.ExecContext(ctx, sqlstr, p.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	p._deleted = true
	return nil
}

// ProductByID retrieves a row from 'product' as a Product.
//
// Generated from index 'product_id_pkey'.
func ProductByID(ctx context.Context, db DB, id int) (*Product, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name, product_key, active, price ` +
		`FROM product ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	p := Product{
		_exists: true,
	}
	qb := sqrl.Expr(sqlstr, id)
	if err := db.QueryRowContext(ctx, &p, qb); err != nil {
		return nil, logerror(err)
	}
	return &p, nil
}

// ProductByProductKey retrieves a row from 'product' as a Product.
//
// Generated from index 'product_key'.
func ProductByProductKey(ctx context.Context, db DB, productKey string) (*Product, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, name, product_key, active, price ` +
		`FROM product ` +
		`WHERE product_key = ?`
	// run
	logf(sqlstr, productKey)
	p := Product{
		_exists: true,
	}
	qb := sqrl.Expr(sqlstr, productKey)
	if err := db.QueryRowContext(ctx, &p, qb); err != nil {
		return nil, logerror(err)
	}
	return &p, nil
}
