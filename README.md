Terraform VKCS Provider
=======================

* Documentation [https://registry.terraform.io/providers/vk-cs/vkcs/latest/docs](https://registry.terraform.io/providers/vk-cs/vkcs/latest/docs)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.1.5 and later
-	[Go](https://golang.org/doc/install) 1.22 (to build the provider plugin)

Using The Provider
------------------
To use the provider, prepare configuration files based on examples from [here](https://github.com/vk-cs/terraform-provider-vkcs/tree/master/examples)

```sh
$ cd $GOPATH/src/github.com/vk-cs/terraform-provider-vkcs/examples/create-vkcs-compute-instance
$ vim provider.tf
$ terraform init
$ terraform plan
```

Provider development
--------------------
To start improving grab the repository, build it and install into local registry repository.
Builds for MacOS, Windows and Linux are available.
The example is for MacOS.
```sh
$ mkdir -p $GOPATH/src/github.com/vk-cs
$ cd $GOPATH/src/github.com/vk-cs
$ git clone git@github.com:vk-cs/terraform-provider-vkcs.git
$ cd $GOPATH/src/github.com/vk-cs/terraform-provider-vkcs
$ make build_darwin
$ mkdir -p ~/.terraform.d/plugins/hub.mcs.mail.ru/repository/vkcs/0.1.0/darwin_amd64/
$ cp terraform-provider-vkcs_darwin ~/.terraform.d/plugins/hub.mcs.mail.ru/repository/vkcs/0.1.0/darwin_amd64/terraform-provider-vkcs_v0.1.0

$ cat <<EOF > main.tf 
terraform {
  required_providers {
    vkcs = {
      source  = "hub.mcs.mail.ru/repository/vkcs"
      version = "~> 0.1.0"
    }
  }
}
EOF
$ terraform init
```

When submitting PR make sure that if golang code has been changed, PR has updates to CHANGELOG.md. Add description of changes under last version with "(unreleased)" mark.

Documenting provider
--------------------
To update documentation contents, please, update "description" field of necessary resource/data_source schema and create/modify documentation templates.
Documentation templates are located in templates/ folder.
PR with renewed provider documentation is generated automatically when updates are merged into master branch.

Thank You!
