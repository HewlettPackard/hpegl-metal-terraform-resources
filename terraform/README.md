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

