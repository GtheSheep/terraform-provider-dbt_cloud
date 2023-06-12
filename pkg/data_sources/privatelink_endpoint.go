package data_sources

import (
	"context"

	"github.com/gthesheep/terraform-provider-dbt-cloud/pkg/dbt_cloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var privatelinkEndpointSchema = map[string]*schema.Schema{
	"id": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The internal ID of the PrivateLink Endpoint",
	},
	"name": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "Given descriptive name for the PrivateLink Endpoint (name and/or private_link_endpoint_url need to be provided to return data for the datasource)",
	},
	"type": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Type of the PrivateLink Endpoint",
	},
	"private_link_endpoint_url": &schema.Schema{
		Type:        schema.TypeString,
		Optional:    true,
		Default:     "",
		Description: "The URL of the PrivateLink Endpoint (private_link_endpoint_url and/or name need to be provided to return data for the datasource)",
	},
	"cidr_range": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The CIDR range of the PrivateLink Endpoint",
	},
	"state": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "PrivatelinkEndpoint state should be 1 = active, as 2 = deleted",
	},
}

func DatasourcePrivatelinkEndpoint() *schema.Resource {
	return &schema.Resource{
		ReadContext: datasourcePrivatelinkEndpointRead,
		Schema:      privatelinkEndpointSchema,
	}
}

func datasourcePrivatelinkEndpointRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	endpointName := d.Get("name").(string)
	privatelinkEndpointURL := d.Get("private_link_endpoint_url").(string)

	privatelinkEndpoint, err := c.GetPrivatelinkEndpoint(endpointName, privatelinkEndpointURL)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("id", privatelinkEndpoint.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", privatelinkEndpoint.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", privatelinkEndpoint.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("private_link_endpoint_url", privatelinkEndpoint.PrivatelinkEndpointURL); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("cidr_range", privatelinkEndpoint.CIDRRange); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("state", privatelinkEndpoint.State); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(privatelinkEndpoint.ID)

	return diags
}
