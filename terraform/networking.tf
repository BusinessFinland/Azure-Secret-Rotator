resource "azurerm_virtual_network" "vnet" {
  name                = "vnet-${var.solution_name}-${var.environment}"
  address_space       = var.vnet_address_space
  location            = var.location
  resource_group_name = azurerm_resource_group.resource_group.name
}

resource "azurerm_subnet" "subnet1" {
  name                 = "snet-${var.solution_name}-${var.environment}-1"
  resource_group_name  = azurerm_resource_group.resource_group.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = var.subnet_address_space1

  service_endpoints = ["Microsoft.Web"]
  # Ensures service delegation to Microsoft.Web for App Service
  delegation {
    name = "delegation"
    service_delegation {
      name    = "Microsoft.Web/serverFarms"
      actions = ["Microsoft.Network/virtualNetworks/subnets/action"]
    }
  }
}

resource "azurerm_subnet" "subnet2" {
  name                 = "snet-${var.solution_name}-${var.environment}-2"
  resource_group_name  = azurerm_resource_group.resource_group.name
  virtual_network_name = azurerm_virtual_network.vnet.name
  address_prefixes     = var.subnet_address_space2

  service_endpoints = ["Microsoft.Web"]
 
}


