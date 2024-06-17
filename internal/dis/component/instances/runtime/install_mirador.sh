#!/bin/bash
set -e

echo "=> Setting up filesystem permissions"
chmod 777 /var/www/data/project/web/sites/default/
trap "chmod 755 /var/www/data/project/web/sites/default/" EXIT

echo "=> Creating 'sites/default/libraries/wisski-mirador-integration/' directory"
mkdir -p /var/www/data/project/web/libraries/wisski-mirador-integration

echo "=> Downloading 'mirador-integration.js'"
curl -L https://raw.githubusercontent.com/rnsrk/wisski-mirador-integration/main/mirador-integration.js -o /var/www/data/project/web/libraries/wisski-mirador-integration/mirador-integration.js

echo "=> Done"