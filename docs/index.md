---
layout: ""
page_title: "Provider: Civo"
description: |-
  The Civo provider is used to interact with the resources supported by Civo. The provider needs to be configured with the proper credentials before it can be used.
---

# Civo Provider

The Civo provider is used to interact with the resources supported by Civo. 

Use the navigation to the left to read about the available resources.

## Configuration

The provider will use the credentials of the [Civo CLI](https://github.com/civo/cli) (stored in ` ~/.civo.json`) if no other credentials have been set up. The provider will use credentials in the following order:

1. Environment variable (`CIVO_TOKEN`).
1. Token provided via a credentials file (See `credentials_file` input [below](#credentials_file))
1. [CLI](https://github.com/civo/cli) configuration (`~/.civo.json`)

That means that if the `CIVO_TOKEN` variable is set, all other credentials will be ignored, and if the `credentials_file` is set, that will be used over the CLI credentials.

### Obtaining a token

First you will need to create a [Civo Account](https://dashboard.civo.com/signup) and then you can do the following:

* If you can want to configure your credentials with the CLI, instructions are [here](https://www.civo.com/docs/overview/civo-cli#add-an-api-key-to-civo-cli)
* To fetch an API key go to the [security section](https://dashboard.civo.com/security) on the dashboard


### Using the CIVO_TOKEN variable

To use the Civo token, export the variable containing your token:

```bash
export CIVO_TOKEN=<your token>
```

### Using a credentials file

The format of the credentials file is as follows:

```json
{
	"apikeys": {
		"tf_key": "write-your-token-here"
	},
	"meta": {
		"current_apikey": "tf_key"
	}
}
```

you will then need to configure the `credentials_file` input to the correct location, for example:

`credentials_file = "/secure/path/civo.json"`

### Using the CLI

If you install the CLI and [configure a token](https://www.civo.com/docs/overview/civo-cli#add-an-api-key-to-civo-cli), there is nothing else you need to do if those are the credentials you wish to use, ideal for local usage. 


## Example Usage

### Simplest usage

In this example the provider will look for credentials set by the CLI and use LON1 region as default, or use the environment variable `CIVO_TOKEN` if set.

```terraform

terraform {
  required_providers {
    civo = {
      source = "civo/civo"
    }
  }
}

provider "civo" {
  region = "LON1"
}

```

### Example with credentials file

In this example we are providing a specific version of the terraform provider and setting a credentials file to use. The credentials file will be used if the environment variable `CIVO_TOKEN` is not set.

```terraform
terraform {
  required_providers {
    civo = {
      source = "civo/civo"
      version = "1.1.0"
    }
  }
}

provider "civo" {
  credentials_file = "/secure/path/civo.json"
  region = "LON1"
}
```

## Argument Reference

### Optional

- `api_endpoint` (String) The Base URL to use for CIVO API.
- `region` (String) This sets the default region for all resources. If no default region is set, you will need to specify individually in every resource.
<a id="credentials_file"></a>
- `credentials_file` (string) specify a location for a file containing your civo credentials token 
- `token` (String) (**Deprecated**) for legacy reasons the user can still specify the token as an input, but in order to avoid storing that in terraform state we have deprecated this and will be remove in future versions - don't use it.
