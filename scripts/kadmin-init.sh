#!/bin/bash

set -euo pipefail

exec /usr/sbin/kadmind -nofork
