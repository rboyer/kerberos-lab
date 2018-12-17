#!/bin/bash

set -euo pipefail

set -x
exec docker run --rm -it --net='container:kdc' local/kerb-client bash
