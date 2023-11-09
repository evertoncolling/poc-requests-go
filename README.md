# poc-requests-go

Repository with some experiments to perform requests to CDF in Go

## Setting up CDF authentication

Create a .env file in the root folder with the following variables (replace`xxxx` with actual values):

```bash
CLIENT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
TENANT_ID=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CDF_CLUSTER=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CDF_PROJECT=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
CLIENT_SECRET=xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx
```

## Running the code

```bash
go run main.go
```
