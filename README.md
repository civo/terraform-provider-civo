Terraform Provider ![Travis build](https://travis-ci.org/civo/terraform-provider-civo.svg?branch=master)
==================

- Website: https://www.terraform.io
- [![Gitter chat](https://badges.gitter.im/hashicorp-terraform/Lobby.png)](https://gitter.im/hashicorp-terraform/Lobby)
- [![Build Status](https://github.com/civo/terraform-provider-civo/workflows/Go/badge.svg)](https://github.com/civo/terraform-provider-civo/actions)
- Mailing list: [Google Groups](http://groups.google.com/group/terraform-tool)

Requirements
------------

-   [Terraform](https://www.terraform.io/downloads.html) 0.13.x
-   [Go](https://golang.org/doc/install) 1.14.x (to build the provider plugin)

Building The Provider
---------------------

Clone repository to: `$GOPATH/src/github.com/civo/terraform-provider-civo`

```sh
$ mkdir -p $GOPATH/src/github.com/terraform-providers; cd $GOPATH/src/github.com/terraform-providers
$ git clone git@github.com:terraform-providers/terraform-provider-civo
```

Enter the provider directory and build the provider

```sh
$ cd $GOPATH/src/github.com/terraform-providers/terraform-provider-civo
$ make build
```

Using the provider
----------------------

When the provider is out of beta the documentation will be at [Civo Provider documentation](https://registry.terraform.io/providers/civo/civo/latest/docs),
but during the beta the best resource for learning about it is [this guide](https://www.civo.com/learn/using-the-civo-terraform-provider)


Developing the Provider
---------------------------

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (version 1.11+ is *required*). You'll also need to correctly setup a [GOPATH](http://golang.org/doc/code.html#GOPATH), as well as adding `$GOPATH/bin` to your `$PATH`.

To compile the provider, run `make build`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

```sh
$ make build
...
$ $GOPATH/bin/terraform-provider-civo
...
```

In order to test the provider, you can simply run `make test`.

```sh
$ make test
```

In order to run the full suite of acceptance tests, run `make testacc`.

*Note:* Acceptance tests create real resources, and often cost money to run.

```sh
$ make testacc
```

In order to run a specific acceptance test, use the `TESTARGS` environment variable. For example, the following command will run `TestAccCivoDomain_Basic` acceptance test only:

```sh
$ make testacc TESTARGS='-run=TestAccCivoDomain_Basic'
```

For information about writting acceptance tests, see the main Terraform [contributing guide](https://github.com/hashicorp/terraform/blob/master/.github/CONTRIBUTING.md#writing-acceptance-tests).
