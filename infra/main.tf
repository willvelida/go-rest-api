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
    azapi = {
      source = "Azure/azapi"
      version = "0.4.0"
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
}

# log analytics workspace
resource "azurerm_log_analytics_workspace" "law" {
  name = "law${random_id.random_deployment_suffix.hex}"
  location = azurerm_resource_group.rg.location
  resource_group_name = azurerm_resource_group.rg.name
  sku = "PerGB2018"
  retention_in_days = 30
}

# Container App Environment
resource "azapi_resource" "env" {
  type = "Microsoft.App/managedEnvironments@2022-06-01-preview"
  name = "env${random_id.random_deployment_suffix.hex}"
  location = azurerm_resource_group.rg.location
  parent_id = azurerm_resource_group.rg.id
  body = jsonencode({
    properties = {
      appLogsConfiguration = {
        destination = "log-analytics"
        logAnalyticsConfiguration = {
          customerId = azurerm_log_analytics_workspace.law.workspace_id
          sharedKey = azurerm_log_analytics_workspace.law.primary_shared_key
        }
      }
    }
  })
}

resource "azapi_resource" "containerapp" {
 type = "Microsoft.App/containerApps@2022-06-01-preview"
 name = "book-api"
 location = azurerm_resource_group.rg.location
 parent_id = azapi_resource.env.id
 body = jsonencode({
  properties = {
    configuration = {
      activeRevisionsMode = "Multiple"
      ingress = {
        allowInsecure = true
        exposedPort = 80
        external = true
        transport = "http"
      }
      registries = [
        {
          server = azurerm_container_registry.acr.login_server
          username = azurerm_container_registry.acr.admin_username
          passwordSecretRef = "containerregistry-password"
        }
      ]
      secrets = [
        {
          name = "containerregistry-password"
          value = azurerm_container_registry.acr.admin_password
        }
      ]
    }
    template = {
      containers = [
        {
          image = "mcr.microsoft.com/azuredocs/containerapps-helloworld:latest"
          name = "book-api"
          resources = {
            cpu = 0.5
            memory = "1.0Gi"
          }
        }
      ]
      scale = {
        minReplicas = 1
        maxReplicas = 1
      }
    }
  }
 })
}