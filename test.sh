#!/bin/bash

WASM_FILE=$1
RPC=$2

archwayd tx wasm store $WASM_FILE \
    --chain-id torii-1 \
    --from cheese \
    --home ~/.archway \
    --gas auto \
    --broadcast-mode sync \
    --node $RPC -y --output json
