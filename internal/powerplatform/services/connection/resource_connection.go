// Copyright (c) Microsoft Corporation.
// Licensed under the MIT license.

package powerplatform

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	api "github.com/microsoft/terraform-provider-power-platform/internal/powerplatform/api"
)

var _ resource.Resource = &ConnectionResource{}
var _ resource.ResourceWithImportState = &ConnectionResource{}

func NewConnectionResource() resource.Resource {
	return &ConnectionResource{
		ProviderTypeName: "powerplatform",
		TypeName:         "_connection",
	}
}

type ConnectionResource struct {
	ConnectionsClient ConnectionsClient
	ProviderTypeName  string
	TypeName          string
}

type ConnectionResourceModel struct {
	Id            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	EnvironmentId types.String `tfsdk:"environment_id"`
	DisplayName   types.String `tfsdk:"display_name"`
	Status        types.List   `tfsdk:"status"`
	//Parameters    types.String `tfsdk:"parameters"`
}

func (r *ConnectionResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + r.TypeName
}

func (r *ConnectionResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "Manages a \"Connection\"",
		MarkdownDescription: "Manages a [Connection](https://learn.microsoft.com/en-us/power-apps/maker/canvas-apps/add-manage-connections). A connection in Power Platform serves as a means to integrate external data sources and services with your Power Platform apps, flows, and other solutions. It acts as a bridge, facilitating secure communication between your solutions and various external systems.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "Unique connection id",
				Description:         "Unique connection id",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"environment_id": schema.StringAttribute{
				MarkdownDescription: "Environment id where the connection is to be created",
				Description:         "Environment id where the connection is to be created",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "Name of the connection. This can be found using `powerplatform_connectors` data source by using the `name` attribute",
				Description:         "Name of the connection. This can be found using `powerplatform_connectors` data source by using the `name` attribute",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"display_name": schema.StringAttribute{
				MarkdownDescription: "Display name of the connection",
				Description:         "Display name of the connection",
				Optional:            true,
				Computed:            true,
			},
			"status": schema.ListAttribute{
				Description:         "List of connection statuses",
				MarkdownDescription: "List of connection statuses",
				ElementType:         types.StringType,
				Computed:            true,
			},
			// "parameters": schema.StringAttribute{
			// 	Description:         "Connection parameters. Json string containing the authentication connection parameters. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
			// 	MarkdownDescription: "Connection parameters. Json string containing the authentication connection parameters, (for example)[https://learn.microsoft.com/en-us/power-automate/desktop-flows/alm/alm-connection#create-a-connection-using-your-service-principal]. Depending on required authentication parameters of a given connector, the connection parameters can vary.",
			// 	Optional:            true,
			// 	Computed:            true,
			// },
		},
	}
}

func (r *ConnectionResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	clientApi := req.ProviderData.(*api.ProviderClient).Api

	if clientApi == nil {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *http.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)

		return
	}
	r.ConnectionsClient = NewConnectionsClient(clientApi)
}

func (r *ConnectionResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan *ConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// var params map[string]interface{} = nil
	// if !plan.Parameters.IsNull() && plan.Parameters.ValueString() != "" {
	// 	err := json.Unmarshal([]byte(plan.Parameters.ValueString()), &params)
	// 	if err != nil {
	// 		resp.Diagnostics.AddError("Failed to convert connection parameters", err.Error())
	// 		return
	// 	}
	// }

	connectionToCreate := ConnectionToCreateDto{
		Properties: ConnectionToCreatePropertiesDto{
			DisplayName: plan.DisplayName.ValueString(),
			Environment: ConnectionToCreateEnvironmentDto{
				Name: plan.EnvironmentId.String(),
				Id:   fmt.Sprintf("/providers/Microsoft.PowerApps/environments/%s", plan.EnvironmentId.ValueString()),
			},
			//ConnectionParameters: params,
		},
	}

	connection, err := r.ConnectionsClient.CreateConnection(ctx, plan.EnvironmentId.ValueString(), plan.Name.ValueString(), connectionToCreate)
	if err != nil {
		resp.Diagnostics.AddError("Failed to create connection", err.Error())
		return
	}

	conectionState := ConvertFromConnectionDto(*connection)
	plan.Id = types.String(conectionState.Id)
	plan.DisplayName = types.String(conectionState.DisplayName)
	plan.Status = conectionState.Status
	//plan.Parameters = types.String(conectionState.Parameters)
	plan.Name = types.String(conectionState.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("CREATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ConnectionResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state *ConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE START: %s", r.TypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connection, err := r.ConnectionsClient.GetConnection(ctx, state.EnvironmentId.ValueString(), state.Name.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when reading %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	conectionState := ConvertFromConnectionDto(*connection)
	state.Id = types.String(conectionState.Id)
	state.DisplayName = types.String(conectionState.DisplayName)
	state.Status = conectionState.Status
	//state.Parameters = types.String(conectionState.Parameters)
	state.Name = types.String(conectionState.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
	tflog.Debug(ctx, fmt.Sprintf("READ RESOURCE END: %s", r.TypeName))
}

func (r *ConnectionResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan *ConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	var state *ConnectionResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	connection, err := r.ConnectionsClient.UpdateConnection(ctx, plan.EnvironmentId.ValueString(), plan.Name.ValueString(), plan.Id.ValueString(), plan.DisplayName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when updating %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	conectionState := ConvertFromConnectionDto(*connection)
	plan.Id = types.String(conectionState.Id)
	plan.DisplayName = types.String(conectionState.DisplayName)
	plan.Status = conectionState.Status
	//plan.Parameters = types.String(conectionState.Parameters)
	plan.Name = types.String(conectionState.Name)

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)

	tflog.Debug(ctx, fmt.Sprintf("UPDATE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ConnectionResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state *ConnectionResourceModel

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE START: %s", r.ProviderTypeName))

	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)

	if resp.Diagnostics.HasError() {
		return
	}

	err := r.ConnectionsClient.DeleteConnection(ctx, state.EnvironmentId.ValueString(), state.Name.ValueString(), state.Id.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(fmt.Sprintf("Client error when deleting %s_%s", r.ProviderTypeName, r.TypeName), err.Error())
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("DELETE RESOURCE END: %s", r.ProviderTypeName))
}

func (r *ConnectionResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
