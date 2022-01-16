package entities

// Code generated by xo. DO NOT EDIT.

import (
	"context"

	"github.com/elgris/sqrl"
)

// UserRole represents a row from 'user_role'.
type UserRole struct {
	ID     int  `json:"ID" db:"id"`          // id
	FkUser int  `json:"FkUser" db:"fk_user"` // fk_user
	FkRole int  `json:"FkRole" db:"fk_role"` // fk_role
	Active bool `json:"Active" db:"active"`  // active
	// xo fields
	_exists, _deleted bool
}

type FilterUserRole struct {
	ID     *int  // id
	FkUser *int  // fk_user
	FkRole *int  // fk_role
	Active *bool // active

}

// Apply filter to sqrl UserRole .
func (ur *UserRole) ApplyFilterSale(sqrlBuilder *sqrl.SelectBuilder, filter FilterUserRole) bool {
	if filter.ID != nil {
		sqrlBuilder.Where(sqrl.Eq{"id": filter.ID})
	}
	if filter.FkUser != nil {
		sqrlBuilder.Where(sqrl.Eq{"fk_user": filter.FkUser})
	}
	if filter.FkRole != nil {
		sqrlBuilder.Where(sqrl.Eq{"fk_role": filter.FkRole})
	}
	if filter.Active != nil {
		sqrlBuilder.Where(sqrl.Eq{"active": filter.Active})
	}

	return true
}

// Exists returns true when the UserRole exists in the database.
func (ur *UserRole) Exists() bool {
	return ur._exists
}

// Deleted returns true when the UserRole has been marked for deletion from
// the database.
func (ur *UserRole) Deleted() bool {
	return ur._deleted
}

// Insert inserts the UserRole to the database.
func (ur *UserRole) Insert(ctx context.Context, db DB) error {
	switch {
	case ur._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case ur._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (primary key generated and returned by database)
	const sqlstr = `INSERT INTO user_role (` +
		`fk_user, fk_role, active` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`
	// run
	logf(sqlstr, ur.FkUser, ur.FkRole, ur.Active)
	res, err := db.ExecContext(ctx, sqlstr, ur.FkUser, ur.FkRole, ur.Active)
	if err != nil {
		return logerror(err)
	}
	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return logerror(err)
	} // set primary key
	ur.ID = int(id)
	// set exists
	ur._exists = true
	return nil
}

// Update updates a UserRole in the database.
func (ur *UserRole) Update(ctx context.Context, db DB) error {
	switch {
	case !ur._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case ur._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE user_role SET ` +
		`fk_user = ?, fk_role = ?, active = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, ur.FkUser, ur.FkRole, ur.Active, ur.ID)
	if _, err := db.ExecContext(ctx, sqlstr, ur.FkUser, ur.FkRole, ur.Active, ur.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the UserRole to the database.
func (ur *UserRole) Save(ctx context.Context, db DB) error {
	if ur.Exists() {
		return ur.Update(ctx, db)
	}
	return ur.Insert(ctx, db)
}

// Upsert performs an upsert for UserRole.
func (ur *UserRole) Upsert(ctx context.Context, db DB) error {
	switch {
	case ur._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO user_role (` +
		`id, fk_user, fk_role, active` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`fk_user = VALUES(fk_user), fk_role = VALUES(fk_role), active = VALUES(active)`
	// run
	logf(sqlstr, ur.ID, ur.FkUser, ur.FkRole, ur.Active)
	if _, err := db.ExecContext(ctx, sqlstr, ur.ID, ur.FkUser, ur.FkRole, ur.Active); err != nil {
		return logerror(err)
	}
	// set exists
	ur._exists = true
	return nil
}

// Delete deletes the UserRole from the database.
func (ur *UserRole) Delete(ctx context.Context, db DB) error {
	switch {
	case !ur._exists: // doesn't exist
		return nil
	case ur._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM user_role ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, ur.ID)
	if _, err := db.ExecContext(ctx, sqlstr, ur.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	ur._deleted = true
	return nil
}

// UserRoleByFkRole retrieves a row from 'user_role' as a UserRole.
//
// Generated from index 'fk_role'.
func UserRoleByFkRole(ctx context.Context, db DB, fkRole int) ([]*UserRole, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, fk_user, fk_role, active ` +
		`FROM user_role ` +
		`WHERE fk_role = ?`
	// run
	logf(sqlstr, fkRole)
	// process
	var res []*UserRole
	qb := sqrl.Expr(sqlstr, fkRole)
	if err := db.QueryContext(ctx, &res, qb); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// UserRoleByFkUser retrieves a row from 'user_role' as a UserRole.
//
// Generated from index 'fk_user'.
func UserRoleByFkUser(ctx context.Context, db DB, fkUser int) ([]*UserRole, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, fk_user, fk_role, active ` +
		`FROM user_role ` +
		`WHERE fk_user = ?`
	// run
	logf(sqlstr, fkUser)
	// process
	var res []*UserRole
	qb := sqrl.Expr(sqlstr, fkUser)
	if err := db.QueryContext(ctx, &res, qb); err != nil {
		return nil, logerror(err)
	}
	return res, nil
}

// UserRoleByID retrieves a row from 'user_role' as a UserRole.
//
// Generated from index 'user_role_id_pkey'.
func UserRoleByID(ctx context.Context, db DB, id int) (*UserRole, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, fk_user, fk_role, active ` +
		`FROM user_role ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	ur := UserRole{
		_exists: true,
	}
	qb := sqrl.Expr(sqlstr, id)
	if err := db.QueryRowContext(ctx, &ur, qb); err != nil {
		return nil, logerror(err)
	}
	return &ur, nil
}

// User returns the User associated with the UserRole's (FkUser).
//
// Generated from foreign key 'user_role_ibfk_1'.
func (ur *UserRole) User(ctx context.Context, db DB) (*User, error) {
	return UserByID(ctx, db, ur.FkUser)
}

// Role returns the Role associated with the UserRole's (FkRole).
//
// Generated from foreign key 'user_role_ibfk_2'.
func (ur *UserRole) Role(ctx context.Context, db DB) (*Role, error) {
	return RoleByID(ctx, db, ur.FkRole)
}
