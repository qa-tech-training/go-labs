provider "simple1" {}

resource "simple1_task" "example" {
  title       = "Test from Registry"
  description = "Created using the local provider"
  completed   = false
}