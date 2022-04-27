#!/bin/bash

#RPC=http://49.12.133.195:26657
RPC=https://rpc.torii-1.archway.tech:443

get_balance() {
    AMOUNT=`archwayd q bank balances $1 \
        --chain-id torii-1 \
        --node $RPC \
        --output json \
        --home ~/.archway | jq .balances[0].amount | sed -e 's/^"//' -e 's/"$//'`
}

get_balance $1

echo $AMOUNT
