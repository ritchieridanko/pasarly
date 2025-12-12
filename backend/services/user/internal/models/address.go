package models

import "time"

type Address struct {
	ID           int64
	Recipient    string
	Phone        string
	Label        string
	Notes        *string
	IsPrimary    bool
	Country      string
	Subdivision1 *string
	Subdivision2 *string
	Subdivision3 *string
	Subdivision4 *string
	Street       string
	Postcode     string
	Latitude     float64
	Longitude    float64
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type CreateAddress struct {
	AuthID       int64
	Recipient    string
	Phone        string
	Label        string
	Notes        *string
	IsPrimary    bool
	Country      string
	Subdivision1 *string
	Subdivision2 *string
	Subdivision3 *string
	Subdivision4 *string
	Street       string
	Postcode     string
	Latitude     float64
	Longitude    float64
}

type UpdateAddress struct {
	AuthID       int64
	AddressID    int64
	Recipient    *string
	Phone        *string
	Label        *string
	Notes        *string
	Country      *string
	Subdivision1 *string
	Subdivision2 *string
	Subdivision3 *string
	Subdivision4 *string
	Street       *string
	Postcode     *string
	Latitude     *float64
	Longitude    *float64
}

type SetPrimaryAddress struct {
	AuthID    int64
	AddressID int64
}
