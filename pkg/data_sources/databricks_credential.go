package data_sources

import (
	"context"
	"fmt"

	"github.com/gthesheep/terraform-provider-dbt-cloud/pkg/dbt_cloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var databricksCredentialSchema = map[string]*schema.Schema{
	"project_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Project ID",
	},
	"credential_id": &schema.Schema{
		Type:        schema.TypeInt,
		Required:    true,
		Description: "Credential ID",
	},
	"adapter_id": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Databricks adapter ID for the credential",
	},
	"target_name": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "Target name",
	},
	"num_threads": &schema.Schema{
		Type:        schema.TypeInt,
		Computed:    true,
		Description: "Number of threads to use",
	},
	"catalog": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The catalog where to create models",
	},
	"schema": &schema.Schema{
		Type:        schema.TypeString,
		Computed:    true,
		Description: "The schema where to create models",
	},
}

func DatasourceDatabricksCredential() *schema.Resource {
	return &schema.Resource{
		ReadContext: databricksCredentialRead,
		Schema:      databricksCredentialSchema,
	}
}

func databricksCredentialRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	credentialID := d.Get("credential_id").(int)
	projectID := d.Get("project_id").(int)

	databricksCredential, err := c.GetDatabricksCredential(projectID, credentialID)
	if err != nil {
		return diag.FromErr(err)
	}

	if err := d.Set("adapter_id", databricksCredential.Adapter_Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_id", databricksCredential.Project_Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("target_name", databricksCredential.Target_Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("num_threads", databricksCredential.Threads); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("catalog", databricksCredential.UnencryptedCredentialDetails["catalog"]); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", databricksCredential.UnencryptedCredentialDetails["schema"]); err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d%s%d", databricksCredential.Project_Id, dbt_cloud.ID_DELIMITER, *databricksCredential.ID))

	return diags
}
