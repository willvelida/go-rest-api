terraform {
  backend "azurerm" {
    resource_group_name = "willvtfstates"
    storage_account_name = "willvstf"
    container_name = "tfstate"
    key = "gowebapp.tfstate"
  }
  required_providers {
    azurerm = {
      source  = "hashicorp/azurerm"
      version = ">= 3.0.0"
    }
  }
  required_version = ">= 0.13"
}

provider "azurerm" {
  features {}
}

data "azurerm_client_config" "current" {}

# Create a resource group
resource "azurerm_resource_group" "rg" {
  name     = var.rg_name
  location = var.location
}

# Create the Linux App Service Plan
resource "azurerm_service_plan" "asp" {
  name                = var.asp_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  os_type             = "Linux"
  sku_name            = "B1"
}

resource "azurerm_linux_web_app" "webapp" {
  name                = var.app_name
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  service_plan_id     = azurerm_service_plan.asp.id
  https_only          = true
  site_config {
    always_on = true
  }
}