package detection

import (
	"net"
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

type Response struct {
	UA      *UserAgentData `json:"ua"`
	Country *CountryData   `json:"country"`
	Domain  string         `json:"urlDomain"`
}

type BidRequestParser struct {
	countryDetector CountryDetector
	uaParser        UAParser
}

func NewBidRequestParser(countryDetector CountryDetector, uaParser UAParser) *BidRequestParser {
	return &BidRequestParser{
		countryDetector: countryDetector,
		uaParser:        uaParser,
	}
}

func (p *BidRequestParser) Parse(req *BidRequest) (resp *Response, err error) {
	resp = new(Response)
	if req.Device.UA == "" {
		return nil, errors.New("empty User Agent in Bid Request")
	}
	resp.UA, err = p.uaParser.Parse(req.Device.UA)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse User Agent")
	}

	if req.Device.IP == "" {
		return nil, errors.New("empty IP in Bid Request")
	}
	ip := net.ParseIP(req.Device.IP)
	if ip == nil {
		return nil, errors.Wrap(err, "failed to parse IP")
	}
	resp.Country, err = p.countryDetector.Detect(ip)
	if ip == nil {
		return nil, errors.Wrap(err, "failed to detect country from IP")
	}

	if req.Site.Page != "" {
		resp.Domain, err = parseDomain(req.Site.Page)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse domain from site URL")
		}
	}

	return
}

func parseDomain(u string) (string, error) {
	parsed, err := url.Parse(u)
	if err != nil {
		return "", err
	}
	parts := strings.Split(parsed.Hostname(), ".")
	if len(parts) < 2 {
		return "", errors.New("invalid hostname")
	}
	return parts[len(parts)-2] + "." + parts[len(parts)-1], nil
}
