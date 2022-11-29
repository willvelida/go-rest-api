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

resource "random_id" "random_deployment_suffix" {
  byte_length = 4
}

# Create a resource group
resource "azurerm_resource_group" "rg" {
  name     = "rg-go-rest-api"
  location = "australiaeast"
}

# azure container registry
resource "azurerm_container_registry" "acr" {
  name                = "acr${random_id.random_deployment_suffix.hex}"
  resource_group_name = azurerm_resource_group.rg.name
  location            = azurerm_resource_group.rg.location
  sku = "Basic"
  admin_enabled = true
  identity {
    type = "SystemAssigned"
  }
}

# Create the Linux App Service Plan
resource "azurerm_service_plan" "asp" {
  name                = "asp${random_id.random_deployment_suffix.hex}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  os_type             = "Linux"
  sku_name            = "B1"
}

resource "azurerm_linux_web_app" "webapp" {
  name                = "app${random_id.random_deployment_suffix.hex}"
  location            = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  service_plan_id     = azurerm_service_plan.asp.id
  https_only          = true
  site_config {
    always_on = true
  }
  identity {
    type = "SystemAssigned"
  }
} 