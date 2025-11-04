# Lab GO04 - Create a Terraform Provider Using the Scaffolding Template

## Objective
Use the Terraform provider scaffolding template to streamline provider creation

## Outcomes
* Create a Terraform provider using a template as a starting point

## High-Level Steps
* Create a new repo based on the scaffolding template
* Edit the files to implement custom provider logic
* Build and test the provider

## Detailed Steps
### Create a New Repository
If you do not have a github account, you will need to create one to complete this lab. Log into your github account, then find the scaffolding template repository: https://github.com/hashicorp/terraform-provider-scaffolding-framework  
Click the green "Use this template" button (top right).
Select "Create a new repository".
Enter your github user and password if prompted.
Name your new repo _terraform-provider-simple1_.
Click "Create repository from template."  
In your lab environment terminal, run:
```bash
git clone https://github.com/<your-github-username>/terraform-provider-simple1.git ~/go-labs/lab04
cd ~/go-labs/lab04
```
### Edit the Template
Focus the new lab04 folder in the editor, then, in go.mod, replace:
```
module github.com/hashicorp/terraform-provider-scaffolding-framework
```
with 
```
module github.com/<your-github-username>/terraform-provider-simple1
```
Making sure to fill in your github username in the module path
Open main.go and edit the contents like so:
```go
package main

import (
	"context"
	"flag"
	"log"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/<your_github_username>/terraform-provider-simple1/simple1"
)

var version = "dev"

func main() {
	var debug bool
	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	err := providerserver.Serve(context.Background(), simple1.New(version), providerserver.ServeOpts{
		Address: "registry.terraform.io/custom/simple1",
		Debug:   debug,
	})

	if err != nil {
		log.Fatal(err.Error())
	}
}
```
Making sure again to substitute your github username and repo information where indicated
Create a new folder, “simple1”, in the repository. This will house your clean provider logic.  
copy provider.go from internal/provider to simple1:
```bash
cp internal/provider/provider.go simple1/
```
update provider.go in the simple1 directory with the following content:
```go
// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package simple1

import (
	"context"
	"net/http"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/function"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure ScaffoldingProvider satisfies various provider interfaces.
var _ provider.Provider = &ScaffoldingProvider{}
var _ provider.ProviderWithFunctions = &ScaffoldingProvider{}
var _ provider.ProviderWithEphemeralResources = &ScaffoldingProvider{}

// ScaffoldingProvider defines the provider implementation.
type ScaffoldingProvider struct {
	// version is set to the provider version on release, "dev" when the
	// provider is built and ran locally, and "test" when running acceptance
	// testing.
	version string
}

// ScaffoldingProviderModel describes the provider data model.
type ScaffoldingProviderModel struct {
	Endpoint types.String `tfsdk:"endpoint"`
}

func (p *ScaffoldingProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "scaffolding"
	resp.Version = p.version
}

func (p *ScaffoldingProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"endpoint": schema.StringAttribute{
				MarkdownDescription: "Example provider attribute",
				Optional:            true,
			},
		},
	}
}

func (p *ScaffoldingProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	var data ScaffoldingProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Configuration values are now available.
	// if data.Endpoint.IsNull() { /* ... */ }

	// Example client configuration for data sources and resources
	client := http.DefaultClient
	resp.DataSourceData = client
	resp.ResourceData = client
}

func (p *ScaffoldingProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewTaskResource,
	}
}

func (p *ScaffoldingProvider) EphemeralResources(ctx context.Context) []func() ephemeral.EphemeralResource {
	return []func() ephemeral.EphemeralResource{}
}

func (p *ScaffoldingProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{}
}

func (p *ScaffoldingProvider) Functions(ctx context.Context) []func() function.Function {
	return []func() function.Function{}
}

func New(version string) func() provider.Provider {
	return func() provider.Provider {
		return &ScaffoldingProvider{
			version: version,
		}
	}
}
```
Summary of changes:
* package name: provider -> simple1
* Return value for resource function: { NewExampleResource } -> { NewTaskResource }
* Return values for EphemeralResource, DataSource and Function functions: { original value } -> {}
### Create TaskResource Type
Add another new file, task_resource.go, in the simple1 directory
Browse to https://github.com/qatip/provider-lab-files.git and copy the contents of task_resource.go. Add them to the file you just created.
### Build the Provider
Run:
```bash
cd ~/go-labs/lab04 # ensure you're in the right place
go mod tidy
go build -o terraform-provider-simple1 -buildvcs=false
```
### Test the Provider
Create a new directory, provider_test:
```bash
mkdir ~/provider_test
```
Copy the binary into that new directory:
```bash
cp terraform-provider-simple1 ~/provider_test/
cd ~/provider_test
```
Focus the new directory in the editor. Create a new _terraform.rc_ file:
```bash
cat - >> terraform.rc <<EOF
provider_installation {
    dev_overrides {
    "simple1" = "$(pwd)"
  }
  direct {
    exclude = ["registry.terraform.io/*/*"]
  }
}
EOF
```
Then create a simple main.tf config:
```terraform
provider "simple1" {}

resource "simple1_task" "example" {
  title       = "Test from Registry"
  description = "Created using the local provider"
  completed   = false
}
```
To test the configuration, start the web app from lab02 in another terminal session, then run:
```bash
export TF_CLI_CONFIG_FILE="$(pwd)/terraform.rc"
terraform plan
terraform apply
```
Verify the creation of the resource:
```bash
curl http://localhost:8080/tasks
```
