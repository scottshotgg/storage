#!/bin/bash

# Generate bin and abi files from contract in Crud.sol
solc contracts/Crud.sol --pretty-json --combined-json abi,asm,ast,bin,bin-runtime,compact-format,devdoc,hashes,interface,metadata,opcodes,srcmap,srcmap-runtime,userdoc --abi --bin --optimize-runs 999999999 --overwrite --gas -o contracts/

# Make etherweb package dir
mkdir -p crud

# Generate Go binding for contract
abigen --abi contracts/Crud.abi --pkg crud --type Crud --out crud/crud.go