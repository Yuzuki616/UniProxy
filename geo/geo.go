package geo

import _ "embed"

//go:embed geoip.dat
var Ip []byte

//go:embed geosite.dat
var Site []byte
