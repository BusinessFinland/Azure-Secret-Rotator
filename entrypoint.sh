#!/bin/sh
# Generate /etc/environment from current environment variables
printenv | grep -Ev '^(HOME|PWD|PATH|SHELL|USER|LOGNAME|LANG|LANGUAGE)' | awk -F= '{print $1 "=" "\"" $2 "\""}' > /etc/environment

chmod 0644 /etc/environment

exec crond -f
