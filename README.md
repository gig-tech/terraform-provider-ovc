[![Build Status](https://travis-ci.org/gig-tech/terraform-provider-ovc.svg?branch=master)](https://travis-ci.org/gig-tech/terraform-provider-ovc)

# OpenvCloud Provider

The OpenvCloud provider is used for interacting with the cloud platform in order to create the desired resources. The provider needs to be configured with credentials.

## Build

To build, make sure the project is in your GOPATH and run the following command:


```sh
mkdir -p ~/.terraform.d/plugins
cd /tmp
git clone https://github.com/gig-tech/terraform-provider-ovc.git
cd terraform-provider-ovc
make build #  This uses go modules and requires go>1.11, if lower install the repo into $GOPATH and run go build
mv terraform-provider-ovc ~/.terraform.d/plugins
```

Put the binary in your plugins folder of terraform. For more information:

https://www.terraform.io/docs/plugins/basics.html#installing-plugins

## Usage

Examples of the provider can be found in the [examples](./examples) folder.

More detailed documentation on the provided resource can be found under [docs/resources.md](./docs/resources.md).  
More detailed documentation on the provided data source can be found under [docs/data_sources.md](./docs/data_sources.md).

## Authentication

Authentication is done with an [https://itsyou.online](itsyouonline) client ID and client secret pair, or with a JWT.

### Authentication with client ID and client secret pair

To configure the ovc provider to authenticate you need to configure the following
arguments:

* server_url - (Required) The server url of the ovc api to connect to
* client_id - (Required) The client_id of the itsyouonline user
* client_secret - (Required) The client_secret of the itsyouonline user

It is advisable to set these arguments as environment variables:

```
export OPENVCLOUD_SERVER_URL="server-url"
export ITSYOU_ONLINE_CLIENT_ID="your-client-id"
export ITSYOU_ONLINE_CLIENT_SECRET="your-client-secret"
```
This way the arguments must not be included in your terraform configuration file.

An example can be found under [examples/client-id-secret](./examples/client-id-secret)

### Authentication with a JWT

The following command is an example how to get a JWT using the `curl` command.
Providing `scope=offline_access` will return a JWT that is refreshable.

```sh
JWT=$(curl --silent -d 'grant_type=client_credentials&client_id='"$CLIENT_ID"'&client_secret='"$CLIENT_SECRET"'&response_type=id_token&scope=offline_access' https://itsyou.online/v1/oauth/access_token)
echo $JWT

By configuration the following parameters you will configure the ovc provider to authenticate
with the JWT
```
* server_url - (Required) The server url of the ovc api to connect to
* client_jwt - (Required) the JWT that is tied to your client_id and client_secret
```

An example can be found under [examples/client-jwt](./examples/client-jwt)
