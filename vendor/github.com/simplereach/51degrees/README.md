#51Degrees for golang
51Degrees UserAgent.  Note that the primary difference between https://github.com/hdczsf/51degrees.go and https://github.com/simplereach/51degrees is that the SimpleReach version supports the 3.2 data file format.

### Installation
`go get github.com/simplereach/51degrees`

### Example
```go
package main

import (
	"fmt"
	"github.com/simplereach/51degrees"
)

func main() {
	properties := "DeviceType, IsMobile, IsSmartPhone, IsTablet, IsTv, HardwareName, HardwareVendor, HardwareModel, BrowserName, BrowserVersion, PlatformName, PlatformVersion, ScreenPixelsWidth, ScreenPixelsHeight"
	testUA := "Mozilla/5.0 (Windows NT 6.1; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/43.0.2357.124 Safari/537.36"
	item, err := NewFiftyoneDegrees("51Degrees-EnterpriseV3_2.dat", properties)
	if err != nil {
		fmt.Println("err=", err)
		return
	}
	fmt.Println(item.Parse(testUA))
	item.Close()
}

```
