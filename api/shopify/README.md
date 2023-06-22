# gqlgenc README

gqlgenc is a tool used to generate Go clients and models for GraphQL. This tool takes a configuration file in YAML format and GraphQL query files as inputs, and it generates typed Go code that can be used to interact with a GraphQL API.

## Configuration

The configuration file (e.g., `.gqlgenc.yml`) contains the necessary information for gqlgenc to generate the Go code. Here is an example configuration file:

```yaml
model:
  # The filename for the generated model code.
  filename: ./gen/models_gen.go
client:
  # The filename for the generated client code.
  filename: ./gen/client_gen.go
endpoint:
  # The GraphQL API endpoint.
  url: https://[YOUR_SHOP_NAME].myshopify.com/admin/api/2023-01/graphql.json
  headers:
    # The access token for the API endpoint.
    X-Shopify-Access-Token: "[YOUR_ACCESS_TOKEN]"
query:
  # The path for the GraphQL query files.
  - "./query/*.graphql"
generate:
  # Generate the V2 client.
  clientV2: true
  # The name of the client interface.
  clientInterfaceName: "ShopifyGraphQLClient"
```

### Fields

#### `model.filename`

The filename for the generated model code.

#### `client.filename`

The filename for the generated client code.

#### `endpoint.url`

The GraphQL API endpoint.

#### `endpoint.headers`

The headers for the API endpoint.

#### `query`

The path for the GraphQL query files.

#### `generate.clientV2`

A boolean value that determines whether to generate the V2 client.

#### `generate.clientInterfaceName`

The name of the client interface.

## Generating Go code using gqlgenc

To generate the Go code using gqlgenc, run the following command in the same directory as the .gqlgenc.yml file:

```console
gqlgenc
```

This command will use the configuration file `.gqlgenc.yml` to generate the Go code. The generated code will be saved in the files specified in the `model.filename` and `client.filename` fields of the configuration file.
