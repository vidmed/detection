package detection

import (
	"encoding/json"
	"strings"

	"github.com/pkg/errors"
	"github.com/simplereach/51degrees"
)

// This struct MUST match the defaultProperties var below
type UserAgentData struct {
	BrowserName     string `json:"browserName"`
	BrowserVersion  string `json:"browserVersion"`
	PlatformName    string `json:"platformName"`
	PlatformVendor  string `json:"platformVendor"`
	PlatformVersion string `json:"platformVersion"`
	IsMobile        string `json:"isMobile"`
	DeviceType      string `json:"deviceType"`
}

// UAParser is
type UAParser interface {
	Parse(userAgent string) (*UserAgentData, error)
}

type FiftyOneDegreesParser struct {
	provider *fiftyonedegrees.FiftyoneDegreesProvider
}

func (p *FiftyOneDegreesParser) Parse(userAgent string) (*UserAgentData, error) {
	uaJSON := p.provider.Parse(userAgent)
	if uaJSON == "" {
		return nil, errors.New("parser empty response")
	}

	uaData := &UserAgentData{}
	dec := json.NewDecoder(strings.NewReader(uaJSON))
	if err := dec.Decode(uaData); err != nil {
		return nil, errors.Wrap(err, "parser json decode errors")
	}
	return uaData, nil
}

func (p *FiftyOneDegreesParser) Close() {
	p.provider.Close()
}

func NewFiftyoneDegreesParser(DBPath string) (*FiftyOneDegreesParser, error) {
	provider, err := fiftyonedegrees.NewFiftyoneDegreesProvider(DBPath,
		`BrowserName, BrowserVersion, PlatformName, PlatformVendor, PlatformVersion, IsMobile, DeviceType`, 0, 0)
	if err != nil {
		return nil, errors.Wrap(err, "error creating FiftyoneDegreesProvider")
	}

	return &FiftyOneDegreesParser{provider: provider}, nil
}
