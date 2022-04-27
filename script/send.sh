#!/bin/bash

VAL_ADDR=archway1pm0yyd2ncc2x67ctuz5p3tcxa59tezx5scp0hj

send_balance() {
    echo "Send from $1 / $2"
    RET=`archwayd tx bank send $1 $VAL_ADDR ${2}utorii \
        --chain-id torii-1 \
        --node https://rpc.torii-1.archway.tech:443 \
        --gas auto \
        --output json \
        -y < passwd`

    CODE=`echo $RET | jq .code`
    TX=`echo $RET | jq .txhash`

    echo $CODE $TX
}

