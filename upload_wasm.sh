#!/bin/bash

#set -x

WASM_FILE=$1
#PASSWD=$2
ACCOUNT=$3
CHAIN_ID=$4
RPC=$5
ARCHWAY_HOME=$6

#echo "===== ===== ===== ===== ====="
#echo $1
#echo $2
#echo `cat $2`
#echo $3
#echo $4
#echo $5
#echo $6
#echo "===== ===== ===== ===== ====="

RESP=`archwayd tx wasm store $WASM_FILE \
    --chain-id $CHAIN_ID \
    --from $ACCOUNT \
    --home ~/.archway \
    --gas auto \
    --broadcast-mode sync \
    --node $RPC -y --output json < ./passwd`

#RESP=`archwayd tx wasm store $WASM_FILE \
#    --chain-id $CHAIN_ID \
#    --from $ACCOUNT \
#    --home ~/.archway \
#    --gas auto \
#    --broadcast-mode sync \
#    --node $RPC -y --output json < ./passwd`

#echo $RESP

#echo "===== ===== ===== ===== ====="
TXHASH=`echo $RESP | jq .txhash | sed -e 's/^"//' -e 's/"$//'`
echo $TXHASH
