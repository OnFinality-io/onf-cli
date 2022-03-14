# OnFinality Command Line Tool

OnFinality is a SaaS Platform to support blockchain developers by providing infrastructure and API services so you can automate blockchain DevOps.
this CLI tool can help users to interact with OnFinality platform in developer prefer way. 

## Features

- manage dedicated nodes
- manage network specs
- check platform information


## Installation
```
curl -s https://raw.githubusercontent.com/OnFinality-io/onf-cli/master/scripts/install/install.sh | bash
```
#### Note:
If you encountered a permission denied failure from the above command, please refer to the command below instead:
```
curl -s https://raw.githubusercontent.com/OnFinality-io/onf-cli/master/scripts/install/install.sh | sudo bash
```

## For development

```bash
git clone https://github.com/OnFinality-io/onf-cli.git
cd onf-cli
make binary-osx
```

## Env configuration

|  Key | Description  |
| ------------ | ------------ |
|  ONF_ACCESS_KEY | How to obtain: https://app.onfinality.io/account  |
|  ONF_SECRET_KEY | How to obtain: https://app.onfinality.io/account  |
| ONF_WORKSPACE_ID | Specify a workspace ID|

