# Samples

The sample folder contains the definition files for different purposes.

    Notice: 
    Before you can create a node in the OnFinality platform, you have to bind a credit card to your workspace. Or please contact to the OnFinality Support for business enquiry.


## Example: Setup a testnet

In the bootstrap-config.json file, we define a new network spec, and define 3 validator nodes and 1 bootnode. Each validator node need to bind a group of session keys what we define in the `sessionKeys` section.

the init session keys are sensitive data, so please DO NOT commit your config file to any public repo.   

```bash 
onf network bootstrap -f ./bootstrap-config.json
```

after run the command above, you can check all the nodes running in your workspace with follow command. 

```bash
onf node list
```
Or login to the https://app.onfinality.io to view nodes in UI.