Terraform VKCS Provider
============================

* Documentation [https://registry.terraform.io/providers/VKCloudSolutions/vkcs/latest/docs](https://registry.terraform.io/providers/VKCloudSolutions/vkcs/latest/docs)

Requirements
------------

-	[Terraform](https://www.terraform.io/downloads.html) 1.0.x
-	[Go](https://golang.org/doc/install) 1.17 (to build the provider plugin)

Using The Provider
----------------------
To use the provider, prepare configuration files based on examples from [here](https://github.com/vk-cs/terraform-provider-vkcs/tree/master/examples)

```sh
$ cd $GOPATH/src/github.com/vk-cs/terraform-provider-vkcs/examples/create-vkcs-compute-instance
$ vim provider.tf
$ terraform init
$ terraform plan
```

Provider development
---------------------
To start improve it grab the repository, build it and install into local registry repository.
Builds for MacOS, Windows and Linux are available.
The example is for MacOS.
```sh
$ mkdir -p $GOPATH/src/github.com/vk-cs
$ cd $GOPATH/src/github.com/vk-cs
$ git clone git@github.com:vk-cs/terraform-provider-vkcs.git
$ cd $GOPATH/src/github.com/vk-cs/terraform-provider-vkcs
$ make build_darwin
$ mkdir -p ~/.terraform.d/plugins/hub.vkcs.mail.ru/repository/vkcs/0.5.8/darwin_amd64/
$ cp terraform-provider-vkcs_darwin ~/.terraform.d/plugins/hub.vkcs.mail.ru/repository/vkcs/0.5.8/darwin_amd64/terraform-provider-vkcs_v0.5.8

$ cat <<EOF > main.tf 
terraform {
  required_providers {
    vkcs = {
      source  = "hub.vkcs.mail.ru/repository/vkcs"
      version = "~> 0.5.8"
    }
  }
}
EOF
$ terraform init
```

Publishing provider
-------------------
Provider publishes via action [release](https://github.com/vk-cs/terraform-provider-vkcs/blob/master/.github/workflows/release.yml).
To call the action create new tag.
```sh
$ git tag v0.5.8
$ git push origin v0.5.8
```

Thank You!
