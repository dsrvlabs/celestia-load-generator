#!/bin/bash

PASSWD_FILE=$1
ACCOUNT=$2
CONTRACT=$3
CHAIN_ID=$4
RPC=$5
ARCHWAY_HOME=$6

TX='{"increment":{}}'

RESP=`archwayd tx wasm execute $CONTRACT "$TX" \
    --from=$ACCOUNT \
    --chain-id=$CHAIN_ID \
    --node=$RPC \
    --gas=auto \
    --gas-prices=1utorii \
    --gas-adjustment=1.5 \
    --output=json -y < ./passwd`

TXHASH=`echo $RESP | jq .txhash | sed -e 's/^"//' -e 's/"$//'`
echo $TXHASH
