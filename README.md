# Installation
1. Install golang
https://golang.org/doc/install
2. `go get github.com/vidmed/detection`
3. `cd $GOPATH/src/github.com/vidmed/detection/cmd`
4. `go build && ./cmd -config=config.toml` 
OR 
`go install && cmd -config=config.toml`

# Config file description
1. ListenAddr - address used by web server to listen to
2. ListenPort - port used by web server to listen to
3. LogLevel - logging level (panic = 0, fatal = 1, error = 2, warning = 3, info = 4, debug = 5)
4. MaxMinDBPath - path to country db file used by MaxMind to detect country using IP
5. FiftyOneDegreesDBPath - path to db file used by FiftyOneDegrees to parse user agent data

# Important information
Default databases for MaxMind and FiftyOneDegrees located in `data` folder. 
They are free and do not provide all necessary information. For production use please obtain commercial databases.