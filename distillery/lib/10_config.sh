#!/bin/bash
set -e

# This is a library file. 
# It should be 'source'd only, if it is not we bail out here. 
if [[ "$0" = "$BASH_SOURCE" ]]; then
   echo "This file should not be executed directly, it should be 'source'd only. "
   exit 1;
fi

# The Path to the configuration file. 
CONFIG_FILE="$SCRIPT_DIR/.env"

# Check that the configuration file exists. 
# If it does not, throw an error
log_info " => Reading configuration file"
if ! [ -f "$CONFIG_FILE" ]; then
   log_error ""
   log_error "Missing configuration, provide a '.env' file"
    exit 1
fi

# 'source' in the configuration file. 
# Ideally we would want to make sure to prevent code-executation within the .env file
# But for the moment let's not. 
source "$CONFIG_FILE"

# Next, validate all the configuration settings. 

# is_valid_slug checks if it's argument is a valid 'slug'. 
# A slug is any non-empty string of alphanumeric characters or '-'s.
# The first character of a slug may not be a dash. 
function is_valid_slug() {
   if [[ "$1" =~ ^[a-zA-Z0-9][-a-zA-Z0-9]*$ ]]; then
      return 0;
   else
      return 1;
   fi;
}

# is_valid_abspath checks if it's argument is an absolute path. 
function is_valid_abspath() {
   if [[ "$1" =~ ^\/(.+)\/([^/]+)$ ]]; then
      return 0;
   else
      return 1;
   fi;
}

# 'is_valid_domain' checks if a number is a valid domain. 
# A domain consists of at least one slug, seperated by '.'s.
# Each token is a slug. 
function is_valid_domain() {
   if [[ "$1" =~ ^([a-zA-Z0-9][-a-zA-Z0-9]*\.)*[a-zA-Z0-9][-a-zA-Z0-9]*$ ]]; then
      return 0;
   else
      return 1;
   fi;
}

# 'is_valid_number' checks if a value is a valid number. 
function is_valid_number() {
   if [[ "$1" =~ ^[1-9][0-9]*$ ]]; then
      return 0;
   else
      return 1;
   fi;
}

# 'is_valid_email' checks if a value is a valid email address
function is_valid_email() {
   if [[ "$1" =~ ^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$ ]]; then
      return 0;
   else
      return 1;
   fi
}

# The 'DEPLOY_ROOT' variable must be an absolute path. 
if ! is_valid_abspath "$DEPLOY_ROOT"; then
   log_error "Variable 'DEPLOY_ROOT' is missing or not a valid path. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   log_info "Please ensure that it does not end in '/'. ";
   exit 1;
fi

# The 'MYSQL_USER_PREFIX' variable must be a valid slug. 
if ! is_valid_slug "$MYSQL_USER_PREFIX"; then
   log_error "Variable 'MYSQL_USER_PREFIX' is missing or not a valid slug. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   exit 1;
fi

# The 'MYSQL_DATABASE_PREFIX' variable must be a valid slug. 
if ! is_valid_slug "$MYSQL_DATABASE_PREFIX"; then
   log_error "Variable 'MYSQL_DATABASE_PREFIX' is missing or not a valid slug. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   exit 1;
fi

# The 'DISTILLERY_BOOKKEEPING_DATABASE' variable must be a valid slug. 
if ! is_valid_slug "$DISTILLERY_BOOKKEEPING_DATABASE"; then
   log_error "Variable 'DISTILLERY_BOOKKEEPING_DATABASE' is missing or not a valid slug. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   exit 1;
fi

# The 'DISTILLERY_BOOKKEEPING_TABLE' variable must be a valid slug. 
if ! is_valid_slug "$DISTILLERY_BOOKKEEPING_TABLE"; then
   log_error "Variable 'DISTILLERY_BOOKKEEPING_TABLE' is missing or not a valid slug. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   exit 1;
fi


# The 'GRAPHDB_USER_PREFIX' variable must be a valid slug. 
if ! is_valid_slug "$GRAPHDB_USER_PREFIX"; then
   log_error "Variable 'DATABASE_PREFIX' is missing or not a valid slug. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   exit 1;
fi

# The 'GRAPHDB_REPO_PREFIX' variable must be a valid slug. 
if ! is_valid_slug "$GRAPHDB_REPO_PREFIX"; then
   log_error "Variable 'GRAPHDB_REPO_PREFIX' is missing or not a valid slug. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   exit 1;
fi

# The 'DEFAULT_DOMAIN' variable must be a valid domain. 
if ! is_valid_domain "$DEFAULT_DOMAIN"; then
   log_error "Variable 'DEFAULT_DOMAIN' is missing or not a valid domain. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   exit 1;
fi

# The 'CERTBOT_EMAIL' variable should either be empty or a valid email
if [ -n "$CERTBOT_EMAIL" ]; then
   if ! is_valid_email "$CERTBOT_EMAIL"; then
         log_error "Variable 'CERTBOT_EMAIL' is not a valid email address. ";
         log_info "Please verify that it is set correctly in '.env' or remove it completly. ";
         exit 1;
   fi;
fi

# The 'PASSWORD_LENGTH' variable must be a valid number. 
if ! is_valid_number "$PASSWORD_LENGTH"; then
   log_error "Variable 'PASSWORD_LENGTH' is missing or not a valid number. ";
   log_info "Please verify that it is set correctly in '.env'. ";
   exit 1;
fi

# paths to composer things
DEPLOY_WEB_DIR="$DEPLOY_ROOT/core/web"
DEPLOY_TRIPLESTORE_DIR="$DEPLOY_ROOT/core/triplestore"
DEPLOY_SQL_DIR="$DEPLOY_ROOT/core/sql"
DEPLOY_INSTANCES_DIR="$DEPLOY_ROOT/instances"


log_ok "Read and validated configuration file. "