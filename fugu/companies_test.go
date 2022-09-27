package fugu

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProducts(t *testing.T) {
	products := Products()
	url := "scripts.simpleanalyticscdn.com"
	product := ProductForUrl(products, url)

	assert.Equal(t, 21, len(products))
	assert.NotNil(t, product)
	assert.Equal(t, "Simple Analytics", product.Company.Name)
}

func TestLongestUrl(t *testing.T) {
	products := Products()
	url1 := "fonts.googleapis.com"
	product1 := ProductForUrl(products, url1)
	assert.NotNil(t, product1)
	assert.Equal(t, "Google", product1.Company.Name)
	assert.Equal(t, "Google Fonts", product1.Name)

	url2 := "googleapis.com"
	product2 := ProductForUrl(products, url2)
	assert.NotNil(t, product2)
	assert.Equal(t, "Google", product2.Company.Name)
	assert.Equal(t, "Google", product2.Name)
}

func TestCompaniesTwoUrls(t *testing.T) {
	url := "d3e54v103j8qbb.cloudfront.net"

	products := Products()
	product := ProductForUrl(products, url)

	assert.NotNil(t, product)
	assert.Equal(t, "Cloudflare", product.Company.Name)
}
