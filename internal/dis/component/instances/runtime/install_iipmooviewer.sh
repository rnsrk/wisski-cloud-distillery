#!/bin/bash
set -e

echo "=> Setting up filesystem permissions"
chmod 777 /var/www/data/project/web/sites/default/
trap "chmod 755 /var/www/data/project/web/sites/default/" EXIT

echo "=> Creating 'sites/default/libraries/iipmooviewer' directory"
mkdir -p /var/www/data/project/web/libraries/iipmooviewer

echo "=> Download 'IIPMooViewer'"
curl -L https://raw.githubusercontent.com/ruven/iipmooviewer/master/js/iipmooviewer-2.0-min.js -o /var/www/data/project/web/libraries/iipmooviewer/iipmooviewer-2.0-min.js

echo "=> Done"


