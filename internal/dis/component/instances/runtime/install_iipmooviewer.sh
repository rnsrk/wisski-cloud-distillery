#!/bin/bash
set -e

echo "=> Setting up filesystem permissions"
chmod 777 /var/www/data/project/web/sites/default/
trap "chmod 755 /var/www/data/project/web/sites/default/" EXIT

echo "=> Download 'IIPMooViewer' repo"
wget https://github.com/ruven/iipmooviewer/archive/refs/heads/master.zip -P /var/www/data/project/web/sites/default/libraries/
echo "=> Unzip 'IIPMooViewer' repo"
unzip /var/www/data/project/web/sites/default/libraries/master.zip -d web/libraries/
echo "=> Remove 'IIPMooViewer' zipped package"
rm -r /var/www/data/project/web/sites/default/libraries/master.zip
echo "=> Rename 'IIPMooViewer' library"
mv /var/www/data/project/web/sites/default/libraries/iipmooviewer-master web/libraries/iipmooviewer

echo "=> Done"