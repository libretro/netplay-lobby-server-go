package domain

import (
	"fmt"
	"github.com/oschwald/maxminddb-golang"
	"net"
)

// GeoIP2Domain abstracts the GeoIP2 country database domain logic.
type GeoIP2Domain struct {
	db *maxminddb.Reader
}

type countryRecord struct {
	Country struct {
		ISOCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

// NewGeoIP2Domain creates a new domain object for the GeoIP2 country database. Need the path to a maxminddb file.
func NewGeoIP2Domain(path string) (*GeoIP2Domain, error) {
	db, err := maxminddb.Open(path)
	if err != nil {
		return nil, fmt.Errorf("can't open geoip2 country database %s: %w", path, err)
	}

	return &GeoIP2Domain{db}, nil
}

// GetCountryCodeForIP returns the two letter country code (ISO 3166-1) for the given IP.
func (d *GeoIP2Domain) GetCountryCodeForIP(ip net.IP) (string, error) {
	record := &countryRecord{}

	err := d.db.Lookup(ip, record)
	if err != nil {
		return "", fmt.Errorf("can't lookup country for IP %s: %w", ip, err)
	}

	return record.Country.ISOCode, nil
}

// Close needs to be called to properly close the internal maxminddb database.
func (d *GeoIP2Domain) Close() {
	d.db.Close()
}
