#!/usr/bin/env bash
set -e
if [ "$#" -eq 0 ]; then
  echo "Usage:"
  echo "  ./convert_p12.sh mumble_client_cert.p12 export_prefix"
  echo "  this will create export_prefix.crt.pem and export_prefix.key.pem that can be used"
  echo "  hit enter when asked for password"
  exit 1
fi
FILE=$1
PREFIX=$2
set -x
openssl pkcs12 -in $FILE -out $PREFIX.crt.pem -clcerts -nokeys
openssl pkcs12 -in $FILE -out $PREFIX.key.pem -nocerts -nodes
