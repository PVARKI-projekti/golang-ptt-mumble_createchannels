#!/usr/bin/env bash
set -e
if [ "$#" -eq 0 ]; then
  echo "Usage:"
  echo "  ./create_cert.sh cert_prefix"
  echo "  hit enter when asked for password"
  exit 1
fi
PREFIX=$1
KEYFILE=$PREFIX.key.pem
if [ -z "$CN_USER_NAME" ] ; then
  CN_USER_NAME="Mumble User"
fi

echo '[req]' >$PREFIX.csr.conf
echo 'default_bits = 2048' >>$PREFIX.csr.conf
echo 'distinguished_name = dn' >>$PREFIX.csr.conf
echo 'prompt             = no' >>$PREFIX.csr.conf
echo '[dn]' >>$PREFIX.csr.conf
echo 'CN="'$CN_USER_NAME'"' >>$PREFIX.csr.conf

echo 'basicConstraints = CA:FALSE' >$PREFIX.ext.conf
echo 'nsCertType = client' >>$PREFIX.ext.conf
echo 'nsComment = "OpenSSL Generated Client Certificate"' >>$PREFIX.ext.conf
echo 'subjectKeyIdentifier = hash' >>$PREFIX.ext.conf
echo 'authorityKeyIdentifier = keyid,issuer' >>$PREFIX.ext.conf
echo 'keyUsage = critical, nonRepudiation, digitalSignature, keyEncipherment' >>$PREFIX.ext.conf
echo 'extendedKeyUsage = clientAuth, emailProtection' >>$PREFIX.ext.conf

openssl req -new -nodes -newkey rsa:2048 -keyout $KEYFILE -out $PREFIX.csr.pem -config $PREFIX.csr.conf
openssl x509 -signkey $KEYFILE -in $PREFIX.csr.pem -req -days 3650 -out $PREFIX.crt.pem -extfile $PREFIX.ext.conf
rm $PREFIX.ext.conf $PREFIX.csr.conf $PREFIX.csr.pem
