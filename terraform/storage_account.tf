resource "azurerm_storage_account" "storage_account" {
  name                     = replace("sa${var.solution_name}${var.environment}", "-", "")
  resource_group_name      = azurerm_resource_group.resource_group.name
  location                 = azurerm_resource_group.resource_group.location
  account_tier             = var.storage_account_tier
  account_replication_type = var.storage_account_replication_type
  tags                     = var.tags

  network_rules {
    default_action = "Deny"
    ip_rules     = ["212.149.156.146"]
  }


}

resource "azurerm_storage_table" "storage_table" {
  name                 = var.storage_table_name
  storage_account_name = azurerm_storage_account.storage_account.name
  depends_on           = [azurerm_storage_account.storage_account]
}

resource "azurerm_storage_table_entity" "storage_table_entity" {
  storage_table_id = azurerm_storage_table.storage_table.id

  partition_key = 0
  row_key       = 0

  entity = {
    AppRegistrationName     = "Example Name"
    SecretName              = "If not specified, default will be used"
    AppRegistrationObjectId = "OBJECT ID!!!"
    GithubRepository        = "Optional"
    KeyVaultName            = "Required"
    NotificationEmail       = "Required"
    Owner                   = "Optional"
    Reprocess               = "automatically changed value"
    SubscriptionID          = "Required"
  }
  depends_on = [azurerm_storage_table.storage_table, azurerm_role_assignment.storage_account_role_assignment]
}

resource "azurerm_role_assignment" "storage_account_role_assignment" {
  scope                = azurerm_storage_account.storage_account.id
  role_definition_name = "Storage Table Data Contributor"
  principal_id         = azurerm_container_group.container_instance.identity[0].principal_id 
  depends_on           = [azurerm_storage_account.storage_account]
}
