package domain

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGeoIP2GetCountryCodeForIP(t *testing.T) {
	geoip2, err := NewGeoIP2Domain("../geoip2/GeoLite2-Country.mmdb")
	assert.NoError(t, err, "Can't create the GeoIP2Domain logic")
	assert.NotNil(t, geoip2)
	if geoip2 != nil {
		germanCode, err := geoip2.GetCountryCodeForIP(net.ParseIP("46.243.122.48"))
		assert.NoError(t, err, "Can't get germany country code")
		assert.Equal(t, "DE", germanCode)

		usCode, err := geoip2.GetCountryCodeForIP(net.ParseIP("54.208.114.32"))
		assert.NoError(t, err, "Can't get US country code")
		assert.Equal(t, "US", usCode)

		localCode, err := geoip2.GetCountryCodeForIP(net.ParseIP("127.0.0.1"))
		assert.NoError(t, err, "Can't get local code")
		assert.Equal(t, "", localCode)

		localCode, err = geoip2.GetCountryCodeForIP(net.ParseIP("10.0.0.1"))
		assert.NoError(t, err, "Can't get local code")
		assert.Equal(t, "", localCode)
	}	
}
