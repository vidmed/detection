package detection

import (
	"net"

	"github.com/oschwald/geoip2-golang"
	"github.com/pkg/errors"
)

const LanguageEN = "en"

type CountryData struct {
	ISOCode string `json:"isoCode"`
	Name    string `json:"name"`
}

type CountryDetector interface {
	Detect(ip net.IP) (*CountryData, error)
}

type MaxMindCountryDetector struct {
	countryDB *geoip2.Reader
}

func (d *MaxMindCountryDetector) Detect(ip net.IP) (*CountryData, error) {
	record, err := d.countryDB.Country(ip)
	if err != nil {
		return nil, errors.New("error getting country from MaxMind")
	}
	return &CountryData{
		ISOCode: record.Country.IsoCode,
		Name:    record.Country.Names[LanguageEN],
	}, nil
}

func NewMaxMindCountryDetector(DBPath string) (*MaxMindCountryDetector, error) {
	countryDB, err := geoip2.Open(DBPath)
	if err != nil {
		return nil, errors.Wrap(err, "error opening MaxMindCountry database")
	}

	return &MaxMindCountryDetector{countryDB: countryDB}, nil
}
