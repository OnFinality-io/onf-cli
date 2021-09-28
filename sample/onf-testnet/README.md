# Samples boostrap parachain connect to relay chain

## onf-cli version
`>= 0.3.5`

## Steps - 1 boostrap relaychain
run `onf network bootstrap -f bootstrap-relaychain-config.yaml`

## Steps - 2 Get relaychain p2p address 
Through command message or login in Onfinality platform to get p2p address

## Steps - 3 boostrap parachain. First modify yaml `--bootnodes` parameters, the address is relaychain p2p address
run `onf network bootstrap -f bootstrap-parachain-config.yaml`
