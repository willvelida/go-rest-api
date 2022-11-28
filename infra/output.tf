output "resource_group_name" {
  value = azurerm_resource_group.rg.name
}

output "app_service_plan_name" {
  value = azurerm_service_plan.asp.name
}

output "app_service_plane_os_type" {
    value = azurerm_service_plan.asp.os_type
}

output "linux_web_app_name" {
  value = azurerm_linux_web_app.webapp.name 
}

output "linux_web_app_asp" {
  value = azurerm_linux_web_app.webapp.service_plan_id
}