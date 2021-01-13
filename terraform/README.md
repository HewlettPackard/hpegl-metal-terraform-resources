# Building

```
cd terraform
make install
```

# Testing
 Quick tests can be executed using

 ```
 make test
 ```

 Running Terraform acceptance level tesing requires that there's a valid token and login to a portal (or the simulator). The tests also asusume the availablility of some specific image names etc.

```
 make acceptance
```

 with example output like

```

 go test -v -i $(go list ./quake | grep -v vendor) 
echo $(go list ./quake | grep -v vendor) | \
	TF_ACC=true xargs -t -n4 go test -v  -timeout=60s -cover
go test -v -timeout=60s -cover github.com/quattronetworks/quake-client/terraform/quake 
=== RUN   TestAvailableResourcesBasic
--- PASS: TestAvailableResourcesBasic (0.27s)
=== RUN   TestAccImages_Basic
--- PASS: TestAccImages_Basic (0.22s)
=== RUN   TestAccUsages_Basic
--- PASS: TestAccUsages_Basic (0.25s)
=== RUN   TestProvider
--- PASS: TestProvider (0.00s)
=== RUN   TestProviderInterface
--- PASS: TestProviderInterface (0.00s)
=== RUN   TestAccQuakeHost
--- PASS: TestAccQuakeHost (0.39s)
=== RUN   TestAccQuattroSSHKey_Basic
--- PASS: TestAccQuattroSSHKey_Basic (0.20s)
=== RUN   TestAccQuattroVolume
--- PASS: TestAccQuattroVolume (6.25s)
PASS
coverage: 67.2% of statements
ok  	github.com/quattronetworks/quake-client/terraform/quake	7.600s	coverage: 67.2% of statements
```

# Using GreenLake tokens

The terraform provider is now capable of using GreenLake tokens.  The information required is stored in a .gltform file.
This file can be written in home or in the directory from which terraform is run.  The file contents:
 
```yaml
space_name: <...>
rest_url: http://localhost:3002
project_id: 65c82181-fefc-4ea7-870e-628225fe7664
access_token: <...>
```

The first field `space_name` is optional, and is only required if the terraform provider is going to be used to create
projects.

To make the provider use GreenLake tokens - i.e. use the information in the .gltform file - the gl_token field must be set
to true in the provider definition stanza:

```hcl
provider "quake" {
  gl_token = true
}
```

Notes:
* Note that the project referred to must be in a GreenLake Organization.
* If creating a project, be warned that ippools information cannot be provided as input (a restriction in the
    client).  Without ippools information host creation will fail.
  

## Setting up simulator to use GreenLake tokens

### Add AWS credentials and scmClientConfig.yaml file to simulator

To enable the Simulator to handle GreenLake tokens you need to add the following to the Simportal container:

* The following file scmClientConfig.yaml:
   ```yaml
    aws:
      creds:
        mount: ./creds
      secrets:
        namespace: integ
    hpe:
      scm:
        url: https://iam.intg.hpedevops.net
      identity:
        url: https://iam.intg.hpedevops.net
    client:
      serviceName: bmaas
      serviceSpace: _hpe_bmaas
      timeout: 30
      base:
        tenantId: root
      log:
        level: info
        formatter: text
    ```

* The ./creds directory must contain AWS creds, the contents are as follows:
    ```bash
    -rw-rw-r-- 1 eamonn eamonn 21 Jan  8 09:26 AWS_ACCESS_KEY_ID
    -rw-rw-r-- 1 eamonn eamonn 10 Jan  8 09:25 AWS_REGION_NAME
    -rw-rw-r-- 1 eamonn eamonn 41 Jan  8 09:26 AWS_SECRET_ACCESS_KEY
    ```
    The contents can be obtained from Bret McKee

* Add the following to steeld_config.yml:
    ```yaml
    glhc_iam_auth_endpoint: "iam.intg.hpedevops.net:443"
    glhc_iam_config_file: "./scmClientConfig.yaml"
    ```

### Add GreenLake organization to portal

A GreenLake organization needs to be added.  The tenant-id to use currently is: bpcfi136mjshu5h8g2sg

### Use GreenLake cli to login to Greenlake, switch to tenant, and get token

The GreenLake cli can be found [here](https://github.com/hpe-hcss/hpecli).  Once it is installed, do the following:

```bash
$ hpe gl login
$ hpe gl login-tenant switch bpcfi136mjshu5h8g2sg
$ hpe gl iam tokens show
```

The last command displays the token to use.  This token needs to be added to the following files:
* ~/.gltform as access_token
* ~/.gljwt: this file is used by qcli to create a project under the GreenLake organization.  The format of the file is:
    ```json
    { "accessToken": {"accessToken":  <....> }}
    ```

### Use qcli to create a project in the GreenLake organization

If the GreenLake organization doesn't have a project, qcli must be used to create one.  First of all use `hpe gl iam spaces list`
to list the Quake tenant GreenLake space.  Get the name (currently "Quake Test Default").  Then use qcli to:

* Create a session for use with the GreenLake token obtained above - use one of the project-ids in the list
    returned from `hpe gl iam spaces list`, the format of the command is:
    ```bash
    $ ./qcli session -g -f ~/.gljwt -u <rest-url> -p <project id>
    ```
* Create a project using the above session, the format of the command is:
    ```bash
    $ ./qcli -d projects create ./test-project.json --space "Quake Test Default"
    ```
    Where the contents of test_project.json are:
    ```json
    {
      "Name": "BT NEM Project",
      "Profile": {
        "TeamName": "Test Team",
        "TeamDesc": "Test Team",
        "Company": "",
        "Address": "Austin, TX",
        "Email": "quaketestbot@gmail.com",
        "EmailVerified": true,
        "PhoneNumber": "5125551212",
        "PhoneVerified": true
      },
      "Limits": {
        "Hosts": 20,
        "Volumes": 20,
        "VolumeCapacity": 0
      }
    }
    ```