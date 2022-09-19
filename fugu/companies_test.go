package fugu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCompanies(t *testing.T) {
	companies := Companies()
	url := "scripts.simpleanalyticscdn.com"
	company := CompanyForUrl(companies, url)

	assert.Equal(t, 9, len(companies))
	assert.NotNil(t, company)
	assert.Equal(t, "Simple Analytics", company.Name)
}

func TestCompaniesTwoUrls(t *testing.T) {
	url := "d3e54v103j8qbb.cloudfront.net"

	companies := Companies()
	company := CompanyForUrl(companies, url)

	assert.Equal(t, 9, len(companies))
	assert.NotNil(t, company)
	assert.Equal(t, "Cloudflare", company.Name)
}
