#!/bin/bash

#ADDR=archway10e428zzqkqpxeexghyg44lpgrmrz0mxqhav0022qaqcvpffywnqswl7y53
ADDR=archway1g32hkj5hh5lcmwztg5nzzr68zf4n04ve04zn9dujygdqrg0dpsgs0t34kx
DATA='{"register": {"name": "fred"}}'

#RPC=http://49.12.133.195:26657
RPC=https://rpc.torii-1.archway.tech:443

archwayd tx wasm execute $ADDR "$DATA" \
    --chain-id=torii-1 \
    --from apple \
    --gas=auto \
    --amount=1utorii \
    --gas-prices=0.25torii --gas-adjustment 1.3 \
    --node=$RPC < <(echo @validator2022)
