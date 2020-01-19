package domain

import (
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const geoip2Path = "../geoip2/GeoLite2-Country.mmdb"

func setupGeoip2Domain(t *testing.T) *GeoIP2Domain {
	geoip2, err := NewGeoIP2Domain(geoip2Path)
	require.NoError(t, err, "Can't create the GeoIP2Domain logic")
	return geoip2
}

func TestGeoIP2GetCountryCodeForIP(t *testing.T) {
	geoip2Domain := setupGeoip2Domain(t)
	assert.NotNil(t, geoip2Domain)
	if geoip2Domain != nil {
		germanCode, err := geoip2Domain.GetCountryCodeForIP(net.ParseIP("46.243.122.48"))
		assert.NoError(t, err, "Can't get germany country code")
		assert.Equal(t, "de", germanCode)

		usCode, err := geoip2Domain.GetCountryCodeForIP(net.ParseIP("54.208.114.32"))
		assert.NoError(t, err, "Can't get US country code")
		assert.Equal(t, "us", usCode)

		localCode, err := geoip2Domain.GetCountryCodeForIP(net.ParseIP("127.0.0.1"))
		assert.NoError(t, err, "Can't get local code")
		assert.Equal(t, "", localCode)

		localCode, err = geoip2Domain.GetCountryCodeForIP(net.ParseIP("192.168.178.2"))
		assert.NoError(t, err, "Can't get local code")
		assert.Equal(t, "", localCode)

		localCode, err = geoip2Domain.GetCountryCodeForIP(net.ParseIP("10.0.0.1"))
		assert.NoError(t, err, "Can't get local code")
		assert.Equal(t, "", localCode)
	}
}
