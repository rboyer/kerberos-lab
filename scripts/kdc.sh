#!/bin/bash

set -euo pipefail

if [[ $# -lt 1 ]]; then
    echo "missing required arg" >&2
    exit 1
fi

case "$1" in
    kdc)
        if [[ ! -f /var/lib/krb5kdc/principal ]]; then
            kdb5_util -r KERB.LOCAL -P hunter2 create -s
        fi
        exec /usr/sbin/krb5kdc -n
        ;;
    kadmin)
        if [[ ! -f /var/lib/krb5kdc/init_done ]]; then
            echo "addprinc -pw admin admin" | kadmin.local
            echo "addprinc -pw demo demo" | kadmin.local
            touch /var/lib/krb5kdc/init_done
        fi
        exec /usr/sbin/kadmind -nofork
        ;;
    *)
        echo "unknown command: $1" >&2
        exit 1
        ;;
esac
