#!/bin/bash

set -euo pipefail

# cat /etc/init.d/krb5-admin-server

if [[ ! -f /var/lib/krb5kdc/principal ]]; then
    kdb5_util -r KERB.LOCAL -P hunter2 create -s
fi
# kdb5_util create -s
# krb5_newrealm

exec /usr/sbin/krb5kdc -n
