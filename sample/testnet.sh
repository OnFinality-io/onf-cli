#!/usr/bin/env bash

CMD="go run ../cmd"

#networkId=$($CMD network-spec create -f ./custom-testnet.json -o json | jq -r ".key")
#echo "$networkId"

networkId="dba52606-71b9-4d70-9139-5b7f04698db2"
#networkId="3ddee574-448f-4500-9f2b-cfbc88eb9270"
$CMD network-spec upload -n "$networkId" -f ./chainspec.json