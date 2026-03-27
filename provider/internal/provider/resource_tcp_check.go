package provider

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64default"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ resource.Resource = &TCPCheckResource{}

func NewTCPCheckResource() resource.Resource {
	return &TCPCheckResource{}
}

type TCPCheckResource struct{}

type TCPCheckResourceModel struct {
	ID       types.String `tfsdk:"id"`
	Host     types.String `tfsdk:"host"`
	Port     types.Int64  `tfsdk:"port"`
	Interval types.Int64  `tfsdk:"interval"`
	Timeout  types.Int64  `tfsdk:"timeout"`
}

func (r *TCPCheckResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_tcp_check"
}

func (r *TCPCheckResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Performs a TCP health check against a host and port with configurable interval and timeout.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:    true,
				Description: "Unique identifier for this health check",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"host": schema.StringAttribute{
				Required:    true,
				Description: "The hostname or IP address to check",
			},
			"port": schema.Int64Attribute{
				Required:    true,
				Description: "The TCP port to check",
			},
			"interval": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(5),
				Description: "The interval in seconds between health check attempts (default: 5)",
			},
			"timeout": schema.Int64Attribute{
				Optional:    true,
				Computed:    true,
				Default:     int64default.StaticInt64(60),
				Description: "The total timeout in seconds for all health check attempts (default: 60)",
			},
		},
	}
}

func (r *TCPCheckResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data TCPCheckResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform TCP health check
	host := data.Host.ValueString()
	port := data.Port.ValueInt64()
	interval := time.Duration(data.Interval.ValueInt64()) * time.Second
	timeout := time.Duration(data.Timeout.ValueInt64()) * time.Second

	tflog.Info(ctx, fmt.Sprintf("Starting TCP health check for %s:%d", host, port))
	tflog.Info(ctx, fmt.Sprintf("Interval: %v, Timeout: %v", interval, timeout))

	// Perform the health check with retries
	if err := performTCPCheck(ctx, host, port, interval, timeout); err != nil {
		resp.Diagnostics.AddError(
			"TCP Health Check Failed",
			fmt.Sprintf("Failed to connect to %s:%d within timeout: %s", host, port, err.Error()),
		)
		return
	}

	// Set ID
	data.ID = types.StringValue(fmt.Sprintf("%s:%d", host, port))

	tflog.Info(ctx, fmt.Sprintf("TCP health check succeeded for %s:%d", host, port))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TCPCheckResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data TCPCheckResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Nothing to do on read - state is already in sync
}

func (r *TCPCheckResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var data TCPCheckResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Perform TCP health check on update
	host := data.Host.ValueString()
	port := data.Port.ValueInt64()
	interval := time.Duration(data.Interval.ValueInt64()) * time.Second
	timeout := time.Duration(data.Timeout.ValueInt64()) * time.Second

	tflog.Info(ctx, fmt.Sprintf("Updating TCP health check for %s:%d", host, port))

	if err := performTCPCheck(ctx, host, port, interval, timeout); err != nil {
		resp.Diagnostics.AddError(
			"TCP Health Check Failed",
			fmt.Sprintf("Failed to connect to %s:%d within timeout: %s", host, port, err.Error()),
		)
		return
	}

	data.ID = types.StringValue(fmt.Sprintf("%s:%d", host, port))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *TCPCheckResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data TCPCheckResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Do nothing on destroy as per requirements
	tflog.Info(ctx, fmt.Sprintf("Delete called for %s:%d - doing nothing as per requirements",
		data.Host.ValueString(), data.Port.ValueInt64()))
}

// performTCPCheck performs TCP health check with retries
func performTCPCheck(ctx context.Context, host string, port int64, interval, timeout time.Duration) error {
	startTime := time.Now()
	endTime := startTime.Add(timeout)
	attempt := 0

	for {
		attempt++
		currentTime := time.Now()

		// Check if we've exceeded the timeout
		if currentTime.After(endTime) {
			return fmt.Errorf("timeout of %v exceeded after %d attempts", timeout, attempt)
		}

		tflog.Debug(ctx, fmt.Sprintf("Attempt %d: Checking %s:%d", attempt, host, port))

		// Attempt TCP connection with a per-attempt timeout
		address := fmt.Sprintf("%s:%d", host, port)
		conn, err := net.DialTimeout("tcp", address, 5*time.Second)

		if err == nil {
			conn.Close()
			elapsed := time.Since(startTime)
			tflog.Info(ctx, fmt.Sprintf("Success: %s:%d is reachable (took %v, %d attempts)", host, port, elapsed, attempt))
			return nil
		}

		tflog.Debug(ctx, fmt.Sprintf("Failed: %s:%d is not reachable - %v", host, port, err))

		// Check if we have time for another attempt
		currentTime = time.Now()
		if currentTime.Add(interval).After(endTime) {
			remaining := endTime.Sub(currentTime)
			return fmt.Errorf("insufficient time for next attempt (%v remaining)", remaining)
		}

		tflog.Debug(ctx, fmt.Sprintf("Waiting %v before next attempt", interval))

		// Sleep with context cancellation support
		select {
		case <-time.After(interval):
			continue
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
