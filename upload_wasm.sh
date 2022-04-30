#!/bin/bash

#set -x

WASM_FILE=$1
#PASSWD=$2
ACCOUNT=$3
CHAIN_ID=$4
RPC=$5
ARCHWAY_HOME=$6

RESP=`archwayd tx wasm store $WASM_FILE \
    --chain-id $CHAIN_ID \
    --from $ACCOUNT \
    --home ~/.archway \
    --gas auto \
    --broadcast-mode sync \
    --node $RPC -y --output json < ./passwd`

TXHASH=`echo $RESP | jq .txhash | sed -e 's/^"//' -e 's/"$//'`
echo $TXHASH
