package fugu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducts(t *testing.T) {
	products := Products()
	url := "scripts.simpleanalyticscdn.com"
	product := ProductForUrl(products, url)

	assert.Equal(t, 16, len(products))
	assert.NotNil(t, product)
	assert.Equal(t, "Simple Analytics", product.Company.Name)
}

func TestCompaniesTwoUrls(t *testing.T) {
	url := "d3e54v103j8qbb.cloudfront.net"

	products := Products()
	product := ProductForUrl(products, url)

	assert.NotNil(t, product)
	assert.Equal(t, "Cloudflare", product.Company.Name)
}
