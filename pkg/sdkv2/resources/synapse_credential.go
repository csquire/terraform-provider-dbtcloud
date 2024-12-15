package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/dbt_cloud"
	"github.com/dbt-labs/terraform-provider-dbtcloud/pkg/helper"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func ResourceSynapseCredential() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceSynapseCredentialCreate,
		ReadContext:   resourceSynapseCredentialRead,
		UpdateContext: resourceSynapseCredentialUpdate,
		DeleteContext: resourceSynapseCredentialDelete,

		Schema: map[string]*schema.Schema{
			"project_id": {
				Type:        schema.TypeInt,
				Required:    true,
				ForceNew:    true,
				Description: "Project ID to create the Synapse credential in",
			},
			"credential_id": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The system Synapse credential ID",
			},
			"user": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				Description:   "The username of the Synapse account to connect to. Only used when connection with AD user/pass",
				ConflictsWith: []string{"tenant_id", "client_id", "client_secret"},
			},
			"password": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Default:       "",
				Description:   "The password for the account to connect to. Only used when connection with AD user/pass",
				ConflictsWith: []string{"tenant_id", "client_id", "client_secret"},
			},
			"tenant_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				Description:   "The tenant ID of the Azure Active Directory instance. This is only used when connecting to Azure SQL with a service principal.",
				ConflictsWith: []string{"user", "password"},
			},
			"client_id": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				Description:   "The client ID of the Azure Active Directory service principal. This is only used when connecting to Azure SQL with an AAD service principal.",
				ConflictsWith: []string{"user", "password"},
			},
			"client_secret": {
				Type:          schema.TypeString,
				Optional:      true,
				Sensitive:     true,
				Default:       "",
				Description:   "The client secret of the Azure Active Directory service principal. This is only used when connecting to Azure SQL with an AAD service principal.",
				ConflictsWith: []string{"user", "password"},
			},
			"schema": {
				Type:          schema.TypeString,
				Required:      true,
				Description:   "The schema where to create the dbt models",
				ConflictsWith: []string{},
			},
			"schema_authorization": {
				Type:          schema.TypeString,
				Optional:      true,
				Default:       "",
				Description:   "Optionally set this to the principal who should own the schemas created by dbt",
				ConflictsWith: []string{},
			},
		},

		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceSynapseCredentialCreate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectId := d.Get("project_id").(int)
	user := d.Get("user").(string)
	password := d.Get("password").(string)
	tenantId := d.Get("tenant_id").(string)
	clientId := d.Get("client_id").(string)
	clientSecret := d.Get("client_secret").(string)
	schema := d.Get("schema").(string)
	schemaAuthorization := d.Get("schema_authorization").(string)

	// FRAMEWORK: Move this logic to the schema validation when moving to the Framework
	userPasswordDefined := user != "" && password != ""
	servicePrincipalDefined := tenantId != "" && clientId != "" && clientSecret != ""

	if !userPasswordDefined && !servicePrincipalDefined {
		diag.FromErr(fmt.Errorf("either user/password or service principal auth must be defined"))
	}

	synapseCredential, err := c.CreateSynapseCredential(
		projectId,
		user,
		password,
		tenantId,
		clientId,
		clientSecret,
		schema,
		schemaAuthorization,
	)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(
		fmt.Sprintf(
			"%d%s%d",
			synapseCredential.Project_Id,
			dbt_cloud.ID_DELIMITER,
			*synapseCredential.ID,
		),
	)

	resourceSynapseCredentialRead(ctx, d, m)

	return diags
}

