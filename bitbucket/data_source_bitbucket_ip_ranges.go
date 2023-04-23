package bitbucket

import (
	"context"
	"encoding/json"
	"fmt"
	"golang.org/x/exp/slices"
	"net/http"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type IpRanges struct {
	CreationDate string        `json:"creationDate"`
	SyncToken    int64         `json:"syncToken"`
	Items        []IpRangeItem `json:"items"`
}

type IpRangeItem struct {
	Network   string   `json:"network"`
	MaskLen   int64    `json:"mask_len"`
	Cidr      string   `json:"cidr"`
	Mask      string   `json:"mask"`
	Region    []string `json:"region"`
	Product   []string `json:"product"`
	Direction []string `json:"direction"`
}

func dataSourceBitbucketIpRanges() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceBitbucketIpRangesRead,
		Schema: map[string]*schema.Schema{
			"ranges": {
				Type:     schema.TypeSet,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"network": {
							Description: "The IP Address.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"cidr": {
							Description: "The CIDR.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"mask": {
							Description: "The Mask.",
							Type:        schema.TypeString,
							Computed:    true,
						},
						"mask_len": {
							Description: "The Mask Length.",
							Type:        schema.TypeInt,
							Computed:    true,
						},
						"regions": {
							Description: "A list of regions this IP Address resides in. Follows AWS Region Code Notation.",
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
						"directions": {
							Description: "A list defining if the IP address is ingress, egress or both.",
							Type:        schema.TypeSet,
							Elem:        &schema.Schema{Type: schema.TypeString},
							Computed:    true,
						},
					},
				},
			},
		},
	}
}

func dataSourceBitbucketIpRangesRead(ctx context.Context, resourceData *schema.ResourceData, meta interface{}) diag.Diagnostics {
	req, err := http.Get("https://ip-ranges.atlassian.com/")
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to get ip ranges with error: %s", err))
	}
	defer req.Body.Close()

	if req.StatusCode != 200 {
		return diag.FromErr(fmt.Errorf("unable to get ip ranges: response code was not 200"))
	}

	var ipRanges IpRanges
	err = json.NewDecoder(req.Body).Decode(&ipRanges)
	if err != nil {
		return diag.FromErr(fmt.Errorf("unable to decode ip ranges with error: %s", err))
	}

	resourceData.SetId(fmt.Sprintf("%d", ipRanges.SyncToken))
	_ = resourceData.Set("ranges", mapToResource(ipRanges.Items))

	return nil
}

func mapToResource(ipRanges []IpRangeItem) []interface{} {
	var resourceList []interface{}

	for _, item := range ipRanges {
		if slices.Contains(item.Product, "bitbucket") {
			resourceList = append(resourceList, map[string]interface{}{
				"network":    item.Network,
				"mask_len":   item.MaskLen,
				"cidr":       item.Cidr,
				"mask":       item.Mask,
				"regions":    item.Region,
				"directions": item.Direction,
			})
		}
	}

	if len(resourceList) == 0 {
		return nil
	}

	return resourceList
}
