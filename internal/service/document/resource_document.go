// Copyright (c) 01Joseph-Hwang10
// SPDX-License-Identifier: MPL-2.0

package document

import (
	"context"

	errs "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/error"
	"github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/mongoclient"
	resourceconfig "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/config"
	resourceid "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/resource/id"
	mdutils "github.com/01Joseph-Hwang10/terraform-provider-mongodb/internal/common/string/markdown"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/boolplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces.
var _ resource.Resource = &DocumentResource{}
var _ resource.ResourceWithImportState = &DocumentResource{}

func NewDocumentResource() resource.Resource {
	return &DocumentResource{}
}

// DocumentResource defines the resource implementation.
type DocumentResource struct {
	config *resourceconfig.ResourceConfig
}

// DocumentResourceModel describes the resource data model.
type DocumentResourceModel struct {
	Id               types.String `tfsdk:"id"`
	Database         types.String `tfsdk:"database"`
	Collection       types.String `tfsdk:"collection"`
	DocumentId       types.String `tfsdk:"document_id"`
	Document         types.String `tfsdk:"document"`
	SyncWithDatabase types.Bool   `tfsdk:"sync_with_database"`
}

func (r *DocumentResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_database_document"
}

func (r *DocumentResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: mdutils.FormatResourceDescription(`
			This resource creates a single document in a collection 
			in a database on the MongoDB server.
		`),

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
				MarkdownDescription: mdutils.FormatSchemaDescription(
					`
						Resource identifier.
						
						ID has a value with a format of the following:

						%s

						Note that this format is used for importing the resource into Terraform state.
						Import the resource using the following command:

						%s
					`,
					mdutils.CodeBlock("", "databases/<database>/collections/<name>/documents/<document_id>"),
					mdutils.CodeBlock("bash", "terraform import mongodb_database_document.<resource_name> databases/<database>/collections/<name>/documents/<document_id>"),
				),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"database": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the database to create the collection in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"collection": schema.StringAttribute{
				Required:            true,
				MarkdownDescription: "Name of the collection to create the document in.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"document_id": schema.StringAttribute{
				MarkdownDescription: "Document ID of the document.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"document": schema.StringAttribute{
				MarkdownDescription: mdutils.FormatSchemaDescription(
					`
						Document to insert into the collection.

						The value of this attribute is a stringified JSON.
						Note that you should escape every double quote in the JSON string.

						In terraform, you can achieve this by simply using the 
						%s function:

						%s

						[EJSON](https://www.mongodb.com/docs/manual/reference/mongodb-extended-json/#std-label-mongodb-extended-json-v2) is supported in this attribute.
					`,
					mdutils.InlineCodeBlock("jsonencode"),
					mdutils.CodeBlock("terraform", "document = jsonencode({ key = \"value\" })"),
				),
				Required: true,
			},
			"sync_with_database": schema.BoolAttribute{
				Computed: true,
				Optional: true,
				Default:  booldefault.StaticBool(true),
				MarkdownDescription: mdutils.FormatSchemaDescription(`
					If this option is true, the provider will ensure that 
					the document in the Terraform state is in sync with the
					document in the database. In other words, it will ensure
					the data consistency between the document in the Terraform
					state and the document in the database. This means that 
					the provider will fail to go through plan or apply stages
					if the document in the database is different from the document
					in the Terraform state.

					In contrast, if this option is false, the provider will
					ignore the consistency between the document in the Terraform
					state and the document in the database.

					This is useful when you want to manage the document whose counterpart
					in the database is managed by another system 
					(i.e. the document can be changed by other systems than Terraform) 
					but still want to perform CRUD operations on the document in the database with Terraform.

					It is IMPORTANT to note that if you once set this option either to true or false,
					you cannot change it back to the other value. 
					This is due to [terraform SDKv2's data consistency rules](https://developer.hashicorp.com/terraform/plugin/sdkv2/resources/data-consistency-errors), 
					keeping the resource state immutable once you set the value from the terraform side, 
					and it is impossible to modify the value from the provider side 
					if there are differences between the state and the database document.

					This value is true by default.
				`),
				PlanModifiers: []planmodifier.Bool{
					boolplanmodifier.UseStateForUnknown(),
				},
			},
		},
	}
}

func (r *DocumentResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	config, diags := resourceconfig.FromProviderData(req.ProviderData)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	r.config = config
}

func (r *DocumentResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	client := mongoclient.New(ctx, r.config.ClientConfig).WithLogger(r.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.Append(
				errs.NewMongoClientError(err).ToDiagnostic(),
			)
			return
		}
		var data DocumentResourceModel

		// Read Terraform plan data into the model
		resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform create operation
		resp.Diagnostics.Append(resourceCreate(client, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	})
}

func (r *DocumentResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	client := mongoclient.New(ctx, r.config.ClientConfig).WithLogger(r.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.Append(
				errs.NewMongoClientError(err).ToDiagnostic(),
			)
			return
		}

		var data DocumentResourceModel

		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform read operation
		resp.Diagnostics.Append(resourceRead(client, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save updated data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	})
}

func (r *DocumentResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	client := mongoclient.New(ctx, r.config.ClientConfig).WithLogger(r.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.Append(
				errs.NewMongoClientError(err).ToDiagnostic(),
			)
			return
		}

		var data DocumentResourceModel

		// Read Terraform plan data into the model
		resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		var state DocumentResourceModel

		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// If the sync_with_database attribute differs
		// between the plan and the state, return an error
		if state.SyncWithDatabase.ValueBool() != data.SyncWithDatabase.ValueBool() {
			resp.Diagnostics.Append(
				errs.NewInvalidResourceConfiguration(
					"sync_with_database attribute cannot be changed once set",
				).ToDiagnostic(),
			)
			return
		}

		// Perform the update operation
		resp.Diagnostics.Append(resourceUpdate(client, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Save updated data into Terraform state
		resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
	})
}

func (r *DocumentResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	client := mongoclient.New(ctx, r.config.ClientConfig).WithLogger(r.config.Logger)
	client.Run(func(client *mongoclient.MongoClient, err error) {
		if err != nil {
			resp.Diagnostics.Append(
				errs.NewMongoClientError(err).ToDiagnostic(),
			)
			return
		}

		var data DocumentResourceModel

		// Read Terraform prior state data into the model
		resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
		if resp.Diagnostics.HasError() {
			return
		}

		// Perform the delete operation
		resp.Diagnostics.Append(resourceDelete(client, &data)...)
	})
}

func (r *DocumentResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	id, err := resourceid.New(req.ID)
	if err != nil {
		resp.Diagnostics.Append(
			errs.NewInvalidImportID(err.Error()).ToDiagnostic(),
		)
		return
	}
	if id.Database() == "" {
		resp.Diagnostics.Append(
			errs.NewInvalidImportID("Database name is required").ToDiagnostic(),
		)
		return
	}
	if id.Collection() == "" {
		resp.Diagnostics.Append(
			errs.NewInvalidImportID("Collection name is required").ToDiagnostic(),
		)
		return
	}
	if id.Document() == "" {
		resp.Diagnostics.Append(
			errs.NewInvalidImportID("Document ID is required").ToDiagnostic(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("database"), id.Database())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("collection"), id.Collection())...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("document_id"), id.Document())...)
}
