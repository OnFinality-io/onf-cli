# Samples - boostrap a dedicated rococo for your team

## onf-cli version
`>= 0.3.3`

## Image
We use latest polkadot image parity/polkadot:v0.9.7
Should be fine to replace it with new version

## Chainspec
We use the built-in rococo-local which defines a two-validator network (alice and bob)

## Steps - 1 Create the network in your workspace
run `onf network-spec create -f create-network.yaml` and note down network-spec's key

## Steps - 2 Launch Alice
update `node-alice.yaml` with your network-spec's key and

run `onf node create -f node-alice.yaml` and note down Alice's nodeId

## Steps - 3 Update network-spec to make Alice bootnode
update `update-network.yaml` with alice's nodeId

run `onf network-spec update -f update-network.yaml -n <your network-spec's key>`

## Step - 4 Launch Bob
update `node-bob.yaml` with your network-spec's key and

run `onf node create -f node-bob.yaml`

## Follow-ups
The network should start producing blocks, you can connect polkadot apps to your node and see.
Also you can start full node in onfinality ui as well as cli, and new nodes will connect to alice and sync by default.
