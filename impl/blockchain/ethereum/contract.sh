#!/bin/bash

# Generate bin and abi files from contract in Etherweb.sol
solc impl/blockchain/ethereum/contracts/Etherweb.sol --pretty-json --combined-json abi,asm,ast,bin,bin-runtime,compact-format,devdoc,hashes,interface,metadata,opcodes,srcmap,srcmap-runtime,userdoc --abi --bin --optimize-runs 999999999 --overwrite --gas -o contracts/

# Make etherweb package dir
mkdir -p etherweb

# Generate Go binding for contract
abigen --abi contracts/Etherweb.abi --pkg etherweb --type Etherweb --out etherweb/etherweb.go