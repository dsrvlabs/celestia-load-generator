#!/bin/bash

RPC=https://rpc.torii-1.archway.tech:443

archwayd q tx $1 \
    --chain-id torii-1 \
    --node $RPC

