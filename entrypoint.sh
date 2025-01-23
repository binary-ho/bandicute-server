#!/bin/sh
set -e

mkdir -p /root/.config/sops/age

echo "$SOPS_PRIVATE_KEY" > /root/.config/sops/age/keys.txt
chmod 600 /root/.config/sops/age/keys.txt

sops -d /app/config/encrypted-property.yml > /app/config/property.yml

./main
