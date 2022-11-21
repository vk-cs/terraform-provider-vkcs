# Examples

## Audience of examples

- Users

    Examples show terraform users how to create vkcs objects with various features.
    Users may just look inside examples or run them to get trained.

- Release/development automatization

    Examples are used by canary tests on vkcs side to verify the latests terraform
    and vkcs provider over vkcs cloud.

    Also examples are used by developers manually or via CI to verify changes.

- Documentation generator

    Documentation is automatically generated using examples that are located here
    and with additional examples from templates folder.

## Running an example

- Choose an example
- Create a folder on a local host
- Set up vkcs terraform provider in the folder using instruction on your account page on vkcs site.
- Download all files of the example into the folder.
- Run `terraform plan` to see what resources with what parameters are going to be deployed onto vkcs.
- Run `terraform apply`. This will deploy the example onto vkcs.
- In case you want to change deployed resources you change example files and run `terraform plan` to see
what is going to be changed on vkcs.
- Then run `terraform apply` again to deploy changes.
- When example deployment is no more needed run `terraform destroy` to destroy it on vkcs.

## Creating a new example

An example must be consistent. I.e. it must describe all required resources/datasources
to run the example in a separate folder.

An example must contain `main.tf` file.

`main.tf` file must contain enough resources to demonstrate an idea of the example itself.
E.g. `main.tf` of compute resource example must contain compute resource.

`main.tf` file must **not** contain resources/datasource which do not direcly relate
to an idea of the example itself.
E.g. `main.tf` of compute resource example must not contain network resources.

Additional resources/datasources of an example must be specified in additional file of
the example. At the moment all such examples use the single `base.tf` file for this purpose,
but this may be changed to several files (`base-network.tf`, `base-compute.tf`, etc).
E.g. compute resource example must place network resources, flavor, and image datasources
into a separate file (or files).

Use only public common cloud objects. I.e. neither stage AZ, nor custom flavors/volume types.

Do not use IDs of common cloud objects. Use datasources instead.
