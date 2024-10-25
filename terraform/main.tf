resource "azurerm_resource_group" "resource_group" {
  name     = "rg-${var.solution_name}-${var.environment}"
  location = var.location
  tags     = var.tags
}

resource "azurerm_application_insights" "app_insights" {
  name                = "ai-${var.solution_name}-${var.environment}"
  location            = azurerm_resource_group.resource_group.location
  resource_group_name = azurerm_resource_group.resource_group.name
  application_type    = "web"
  tags                = var.tags

  depends_on = [azurerm_resource_group.resource_group]
}


resource "azurerm_container_group" "container_instance" {
  name                = "ci-${var.solution_name}-${var.environment}"
  location            = azurerm_resource_group.resource_group.location
  resource_group_name = azurerm_resource_group.resource_group.name
  os_type             = "Linux"
  tags                = var.tags

  ip_address_type = "None"

  identity {
    type = "SystemAssigned"
  }

  container {
    name   = "azure-secret-rotator"
    image  = "containerjrdev.azurecr.io/samples/azure-secret-rotator:1.0.0"
    cpu    = "0.25"
    memory = "0.5"

    # Secure environment variable for the container using secure value from Azure Key Vault
    secure_environment_variables = {
      "DOCKER_REGISTRY_SERVER_PASSWORD" = var.DOCKER_REGISTRY_SERVER_PASSWORD
    }
  }

  image_registry_credential {
    server   = var.DOCKER_REGISTRY_SERVER_URL
    username = var.DOCKER_REGISTRY_SERVER_USERNAME
    password = var.DOCKER_REGISTRY_SERVER_PASSWORD
  }

  # If you need to collect logs you can use log analytics
  

  depends_on = [azurerm_resource_group.resource_group, azurerm_subnet.subnet1]
}


