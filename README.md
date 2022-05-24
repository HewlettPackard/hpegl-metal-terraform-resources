# hpegl-metal-terraform-resources

- [hpegl-metal-terraform-resources](#hpegl-metal-terraform-resources)
- [Terraform resources for HPEGL Metal](#terraform-resources-for-hpegl-metal)
  - [Requirements](#requirements)
  - [Usage](#usage)
  - [Building the resources as a stand-alone provider](#building-the-resources-as-a-stand-alone-provider)
  - [Testing stand-alone provider](#testing-stand-alone-provider)
    - [Unit tests](#unit-tests)
    - [Acceptance tests](#acceptance-tests)

# Terraform resources for HPEGL Metal

Terraform Metal resources is a plugin for HPEGL terraform provider that allows the full lifecycle management of HPEGL
Metal resources. This provider is maintained by [HPEGL Metal resources team](mailTo:quake-core@hpe.com).

## Requirements

1. Terraform version >= v0.13 [install terraform](https://learn.hashicorp.com/tutorials/terraform/install-cli)
2. A Service Client to authenticate against GreenLake.
3. Terraform basics. [Terraform Introduction](https://www.terraform.io/intro/index.html)

## Usage

See the Terraform provider for
hpegl [documentation](https://registry.terraform.io/providers/HewlettPackard/hpegl/latest/docs)
to get started using the provider.

## Building the resources as a stand-alone provider

```bash
$ git clone https://github.com/hewlettpackard/hpegl-metal-terraform-resources.git
$ cd hpegl-metal-terraform-resources
$ make build
```

Note: For debugging the provider please refer to the
[debugging guide](https://medium.com/@gandharva666/debugging-terraform-using-jetbrains-goland-f9a7e992cb1d)

## Testing stand-alone provider

### Unit tests
Unit tests can be executed using
 ```
 make test
 ```

### Acceptance tests
Running Terraform acceptance level testing requires a Metal service endpoint and a Project_Owner membership.  
The tests as of now work with a Metal simulator and assume that the required environment is already available.
* Hoster _**TestHoster1**_ and  Project **_TestTeam1_**
* Metal issued JWT tokens
* Also assumes the availability of certain resources like services, networks, etc.   

**To run the acceptance test,**
1. When the Metal token is used, the Plugin reads the token, URL, and membership details from the file  _**~/.qjwt**_.
Create this file in the format:

```yaml
rest_url: http://172.25.0.2:3002
user: h1@hpe.com
jwt: eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCIsImtpZCI6IlJFTk.dlfkjsj.dfsdf
member_id: 835590C1-AFF7-438B-BBBD-D6184157CB41
```

2. Set the environment variable HPEGL_METAL_GL_TOKEN to `false` to indicate Metal authentication mode.
```bash
export HPEGL_METAL_GL_TOKEN=false
```

3. Run test
```
make acceptance
```
