variable "solution_name" {
  type    = string
  default = "azure-secret-rotator"
}

variable "environment" {
  type    = string
  default = "dev"
}

variable "vnet_address_space" {
  default = ["10.0.0.0/22"]
}

variable "subnet_address_space1" {
  default = ["10.0.0.0/24"]
}

variable "subnet_address_space2" {
  default = ["10.0.2.0/23"]
}

variable "location" {
  type    = string
  default = "West Europe"
}

variable "tags" {
  type = map(string)
  default = {
  }
}

variable "storage_account_tier" {
  type    = string
  default = "Standard"
}

variable "storage_account_replication_type" {
  type    = string
  default = "LRS"
}

variable "storage_table_name" {
  type    = string
  default = "appRegistrations"
}

variable "key_vault_sku" {
  type    = string
  default = "standard"
}

variable "key_vault_name" {
  type    = string
  default = "kvazsecrotator"
}

variable "tenant_id" {
  type    = string
  default = ""
}

variable "subscription_id" {
  type    = string
  default = ""
}

variable "secret_names" {
  type    = set(string)
  default = [ 
              "TENANT-ID",
              "AZURE-STORAGE-ACCOUNT",
              "AZURE-TABLE-NAME",
              "GITHUB-APP-ID",
              "GITHUB-INSTALLATION-ID",
              "ORGANIZATION-NAME",
              "SENDGRID-API-KEY"
            ]
}

variable "DOCKER_REGISTRY_SERVER_USERNAME" {
  type = string
  default = ""
}

variable "DOCKER_REGISTRY_SERVER_PASSWORD" {
  type = string
  default = ""
}

variable "DOCKER_REGISTRY_SERVER_URL" {
  type = string
  default = ""
}

variable "principal_id" {
  type = string
  default = ""
}
