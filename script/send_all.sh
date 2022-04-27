#!/bin/bash

source ./get_balances.sh
source ./send.sh

ACCS=`archwayd keys list \
    --home ~/.archway \
    --output json < passwd | jq .[].address | sed -e 's/^"//' -e 's/"$//'`

for ACC in ${ACCS[@]}; do
    get_balance $ACC

    if [ "$AMOUNT" != null ]; then
        echo "$ACC => Send"
        send_balance $ACC $AMOUNT
    else
        echo "$ACC => PASS"
    fi
done

