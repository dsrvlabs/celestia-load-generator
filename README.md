# WASM Load Generator

## Description

This repository provides Load Generator tools for Archway testnet challenge.

There are two types of test that are supported by this tool.
- upload contract test.
- call contract test.

Basically, this tool creates multiple txs continuously with multiple accounts.
When test start, go routines will create per account and every go routines create txs simuletainously.

**Upload Contract**
This repository provides sample wasm contract for testing, but you can choose another contract if you want.

**Call Contract**
Sample contract which makes high CPU load is already deployed. But you can also choose another contract address.

## Pre-requisites
- golang +1.18 for building application.
- `archwayd` should be installed for sending TX.
- Account addresses for testing with

## How to Use

### Step 1: Build load generator

Clone repository.

```
~$ git clone https://github.com/rootwarp/wasm-load-generator
```

And build.

```
~$ cd wasm-load-generator
~$ go build
~$ go install
```

### Step 2: Prepare wallets on keyring

To create a lot of txs, you need available account as many as possible.
And make sure that keys are stored in same machine.

### Step 3: Set password file

You need to set passwd file to unlock your wallet when the application try to send tx.

```
~$ echo "<FILL YOUR PASSWORD>" > passwd
```

### Step 4: Set account file

As I mentioned, this load generator uses multiple accounts for testing simuletaneously.

For example, `accounts.txt` file like below.

```
~$ cat accounts.txt

archway1jduy83242hv60p4k4kn8dfx9mv95qgrq9lrpt9
archway1pm0yyd2ncc2x67ctuz5p3tcxa59tezx5scp0hj
archway179jgnt5ckmnjrxg2rykycv4vkh2y8nh8h30sga
archway1fqdch0dl4r43wp3yw6f5nkp7jca3xn4m3s2plh
archway13yazsem53w0lpsf0j0n0l7038jvs5xadr85zrt
```

### Step 5: Run test

#### Upload Test

This test upload wasm contract file with multiple account continuously.
For this test, tester needs contract file to be uploaded, password and account files as described above.

For simplicity, this repository contains `test_contract.wasm` file to test.

Run command below to start test, you can change flags according to your environments.

```
~$ wasm-load-tester upload --wasm ./test_contract.wasm --password ./passwd --account ./accounts.txt --chain-id torii-1 --node https://rpc.torii-1.archway.tech:443
```

#### Contract Call Test

This test calls heavy contract to test cpu load of validator.
For testing, prepared contract should be uploaded and instantiated first.

For convenience, test contract prepared at `archway15mha747mukh5un0nsw5jgn7et3geujc6nfp2atppj6y68378z55qqarxyk`.

Run command below to start test.

```
~$ ./wasm-load-tester call --account ./account_single.txt --chain-id torii-1 --contract archway15mha747mukh5un0nsw5jgn7et3geujc6nfp2atppj6y68378z55qqarxyk --password ./passwd --node https://rpc.torii-1.archway.tech:443
```

## TODOs
- Remove script code.
- Fix TX sequence bug.
