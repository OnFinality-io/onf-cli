# Samples boostrap parachain connect to relay chain

## onf-cli version
`>= 0.3.5`

## Step - 1 Prepare your chainspec file

We provide some sample chainspec files in the folder `rococo-local.json` (relay) and `karura-dev-2000.json` (para) that you can use directly; or you can generate the chainspec file in your way.


## Step - 2 boostrap relaychain

run `onf network bootstrap -f bootstrap-relaychain-config.yaml`

After you run the command, you will create one network spec in the workspace, a new custom network with 3 new validator nodes. If you want to add bootnodes for the new network, you can adjust the `count` property in the bootNode section. Then all the new created bootnodes libp2p address will be updated to the network spec.


## Step - 3 boostrap parachain. 

Firstly, you need to modify `--bootnodes` parameters in the networkspec definition section to replace with addresses from the bootnodes you created in the previous step. you can get the libp2p address in the OnFinality webapp, or via cli tool.

run `onf network bootstrap -f bootstrap-parachain-config.yaml`

Now you have another 3 validator nodes for the parachains.
