# yocoin

## Build binary
```shell script
    make yocoin
```

## Starting node

1. Create new account(-s)
```
    build/bin/yocoin account new
```

1. Create **genesis.json** file from **genesis.exaple.json** and add new account to **"alloc"** block


2. Initialize private net
```
    build/bin/yocoin --datadir "./datadir" init genesis.json
```

3. Run new node with console mode
```
    build/bin/yocoin --datadir "./datadir" --maxpeers 3 --port "30301" --rpc --rpccorsdomain "*" console
```

4. Chech mining base address
```
    > eth.coinbase
```
4. Start mining to generate DAG
```
    > miner.start()
```
