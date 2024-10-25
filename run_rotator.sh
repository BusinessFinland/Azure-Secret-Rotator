#!/bin/sh

# Source environment variables
if [ -f /etc/environment ]; then
  . /etc/environment
fi

export "KEYVAULT_NAME=kvazsecrotator"


# Run the application
/azure-secret-rotator >> /var/log/azure-secret-rotator.log 2>&1
