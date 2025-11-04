# Lab GO03 - Build a Custom Terraform Provider

## Objective
To build and use a custom terraform provider to interact with an API

## Outcomes
* Understand the structure of a terraform provider
* Build a provider using golang tooling
* Use a custom provider in a tf file

## High-Level Steps
* Review the codebase
* Build and relocate a binary
* Create a terraform configuration which references the new provider
* Create, update and delete objects using the custom provider

## Detailed Steps
### Examine the Provisioned Files
Change directory into the lab03 folder:
```bash
cd ~/go-labs/lab03
```
Focus the lab03 folder in the editor's file explorer. Review the following files:
#### File: main.go
Purpose: This is the entry point for the Terraform provider binary. When Terraform initializes your provider, it invokes this function, which starts the plugin server.
#### File: provider/provider.go
Purpose: Registers your provider with Terraform (TypeName), sets its version (Version) and declares available resources and data sources (we only declare NewTaskResource() now)
#### File: provider/task_resource.go
Purpose: This defines the resource that the provider exposes, in this case, a simple task resource interacting with an API
Note: There is a Deep-Dive section at the end of these lab instructions that details these files.
### Initialize Golang
To set up a go.mod and go.sum file for the project, run the following: 
```bash
go mod init mynewprovider
```
Then, run: 
```bash
go mod tidy
```
This will download any missing dependencies (like the Terraform Plugin Framework) and clean up unused imports and updates the go.sum file
### Build the Provider Binary
Now you're ready to compile your provider into an executable that Terraform can use. Run:
```bash
go build -o terraform-provider-mynewprovider -buildvcs=false
```
This tells Go: 
	go build: compile the code in this module
    -o terraform-provider-mynewprovider: name the output binary with Terraform’s expected naming convention
    -buildvcs=false: don't use git commit info to automatically version the binary
### Move the provider binary
Move the new binary into terraform's default search path:
```bash
mkdir -p ~/.terraform.d/plugins/registry.terraform.io/custom/mynewprovider/0.1.0/linux_amd64
cp terraform-provider-mynewprovider ~/.terraform.d/plugins/registry.terraform.io/custom/mynewprovider/0.1.0/linux_amd64/
```
### Test Terraform Configuration
In the editor, navigate to the testing folder. Open main.tf and populate with:
```terraform
terraform {
  required_providers {
    mynewprovider = {
      source  = "custom/mynewprovider"
      version = "0.1.0"
    }
  }
}

provider "mynewprovider" {}

resource "mynewprovider_task" "example" {
  title       = "From custom provider"
  description = "Created using mynewprovider"
  completed   = false
}
```
In the terminal, navigate to the testing directory and run:
```bash
terraform init
terraform plan
terraform apply
```
Init and Plan will succeed but apply will fail as there is no target (mock) api running
In a new terminal (keeping existing session open), navigate back to the lab02 directory and start the app from earlier:
```bash
go run main.go
```
Leave the app running. Back in the original terminal, re-run the terraform apply. It should now succeed. Verify the task creation:
```bash
curl http://localhost:8080/tasks
```
Update the Terraform configuration in main.tf, changing any of the values
Re-run terraform apply. To verify the update:
```bash
curl http://localhost:8080/tasks
```
The task attributes should have been successfully updated

## The Sample Provider - a Deeper Dive
Understanding Core Files in a Custom Terraform Provider
### main.go — Entry Point for the Terraform Plugin
Purpose:
This is the first file that runs when Terraform loads the provider binary. It registers the provider with Terraform by calling providerserver.Serve.  
Key Functionality:
```go
func main() {
  providerserver.Serve(
    context.Background(),
  provider.New,
  providerserver.ServeOpts{
    Address: "registry.terraform.io/custom/mynewprovider",
  },
  )
}
```
Explanation:
- provider.New: Returns an instance of your provider (defined in provider.go).
- ServeOpts{ Address: ... }: Identifies the provider namespace/type.
- providerserver.Serve: Boots the plugin server so Terraform can communicate with it.

### provider.go — Declares the Provider and Its Resources
Purpose:
This file defines what the provider exposes to Terraform — metadata, available resources, and configuration.  
Key Functionality:
```go
func (p *myNewProvider) Metadata(...) {
  resp.TypeName = "mynewprovider"
  resp.Version = "0.1.0"
}

func (p *myNewProvider) Resources(...) []func() resource.Resource {
  return []func() resource.Resource{
    NewTaskResource,
  }
}
```
Explanation:
- Metadata: Declares provider name and version.
- Resources: Registers resource constructors (e.g. task resource).
- Configure: (Optional) Injects configuration like API endpoints.
### task_resource.go — Implements the Custom Resource Logic
Purpose:
Defines the Create, Read, Update, Delete (CRUD) behavior of your custom resource.  
Key Structure:
```go
type taskResource struct{}
type taskModel struct {
  ID          types.String `tfsdk:"id"`
  Title        types.String `tfsdk:"title"`
  Description  types.String `tfsdk:"description"`
  Completed  types.Bool  `tfsdk:"completed"`
}

func (r *taskResource) Create(...) { ... }
func (r *taskResource) Read(...)   { ... }
func (r *taskResource) Update(...) { ... }
func (r *taskResource) Delete(...) { ... }
```
Explanation:
- taskModel: Maps to Terraform config.
- Create: Sends POST to API.
- Read: Sends GET to API.
- Update: Sends PUT to API.
- Delete: Sends DELETE to API.

How They Work Together
| File | Role |	Who Calls It | Interacts With |
| :--- | :--- | :--- | :--- |
| main.go | Starts the provider plugin | Terraform CLI | Calls provider.New |
| provider.go | Registers provider & resources | main.go | Returns taskResource |
| task_resource.go | Implements mynewprovider_task | provider.go |Talks to API |

Terraform executes the binary → runs main.go → which returns provider.go → which loads task_resource.go.