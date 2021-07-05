# Samples - boostrap a two validator substrate-based network

## onf-cli version
`>= 0.3.3`

## Image
We use substrate v3 release image parity/substrate:v3.0.0
Should be fine to replace it with new version

## Chainspec
This example we export the built-in local chainspec to chainspec.json.
You are free to make changes and use your own pre-prepared chainspec.

## Steps - 1 Prepare chainspec.json
The `chainspec.json` in this folder is generated via `docker run --rm parity/substrate:v3.0.0 build-spec --chain=local > chainspec.json`

## Steps - 2 bootstrap config (optional)
You can change validator count or add dedicated bootnodes by updating bootstrap-config.yaml

## Steps - 3 Launch
run `onf network bootstrap -f bootstrap-config.yaml`

## Follow-ups
The network should start producing blocks, you can connect polkadot apps to your node and see.
Also you can start full node in onfinality ui as well as cli, and new nodes will connect to alice and sync by default.
