package entities

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql/driver"
	"fmt"
)

// SaleStatus is the 'sale_status' enum type from schema 'app_db'.
type SaleStatus uint16

// SaleStatus values.
const (
	// SaleStatusPending is the 'pending' sale_status.
	SaleStatusPending SaleStatus = 1
	// SaleStatusSaled is the 'saled' sale_status.
	SaleStatusSaled SaleStatus = 2
	// SaleStatusFailed is the 'failed' sale_status.
	SaleStatusFailed SaleStatus = 3
)

// String satisfies the fmt.Stringer interface.
func (ss SaleStatus) String() string {
	switch ss {
	case SaleStatusPending:
		return "pending"
	case SaleStatusSaled:
		return "saled"
	case SaleStatusFailed:
		return "failed"
	}
	return fmt.Sprintf("SaleStatus(%d)", ss)
}

// MarshalText marshals SaleStatus into text.
func (ss SaleStatus) MarshalText() ([]byte, error) {
	return []byte(ss.String()), nil
}

// UnmarshalText unmarshals SaleStatus from text.
func (ss *SaleStatus) UnmarshalText(buf []byte) error {
	switch str := string(buf); str {
	case "pending":
		*ss = SaleStatusPending
	case "saled":
		*ss = SaleStatusSaled
	case "failed":
		*ss = SaleStatusFailed
	default:
		return ErrInvalidSaleStatus(str)
	}
	return nil
}

// Value satisfies the driver.Valuer interface.
func (ss SaleStatus) Value() (driver.Value, error) {
	return ss.String(), nil
}

// Scan satisfies the sql.Scanner interface.
func (ss *SaleStatus) Scan(v interface{}) error {
	if buf, ok := v.([]byte); ok {
		return ss.UnmarshalText(buf)
	}
	return ErrInvalidSaleStatus(fmt.Sprintf("%T", v))
}

// NullSaleStatus represents a null 'sale_status' enum for schema 'app_db'.
type NullSaleStatus struct {
	SaleStatus SaleStatus
	// Valid is true if SaleStatus is not null.
	Valid bool
}

// Value satisfies the driver.Valuer interface.
func (nss NullSaleStatus) Value() (driver.Value, error) {
	if !nss.Valid {
		return nil, nil
	}
	return nss.SaleStatus.Value()
}

// Scan satisfies the sql.Scanner interface.
func (nss *NullSaleStatus) Scan(v interface{}) error {
	if v == nil {
		nss.SaleStatus, nss.Valid = 0, false
		return nil
	}
	err := nss.SaleStatus.Scan(v)
	nss.Valid = err == nil
	return err
}

// ErrInvalidSaleStatus is the invalid SaleStatus error.
type ErrInvalidSaleStatus string

// Error satisfies the error interface.
func (err ErrInvalidSaleStatus) Error() string {
	return fmt.Sprintf("invalid SaleStatus(%s)", string(err))
}
