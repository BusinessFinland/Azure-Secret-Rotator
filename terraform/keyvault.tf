resource "azurerm_key_vault" "key_vault" {
  name                      = var.key_vault_name
  location                  = azurerm_resource_group.resource_group.location
  resource_group_name       = azurerm_resource_group.resource_group.name
  sku_name                  = var.key_vault_sku
  tenant_id                 = var.tenant_id
  enable_rbac_authorization = true
  tags                      = var.tags

  network_acls {
    bypass         = "AzureServices"
    default_action = "Allow"
    ip_rules       = ["87.94.154.178", "212.149.156.146"] 
  }

  depends_on = [azurerm_container_group.container_instance, azurerm_application_insights.app_insights]
}

resource "azurerm_role_assignment" "key_vault_role_assignment" {
  scope                = azurerm_key_vault.key_vault.id
  role_definition_name = "Key Vault Secrets User"
  principal_id         = azurerm_container_group.container_instance.identity[0].principal_id 

  depends_on = [azurerm_container_group.container_instance, azurerm_key_vault.key_vault]
}

resource "azurerm_role_assignment" "key_vault_role_assignment2" {
  scope                = azurerm_key_vault.key_vault.id
  role_definition_name = "Key Vault Secrets Officer"
  principal_id         = var.principal_id 

  depends_on = [azurerm_key_vault.key_vault]
}

data "azurerm_key_vault_secret" "existing" {
  for_each     = var.secret_names
  key_vault_id = azurerm_key_vault.key_vault.id
  name         = each.value
}

resource "azurerm_key_vault_secret" "key_vault_secrets" {
  for_each     = var.secret_names
  key_vault_id = azurerm_key_vault.key_vault.id
  name         = each.value

  # The conditional expression checks if the existing secret is not found (null) or has the value "To Be Replaced".
  # In both those cases, it will use the value from the variable `secret_values`. Otherwise, it retains the current value.
  value = coalesce(
    try(
      data.azurerm_key_vault_secret.existing[each.value].value != "To Be Replaced" ? 
        data.azurerm_key_vault_secret.existing[each.value].value : 
        null,
      null
    ),
  )
  
  depends_on = [azurerm_role_assignment.key_vault_role_assignment2]
}

resource "azurerm_key_vault_secret" "infrastructure_outputs" {
  key_vault_id = azurerm_key_vault.key_vault.id
  name         = "APPINSIGHTS-INSTRUMENTATIONKEY-CONNECTIONSTRING"

  value = azurerm_application_insights.app_insights.connection_string 
  depends_on = [azurerm_role_assignment.key_vault_role_assignment2]
}


  ## resource "azurerm_private_endpoint" "private_endpoint" {
  ##   name                = "pe-${var.solution_name}-${var.environment}"
  ##   location            = azurerm_resource_group.resource_group.location
  ##   resource_group_name = azurerm_resource_group.resource_group.name
  ##   subnet_id           = azurerm_subnet.subnet2.id
  ## 
  ##   private_service_connection {
  ##     name                           = "psc-${var.solution_name}-${var.environment}"
  ##     private_connection_resource_id = azurerm_key_vault.key_vault.id 
  ##     is_manual_connection           = false
  ##     subresource_names              = ["vault"]
  ##   }
  ## }
  ## 
  ## # Create a Private DNS Zone for Azure Key Vault
  ## resource "azurerm_private_dns_zone" "private_dns_zone" {
  ##   name                = "privatelink.vaultcore.azure.net"
  ##   resource_group_name = azurerm_resource_group.resource_group.name
  ## }
  ## 
  ## # Link the Private DNS Zone to the VNet
  ## resource "azurerm_private_dns_zone_virtual_network_link" "dns_zone_vnet_link" {
  ##   name                  = "pdzvn-link-${var.solution_name}-${var.environment}"
  ##   resource_group_name   = azurerm_resource_group.resource_group.name
##   private_dns_zone_name = azurerm_private_dns_zone.private_dns_zone.name
##   virtual_network_id    = azurerm_virtual_network.vnet.id
##   registration_enabled  = false
## }
## 
## # Configure the Private Endpoint's DNS record within the Private DNS Zone
## resource "azurerm_private_dns_a_record" "private_dns_a_record" {
##   name                = split(".", azurerm_key_vault.key_vault.name)[0]
##   zone_name           = azurerm_private_dns_zone.private_dns_zone.name
##   resource_group_name = azurerm_private_dns_zone.private_dns_zone.resource_group_name
##   ttl                 = 300
##   records             = [azurerm_private_endpoint.private_endpoint.private_service_connection[0].private_ip_address]
## }
