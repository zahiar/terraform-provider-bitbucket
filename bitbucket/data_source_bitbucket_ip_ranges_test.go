package bitbucket

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/stretchr/testify/assert"
)

func TestAccBitbucketIpRangesDataSource_basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t) },
		ProviderFactories: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: `data "bitbucket_ip_ranges" "testacc" {}`,
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.bitbucket_ip_ranges.testacc", "ranges.#"),

					resource.TestCheckResourceAttrSet("data.bitbucket_ip_ranges.testacc", "ranges.0.cidr"),
					resource.TestCheckResourceAttrSet("data.bitbucket_ip_ranges.testacc", "ranges.0.directions.#"),
					resource.TestCheckResourceAttrSet("data.bitbucket_ip_ranges.testacc", "ranges.0.mask"),
					resource.TestCheckResourceAttrSet("data.bitbucket_ip_ranges.testacc", "ranges.0.mask_len"),
					resource.TestCheckResourceAttrSet("data.bitbucket_ip_ranges.testacc", "ranges.0.network"),
					resource.TestCheckResourceAttrSet("data.bitbucket_ip_ranges.testacc", "ranges.0.regions.#"),
				),
			},
		},
	})
}

func TestMapToResource(t *testing.T) {
	var emptyRange []IpRangeItem
	assert.Empty(t, mapToResource(emptyRange))

	var nonBitBucketIpRange []IpRangeItem
	nonBitBucketIpRange = append(nonBitBucketIpRange, IpRangeItem{Product: []string{"notBitbucket"}})
	assert.Empty(t, mapToResource(nonBitBucketIpRange))

	var bitbucketIpRanges []IpRangeItem
	bitbucketIpRanges = append(
		bitbucketIpRanges,
		IpRangeItem{
			Network:   "123.456.789.123",
			MaskLen:   8,
			Cidr:      "123.456.789.123/8",
			Mask:      "255.255.255.255",
			Region:    []string{"eu-west-1"},
			Product:   []string{"bitbucket"},
			Direction: []string{"ingress"},
		},
		IpRangeItem{
			Network:   "456.789.123.456",
			MaskLen:   16,
			Cidr:      "456.789.123.456/16",
			Mask:      "255.255.255.255",
			Region:    []string{"eu-west-2"},
			Product:   []string{"bitbucket"},
			Direction: []string{"egress"},
		},
	)

	var expected []interface{}
	expected = append(
		expected,
		map[string]interface{}{
			"network":    "123.456.789.123",
			"mask_len":   int64(8),
			"cidr":       "123.456.789.123/8",
			"mask":       "255.255.255.255",
			"regions":    []string{"eu-west-1"},
			"directions": []string{"ingress"},
		},
	)
	expected = append(
		expected,
		map[string]interface{}{
			"network":    "456.789.123.456",
			"mask_len":   int64(16),
			"cidr":       "456.789.123.456/16",
			"mask":       "255.255.255.255",
			"regions":    []string{"eu-west-2"},
			"directions": []string{"egress"},
		},
	)
	assert.Equal(t, expected, mapToResource(bitbucketIpRanges))
}
