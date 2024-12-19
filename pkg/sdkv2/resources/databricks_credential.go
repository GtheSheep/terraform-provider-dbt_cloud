package resources

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/dbt_cloud"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	adapterTypes = []string{
		"databricks",
		"spark",
	}
)

func isLegacyDatabricksConnection(d *schema.ResourceData) bool {
	return d.Get("adapter_id").(int) != 0
}

func ResourceDatabricksCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceDatabricksCredentialCreate,
		ReadContext:   resourceDatabricksCredentialRead,
		UpdateContext: resourceDatabricksCredentialUpdate,
		DeleteContext: resourceDatabricksCredentialDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID to create the Databricks credential in",
			},
			"adapter_id": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "Databricks adapter ID for the credential (do not fill in when using global connections, only to be used for connections created with the legacy connection resource `dbtcloud_connection`)",
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The system Databricks credential ID",
			},
			"target_name": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "default",
				Description: "Target name",
				Deprecated:  "This field is deprecated at the environment level (it was never possible to set it in the UI) and will be removed in a future release. Please remove it and set the target name at the job level or leverage environment variables.",
			},
			"token": {
				Type:        schema.TypeString,
				Required:    true,
				Sensitive:   true,
				Description: "Token for Databricks user",
			},
			"catalog": {
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "The catalog where to create models (only for the databricks adapter)",
			},
			"schema": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The schema where to create models",
			},
			"adapter_type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The type of the adapter (databricks or spark)",
				ValidateFunc: validation.StringInSlice(adapterTypes, false),
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceDatabricksCredentialCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	if isLegacyDatabricksConnection(d) {
		return resourceDatabricksCredentialCreateLegacy(ctx, d, m)
	} else {
		return resourceDatabricksCredentialCreateGlobConn(ctx, d, m)
	}
}

func resourceDatabricksCredentialCreateLegacy(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectId := d.Get("project_id").(int)
	adapterId := d.Get("adapter_id").(int)
	targetName := d.Get("target_name").(string)
	token := d.Get("token").(string)
	catalog := d.Get("catalog").(string)
	schema := d.Get("schema").(string)
	adapterType := d.Get("adapter_type").(string)

	databricksCredential, err := c.CreateDatabricksCredentialLegacy(
		projectId,
		"adapter",
		targetName,
		adapterId,
		token,
		catalog,
		schema,
		adapterType,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(
		fmt.Sprintf(
			"%d%s%d",
			databricksCredential.Project_Id,
			dbt_cloud.ID_DELIMITER,
			*databricksCredential.ID,
		),
	)

	resourceDatabricksCredentialRead(ctx, d, m)

	return diags
}

func resourceDatabricksCredentialCreateGlobConn(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	projectId := d.Get("project_id").(int)
	targetName := d.Get("target_name").(string)
	token := d.Get("token").(string)
	catalog := d.Get("catalog").(string)
	schema := d.Get("schema").(string)
	adapterType := d.Get("adapter_type").(string)

	// for now, just supporting databricks
	if adapterType == "spark" {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Spark adapter is not supported currently for global connections credentials. Please raise a GitHub issue if you need it",
		})
		return diags
	}

	databricksCredential, err := c.CreateDatabricksCredential(
		projectId,
		token,
		schema,
		targetName,
		catalog,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(
		fmt.Sprintf(
			"%d%s%d",
			databricksCredential.Project_Id,
			dbt_cloud.ID_DELIMITER,
			*databricksCredential.ID,
		),
	)

	resourceDatabricksCredentialRead(ctx, d, m)

	return diags
}

func resourceDatabricksCredentialRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics
	projectId, databricksCredentialId, err := helper.SplitIDToInts(
		d.Id(),
		"dbtcloud_databricks_credential",
	)
	if err != nil {
		return diag.FromErr(err)
	}

	databricksCredential, err := c.GetDatabricksCredential(projectId, databricksCredentialId)
	if err != nil {
		if strings.HasPrefix(err.Error(), "resource-not-found") {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	if err := d.Set("credential_id", databricksCredentialId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_id", databricksCredential.Project_Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("adapter_id", databricksCredential.Adapter_Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("target_name", databricksCredential.Target_Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("token", d.Get("token").(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("catalog", databricksCredential.UnencryptedCredentialDetails.Catalog); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", databricksCredential.UnencryptedCredentialDetails.Schema); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("adapter_type", d.Get("adapter_type").(string)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDatabricksCredentialUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	if isLegacyDatabricksConnection(d) {
		return resourceDatabricksCredentialUpdateLegacy(ctx, d, m)
	} else {
		return resourceDatabricksCredentialUpdateGlobConn(ctx, d, m)
	}
}

func resourceDatabricksCredentialUpdateLegacy(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)
	projectId, databricksCredentialId, err := helper.SplitIDToInts(
		d.Id(),
		"dbtcloud_databricks_credential",
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("adapter_id") || d.HasChange("token") || d.HasChange("target_name") ||
		d.HasChange("catalog") ||
		d.HasChange("schema") ||
		d.HasChange("adapter_type") {
		databricksCredential, err := c.GetDatabricksCredential(projectId, databricksCredentialId)
		if err != nil {
			return diag.FromErr(err)
		}
		if d.HasChange("adapter_id") {
			adapterId := d.Get("adapter_id").(int)
			databricksCredential.Adapter_Id = adapterId
		}
		if d.HasChange("target_name") {
			targetName := d.Get("target_name").(string)
			databricksCredential.Target_Name = targetName
		}

		// we need to fill in the DatabricksCredentialFieldMetadata for all fields, except token if it was not changed
		validation := dbt_cloud.AdapterCredentialFieldMetadataValidation{
			Required: false,
		}
		tokenMetadata := dbt_cloud.AdapterCredentialFieldMetadata{
			Label:       "Token",
			Description: "Personalized user token.",
			Field_Type:  "text",
			Encrypt:     true,
			Validation:  validation,
		}
		catalogMetadata := dbt_cloud.AdapterCredentialFieldMetadata{
			Label:       "Catalog",
			Description: "Catalog name if Unity Catalog is enabled in your Databricks workspace.  Only available in dbt version 1.1 and later.",
			Field_Type:  "text",
			Encrypt:     false,
			Validation:  validation,
		}
		schemaMetadata := dbt_cloud.AdapterCredentialFieldMetadata{
			Label:       "Schema",
			Description: "User schema.",
			Field_Type:  "text",
			Encrypt:     false,
			Validation:  validation,
		}
		threadsMetadata := dbt_cloud.AdapterCredentialFieldMetadata{
			Label:       "Threads",
			Description: "The number of threads to use for your jobs.",
			Field_Type:  "number",
			Encrypt:     false,
			Validation:  validation,
		}

		credentialsFieldToken := dbt_cloud.AdapterCredentialField{
			Metadata: tokenMetadata,
			Value:    d.Get("token").(string),
		}
		credentialsFieldCatalog := dbt_cloud.AdapterCredentialField{
			Metadata: catalogMetadata,
			Value:    d.Get("catalog").(string),
		}
		credentialsFieldSchema := dbt_cloud.AdapterCredentialField{
			Metadata: schemaMetadata,
			Value:    d.Get("schema").(string),
		}
		credentialsFieldThreads := dbt_cloud.AdapterCredentialField{
			Metadata: threadsMetadata,
			Value:    dbt_cloud.NUM_THREADS_CREDENTIAL,
		}

		credentialFields := map[string]dbt_cloud.AdapterCredentialField{}

		// only databricks accepts a catalog, not spark
		if d.Get("adapter_type").(string) == "databricks" {
			credentialFields["catalog"] = credentialsFieldCatalog

			// for databricks, we update token only if it was changed
			if d.HasChange("token") {
				credentialFields["token"] = credentialsFieldToken
			}
		}

		// spark requires sending all the details
		if d.Get("adapter_type").(string) == "spark" {
			credentialFields["token"] = credentialsFieldToken
			credentialFields["threads"] = credentialsFieldThreads
			credentialFields["schema"] = credentialsFieldSchema
		}

		credentialFields["schema"] = credentialsFieldSchema

		credentialDetails := dbt_cloud.AdapterCredentialDetails{
			Fields:      credentialFields,
			Field_Order: []string{},
		}

		databricksCredential.Credential_Details = credentialDetails

		_, err = c.UpdateDatabricksCredentialLegacy(
			projectId,
			databricksCredentialId,
			*databricksCredential,
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDatabricksCredentialRead(ctx, d, m)
}

func resourceDatabricksCredentialUpdateGlobConn(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)
	projectId, databricksCredentialId, err := helper.SplitIDToInts(
		d.Id(),
		"dbtcloud_databricks_credential",
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("token") ||
		d.HasChange("target_name") ||
		d.HasChange("catalog") ||
		d.HasChange("schema") {

		patchCredentialsDetails, err := dbt_cloud.GenerateDatabricksCredentialDetails(
			d.Get("token").(string),
			d.Get("schema").(string),
			d.Get("target_name").(string),
			d.Get("catalog").(string),
		)

		for key := range patchCredentialsDetails.Fields {
			if d.Get(key) == nil || !d.HasChange(key) {
				delete(patchCredentialsDetails.Fields, key)
			}
		}

		databricksPatch := dbt_cloud.DatabricksCredentialGLobConnPatch{
			ID:                databricksCredentialId,
			CredentialDetails: patchCredentialsDetails,
		}

		_, err = c.UpdateDatabricksCredentialGlobConn(
			projectId,
			databricksCredentialId,
			databricksPatch,
		)

		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceDatabricksCredentialRead(ctx, d, m)
}

func resourceDatabricksCredentialDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	if isLegacyDatabricksConnection(d) {
		return resourceDatabricksCredentialDeleteLegacy(ctx, d, m)
	} else {
		return resourceDatabricksCredentialDeleteGlobConn(ctx, d, m)
	}
}

func resourceDatabricksCredentialDeleteLegacy(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	projectId, databricksCredentialId, err := helper.SplitIDToInts(
		d.Id(),
		"dbtcloud_databricks_credential",
	)
	if err != nil {
		return diag.FromErr(err)
	}

	databricksCredential, err := c.GetDatabricksCredential(projectId, databricksCredentialId)
	if err != nil {
		return diag.FromErr(err)
	}

	databricksCredential.State = dbt_cloud.STATE_DELETED

	// those values don't mean anything for delete operation but they are required by the API
	validation := dbt_cloud.AdapterCredentialFieldMetadataValidation{
		Required: false,
	}
	catalogMetadata := dbt_cloud.AdapterCredentialFieldMetadata{
		Label:       "Catalog",
		Description: "Catalog name if Unity Catalog is enabled in your Databricks workspace.  Only available in dbt version 1.1 and later.",
		Field_Type:  "text",
		Encrypt:     false,
		Validation:  validation,
	}
	schemaMetadata := dbt_cloud.AdapterCredentialFieldMetadata{
		Label:       "Schema",
		Description: "User schema.",
		Field_Type:  "text",
		Encrypt:     false,
		Validation:  validation,
	}
	credentialsFieldCatalog := dbt_cloud.AdapterCredentialField{
		Metadata: catalogMetadata,
		Value:    "NA",
	}
	credentialsFieldSchema := dbt_cloud.AdapterCredentialField{
		Metadata: schemaMetadata,
		Value:    "NA",
	}
	credentialFields := map[string]dbt_cloud.AdapterCredentialField{}
	credentialFields["catalog"] = credentialsFieldCatalog
	credentialFields["schema"] = credentialsFieldSchema

	credentialDetails := dbt_cloud.AdapterCredentialDetails{
		Fields:      credentialFields,
		Field_Order: []string{},
	}

	databricksCredential.Credential_Details = credentialDetails

	_, err = c.UpdateDatabricksCredentialLegacy(
		projectId,
		databricksCredentialId,
		*databricksCredential,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceDatabricksCredentialDeleteGlobConn(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	projectId, databricksCredentialId, err := helper.SplitIDToInts(
		d.Id(),
		"dbtcloud_databricks_credential",
	)
	if err != nil {
		return diag.FromErr(err)
	}

	_, err = c.DeleteCredential(
		strconv.Itoa(databricksCredentialId),
		strconv.Itoa(projectId),
	)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
