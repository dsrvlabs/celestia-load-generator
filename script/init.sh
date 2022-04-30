#!/bin/bash

DATA='{"purchase_price":{"amount":"100","denom":"utorii"},"transfer_price":{"amount":"999","denom":"utorii"}}'

archwayd tx wasm instantiate 14938 "$DATA" \
    --chain-id=torii-1 \
    --from=cheese \
    --label="test" \
    --gas=auto \
    --no-admin \
    --node=https://rpc.torii-1.archway.tech:443
