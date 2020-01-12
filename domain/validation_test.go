package domain

import (
	"testing"
	"net"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var testCoreWhitelist = []string{
	"snes",
	"bsnes",
}

var testStringBlacklist = []string{
	".*badWord.*",
	"^prefixTest.*$",
	"\\s{3,}",
}

var testIPBlacklist = []string{
	"127.0.0.1",
	"2001:db8:0:8d3:0:8a2e:70:7344",
}

func TestValidationDomainValidCreation(t *testing.T) {
	_, err := NewValidationDomain(testCoreWhitelist, testStringBlacklist, testIPBlacklist)
	require.NoError(t, err)
}

func TestValidationDomainInvalidIP(t *testing.T) {
	_, err := NewValidationDomain(testCoreWhitelist, testStringBlacklist, []string{"256.123.12.3"})
	require.Error(t, err)

	_, err = NewValidationDomain(testCoreWhitelist, testStringBlacklist, []string{"2001:db8:0:8d3:0:8a2ef:70:7344"})
	require.Error(t, err)
}

func TestValidationDomainRegexpShouldNotCompile(t *testing.T) {
	_, err := NewValidationDomain(testCoreWhitelist, []string{"["}, testIPBlacklist)
	require.Error(t, err)

	_, err = NewValidationDomain(testCoreWhitelist, []string{"[0-9]++"}, testIPBlacklist)
	require.Error(t, err)
}

func TestValidationDomainValidateCore(t *testing.T) {
	validationDomain, err := NewValidationDomain(testCoreWhitelist, testStringBlacklist, testIPBlacklist)
	require.NoError(t, err)

	assert.True(t, validationDomain.ValidateString("bsnes"))
	assert.False(t, validationDomain.ValidateCore("dolphin"))
	assert.False(t, validationDomain.ValidateCore("zsnes"))

}
func TestValidationDomainValidateString(t *testing.T) {
	validationDomain, err := NewValidationDomain(testCoreWhitelist, testStringBlacklist, testIPBlacklist)
	require.NoError(t, err)

	assert.False(t, validationDomain.ValidateString("non ascii ä"))
	assert.False(t, validationDomain.ValidateString("utf-8 𝄞"))
	assert.False(t, validationDomain.ValidateString("   spaces"))
	assert.False(t, validationDomain.ValidateString("spaces   "))
	assert.True(t, validationDomain.ValidateString("mario"))
	assert.True(t, validationDomain.ValidateString("zelda"))
	assert.False(t, validationDomain.ValidateString("prefixTestZelda"))
	assert.True(t, validationDomain.ValidateString("ZeldaprefixTest"))
}

func TestValidationDomainValidateIP(t *testing.T) {
	validationDomain, err := NewValidationDomain(testCoreWhitelist, testStringBlacklist, testIPBlacklist)
	require.NoError(t, err)

	assert.True(t, validationDomain.ValdateIP(net.ParseIP("192.168.178.2")))
	assert.True(t, validationDomain.ValdateIP(net.ParseIP("8.8.8.8")))
	assert.True(t, validationDomain.ValdateIP(net.ParseIP("88.12.123.77")))
	assert.True(t, validationDomain.ValdateIP(net.ParseIP("2001:db8::1428:57ab")))
	assert.True(t, validationDomain.ValdateIP(net.ParseIP("2001:db8:0:0:0:8d3:0:0")))
	assert.False(t, validationDomain.ValdateIP(net.ParseIP("127.0.0.1")))
	assert.False(t, validationDomain.ValdateIP(net.ParseIP("2001:db8:0:8d3:0:8a2e:70:7344")))
}