func resourceSynapseCredentialRead(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	// Warning or errors can be collected in a slice type
	var diags diag.Diagnostics

	projectId, synapseCredentialId, err := helper.SplitIDToInts(
		d.Id(),
		"dbtcloud_synapse_credential",
	)
	if err != nil {
		return diag.FromErr(err)
	}

	synapseCredential, err := c.GetSynapseCredential(projectId, synapseCredentialId)
	if err != nil {
		if strings.HasPrefix(err.Error(), "resource-not-found") {
			d.SetId("")
			return diags
		}
		return diag.FromErr(err)
	}

	// set the ones that don't come back from the API

	if err := d.Set("project_id", synapseCredential.Project_Id); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("credential_id", synapseCredentialId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("user", synapseCredential.UnencryptedCredentialDetails.User); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("tenant_id", synapseCredential.UnencryptedCredentialDetails.TenantId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_id", synapseCredential.UnencryptedCredentialDetails.ClientId); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema", synapseCredential.UnencryptedCredentialDetails.Schema); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("schema_authorization", synapseCredential.UnencryptedCredentialDetails.SchemaAuthorization); err != nil {
		return diag.FromErr(err)
	}

	// set the ones that don't come back from the API
	if err := d.Set("password", d.Get("password").(string)); err != nil {
		return diag.FromErr(err)
	}
	if err := d.Set("client_secret", d.Get("client_secret").(string)); err != nil {
		return diag.FromErr(err)
	}

	return diags
}

func resourceSynapseCredentialUpdate(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	projectId, synapseCredentialId, err := helper.SplitIDToInts(
		d.Id(),
		"dbtcloud_synapse_credential",
	)
	if err != nil {
		return diag.FromErr(err)
	}

	if d.HasChange("user") ||
		d.HasChange("password") ||
		d.HasChange("tenant_id") ||
		d.HasChange("client_id") ||
		d.HasChange("client_secret") ||
		d.HasChange("schema") ||
		d.HasChange("schema_authorization") {

		user := d.Get("user").(string)
		password := d.Get("password").(string)
		tenantId := d.Get("tenant_id").(string)
		clientId := d.Get("client_id").(string)
		clientSecret := d.Get("client_secret").(string)
		schema := d.Get("schema").(string)
		schemaAuthorization := d.Get("schema_authorization").(string)

		// FRAMEWORK: Move this logic to the schema validation when moving to the Framework
		userPasswordDefined := user != "" && password != ""
		servicePrincipalDefined := tenantId != "" && clientId != "" && clientSecret != ""

		if !userPasswordDefined && !servicePrincipalDefined {
			return diag.FromErr(
				fmt.Errorf("either user/password or service principal auth must be defined"),
			)
		}

		synapseCredential, err := c.GetSynapseCredential(projectId, synapseCredentialId)
		if err != nil {
			return diag.FromErr(err)
		}

		synapseCredentialDetails, err := dbt_cloud.GenerateSynapseCredentialDetails(
			user,
			password,
			tenantId,
			clientId,
			clientSecret,
			schema,
			schemaAuthorization,
		)
		if err != nil {
			return diag.FromErr(err)
		}

		synapseCredential.CredentialDetails = synapseCredentialDetails

		_, err = c.UpdateSynapseCredential(
			projectId,
			synapseCredentialId,
			*synapseCredential,
		)
		if err != nil {
			return diag.FromErr(err)
		}
	}

	return resourceSynapseCredentialRead(ctx, d, m)
}

func resourceSynapseCredentialDelete(
	ctx context.Context,
	d *schema.ResourceData,
	m interface{},
) diag.Diagnostics {
	c := m.(*dbt_cloud.Client)

	var diags diag.Diagnostics

	projectId, synapseCredentialId, err := helper.SplitIDToInts(
		d.Id(),
		"dbtcloud_synapse_credential",
	)
	if err != nil {
		return diag.FromErr(err)
	}

	synapseCredential, err := c.GetSynapseCredential(projectId, synapseCredentialId)
	if err != nil {
		return diag.FromErr(err)
	}

	synapseCredential.State = dbt_cloud.STATE_DELETED

	// These values don't mean anything for the delete operation but they are required by the API
	// TODO(cwalden): do we need to handle this err?
	emptySynapseCredentialDetails, _ := dbt_cloud.GenerateSynapseCredentialDetails(
		"",
		"",
		"",
		"",
		"",
		"",
		"",
	)

	synapseCredential.CredentialDetails = emptySynapseCredentialDetails

	_, err = c.UpdateSynapseCredential(projectId, synapseCredentialId, *synapseCredential)
	if err != nil {
		return diag.FromErr(err)
	}

	return diags
}
