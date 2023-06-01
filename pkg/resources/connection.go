package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/gthesheep/terraform-provider-dbt-cloud/pkg/dbt_cloud"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"
)

var (
	connectionTypes = []string{
		"snowflake",
		"bigquery",
		"redshift",
		"postgres",
		"alloydb",
		"adapter",
	}
)

func ResourceConnection() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceConnectionCreate,
		ReadContext:   resourceConnectionRead,
		UpdateContext: resourceConnectionUpdate,
		DeleteContext: resourceConnectionDelete,

		Schema: map[string]*schema.Schema{
			"connection_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Connection Identifier",
			},
			"is_active": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     true,
				Description: "Whether the connection is active",
			},
			"project_id": &schema.Schema{
				Type:        schema.TypeInt,
				Required:    true,
				Description: "Project ID to create the connection in",
			},
			"name": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Connection name",
			},
			"type": &schema.Schema{
				Type:         schema.TypeString,
				Required:     true,
				Description:  "The type of connection",
				ValidateFunc: validation.StringInSlice(connectionTypes, false),
			},
			"account": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Account name for the connection",
				ConflictsWith: []string{"host_name"},
			},
			"host_name": &schema.Schema{
				Type:          schema.TypeString,
				Optional:      true,
				Description:   "Host name for the connection, including Databricks cluster",
				ConflictsWith: []string{"account"},
			},
			"port": &schema.Schema{
				Type:        schema.TypeInt,
				Optional:    true,
				Default:     nil,
				Description: "Port number to connect via",
			},
			"database": &schema.Schema{
				Type:        schema.TypeString,
				Required:    true,
				Description: "Database name for the connection",
			},
			"warehouse": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Warehouse name for the connection",
			},
			"role": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "Role name for the connection",
			},
			"allow_sso": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the connection should allow SSO",
			},
			"allow_keep_alive": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not the connection should allow client session keep alive",
			},
			"oauth_client_id": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "OAuth client identifier",
			},
			"oauth_client_secret": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Default:     "",
				Description: "OAuth client secret",
			},
			"tunnel_enabled": &schema.Schema{
				Type:        schema.TypeBool,
				Optional:    true,
				Default:     false,
				Description: "Whether or not tunneling should be enabled on your database connection",
			},
			"http_path": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The HTTP path of the Databricks cluster or SQL warehouse",
			},
			"catalog": &schema.Schema{
				Type:        schema.TypeString,
				Optional:    true,
				Description: "Catalog name if Unity Catalog is enabled in your Databricks workspace",
			},
			"adapter_id": &schema.Schema{
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Adapter id created for the Databricks connection",
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceConnectionCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	isActive := d.Get("is_active").(bool)
	projectId := d.Get("project_id").(int)
	name := d.Get("name").(string)
	connectionType := d.Get("type").(string)
	account := d.Get("account").(string)
	database := d.Get("database").(string)
	warehouse := d.Get("warehouse").(string)
	role := d.Get("role").(string)
	allowSSO := d.Get("allow_sso").(bool)
	allowKeepAlive := d.Get("allow_keep_alive").(bool)
	oAuthClientID := d.Get("oauth_client_id").(string)
	oAuthClientSecret := d.Get("oauth_client_secret").(string)
	hostName := d.Get("host_name").(string)
	port := d.Get("port").(int)
	tunnelEnabled := d.Get("tunnel_enabled").(bool)
	httpPath := d.Get("http_path").(string)
	catalog := d.Get("catalog").(string)

	connection, err := c.CreateConnection(projectId, name, connectionType, isActive, account, database, warehouse, role, &allowSSO, &allowKeepAlive, oAuthClientID, oAuthClientSecret, hostName, port, &tunnelEnabled, httpPath, catalog)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fmt.Sprintf("%d%s%d", connection.ProjectID, dbt_cloud.ID_DELIMITER, *connection.ID))

	resourceConnectionRead(ctx, d, m)

	return diags
}

func resourceConnectionRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	projectIdString := strings.Split(d.Id(), dbt_cloud.ID_DELIMITER)[0]
	connectionIdString := strings.Split(d.Id(), dbt_cloud.ID_DELIMITER)[1]

	connection, err := c.GetConnection(connectionIdString, projectIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	// TODO: Remove when done better
	connection.Details.OAuthClientID = d.Get("oauth_client_id").(string)
	connection.Details.OAuthClientSecret = d.Get("oauth_client_secret").(string)

	if err := d.Set("connection_id", connection.ID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("is_active", connection.State == dbt_cloud.STATE_ACTIVE); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("project_id", connection.ProjectID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("name", connection.Name); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("type", connection.Type); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("account", connection.Details.Account); err != nil {
		return diag.FromErr(err)
	}

	databaseName := connection.Details.DBName
	if connection.Type == "snowflake" {
		databaseName = connection.Details.Database
	}
	if err := d.Set("database", databaseName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("warehouse", connection.Details.Warehouse); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("role", connection.Details.Role); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allow_sso", connection.Details.AllowSSO); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("allow_keep_alive", connection.Details.ClientSessionKeepAlive); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("oauth_client_id", connection.Details.OAuthClientID); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("oauth_client_secret", connection.Details.OAuthClientSecret); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("port", connection.Details.Port); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tunnel_enabled", connection.Details.TunnelEnabled); err != nil {
		return diag.FromErr(err)
	}
	httpPath := ""
	catalog := ""
	hostName := connection.Details.Host
	if connection.Details.AdapterDetails != nil {
		httpPath = connection.Details.AdapterDetails.Fields["http_path"].Value.(string)
		catalog = connection.Details.AdapterDetails.Fields["catalog"].Value.(string)
		hostName = connection.Details.AdapterDetails.Fields["host"].Value.(string)
	}
	if err := d.Set("host_name", hostName); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("http_path", httpPath); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("catalog", catalog); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("adapter_id", connection.Details.AdapterId); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceConnectionUpdate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	projectIdString := strings.Split(d.Id(), dbt_cloud.ID_DELIMITER)[0]
	connectionIdString := strings.Split(d.Id(), dbt_cloud.ID_DELIMITER)[1]

	// TODO: add more changes here

	if d.HasChange("name") || d.HasChange("type") || d.HasChange("account") || d.HasChange("database") || d.HasChange("warehouse") || d.HasChange("role") || d.HasChange("allow_sso") || d.HasChange("allow_keep_alive") || d.HasChange("oauth_client_id") || d.HasChange("oauth_client_secret") {
		connection, err := c.GetConnection(connectionIdString, projectIdString)
		if err != nil {
			return diag.FromErr(err)
		}

		if d.HasChange("name") {
			name := d.Get("name").(string)
			connection.Name = name
		}
		if d.HasChange("type") {
			connectionType := d.Get("type").(string)
			connection.Type = connectionType
		}
		if d.HasChange("account") {
			account := d.Get("account").(string)
			connection.Details.Account = account
		}
		if d.HasChange("database") {
			database := d.Get("database").(string)
			if connection.Type == "snowflake" {
				connection.Details.Database = database
			} else {
				connection.Details.DBName = database
			}
		}
		if d.HasChange("warehouse") {
			warehouse := d.Get("warehouse").(string)
			connection.Details.Warehouse = warehouse
		}
		if d.HasChange("role") {
			role := d.Get("role").(string)
			connection.Details.Role = role
		}
		if d.HasChange("allow_sso") {
			allowSSO := d.Get("allow_sso").(bool)
			connection.Details.AllowSSO = &allowSSO
		}
		if d.HasChange("allow_keep_alive") {
			allowKeepAlive := d.Get("allow_keep_alive").(bool)
			connection.Details.ClientSessionKeepAlive = &allowKeepAlive
		}
		if d.HasChange("oauth_client_id") {
			oAuthClientID := d.Get("oauth_client_id").(string)
			connection.Details.OAuthClientID = oAuthClientID
		}
		if d.HasChange("oauth_client_secret") {
			oAuthClientSecret := d.Get("oauth_client_secret").(string)
			connection.Details.OAuthClientSecret = oAuthClientSecret
		}
		if d.HasChange("host_name") {
			hostName := d.Get("host_name").(string)
			connection.Details.Host = hostName
		}
		if d.HasChange("port") {
			port := d.Get("port").(int)
			connection.Details.Port = port
		}

		_, err = c.UpdateConnection(connectionIdString, projectIdString, *connection)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceConnectionRead(ctx, d, m)
}

func resourceConnectionDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	projectIdString := strings.Split(d.Id(), dbt_cloud.ID_DELIMITER)[0]
	connectionIdString := strings.Split(d.Id(), dbt_cloud.ID_DELIMITER)[1]

	_, err := c.DeleteConnection(connectionIdString, projectIdString)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
