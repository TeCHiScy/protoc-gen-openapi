# protoc-gen-openapi

This directory contains a protoc plugin that generates an
OpenAPI description for a REST API that corresponds to a
Protocol Buffer service.

Installation:

    go install github.com/google/gnostic/cmd/protoc-gen-openapi

Usage:

	protoc sample.proto -I=. --openapi_out=.

This runs the plugin for a file named `sample.proto` which 
refers to additional .proto files in the same directory as
`sample.proto`. Output is written to the current directory.

## options

1. `version`: version number text, e.g. 1.2.3
   - **default**: `0.0.1`
2. `title`: name of the API
   - **default**: empty string or service name if there is only one service
3. `description`: description of the API
   - **default**: empty string or service description if there is only one service
4. `naming`: naming convention. Use "proto" for passing names directly from the proto files
   - **default**: `json`
   - `json`: will turn field `updated_at` to `updatedAt`
   - `proto`: keep field `updated_at` as it is
5. `fq_schema_naming`: schema naming convention. If "true", generates fully-qualified schema names by prefixing them with the proto message package name
   - **default**: false
   - `false`: keep message `Book` as it is
   - `true`: turn message `Book` to `google.example.library.v1.Book`, it is useful when there are same named message in different package
6. `enum_type`: type for enum serialization. Use "string" for string-based serialization
   - **default**: `integer`
   - `integer`: setting type to `integer`
      ```yaml
      schema:
        type: integer
        format: enum
      ```
   - `string`: setting type to `string`, and list available values in `enum`
      ```yaml
      schema:
        enum:
          - UNKNOWN_KIND
          - KIND_1
          - KIND_2
        type: string
        format: enum
      ```
7. `depth`: depth of recursion for circular messages
   - **default**: 2, this depth only used in query parameters, usually 2 is enough
8. `default_response`: add default response. If "true", automatically adds a default response to operations which use the google.rpc.Status message.
   Useful if you use envoy or grpc-gateway to transcode as they use this type for their default error responses.
   - **default**: true, this option will add this default response for each method as following:
      ```yaml
      default:
        description: Default error response
          content:
            application/json:
              schema:
                $ref: '#/components/schemas/google.rpc.Status'
      ```

# protoc-gen-openapi

Contains a protoc plugin that generates openapi v3 documents

**Forked from [github.com/google/gnostic/cmd/protoc-gen-openapi](https://github.com/google/gnostic/tree/main/cmd/protoc-gen-openapi)**

Installation:

    go install github.com/kollalabs/protoc-gen-openapi@latest

Usage:

    protoc sample.proto -I. --openapi_out=version=1.2.3:.

## Testing

To see output during tests use `log.Print*`

```
go clean -testcache && go test
```

## Added Features
We have added some features that the Gnostic team most likely doesn't want to add :-)
Some are fairly Kolla specific, sorry. We try to hide Kolla specific functionality
in a way that won't trip anyone up.

- [protoc-gen-openapi](#protoc-gen-openapi)
  - [options](#options)
- [protoc-gen-openapi](#protoc-gen-openapi-1)
  - [Testing](#testing)
  - [Added Features](#added-features)
    - [Better Enum Support](#better-enum-support)
    - [Summary Field](#summary-field)
    - [Validation](#validation)
      - [Example](#example)
      - [Supported Validators](#supported-validators)
    - [Google Field Behavior Annotations](#google-field-behavior-annotations)
    - [OAS3 header support](#oas3-header-support)
- [https://github.com/kollalabs/protoc-gen-openapi?tab=readme-ov-file](#httpsgithubcomkollalabsprotoc-gen-openapitabreadme-ov-file)

### Better Enum Support
Enums work better by using string values of proto enums instead of ints.

### Summary Field

Sometimes you want more control over certain properties in the OpenAPI manifest. In our
case we wanted to use the `summary` property on routes to look nice for generating
documentation from the OpenAPI manifest. Normally the summary comes simply from the
name of the route. We added a feature that parses the comment over the proto service
method and looks for a pipe character ("`|`") and if it sees it, it will take anything to
the left of it and put it in the `summary` field, and anything to the right of it will
be the `description`. If no pipe is found it puts the whole comment in the description
like normal. From `/examples/tests/summary/message.proto`:

```proto
service Messaging {
    // Update Message Summary | This function updates a message.
    rpc UpdateMessage(Message) returns(Message) {
        option(google.api.http) = {
            patch: "/v1/messages/{message_id}"
            body: "text"
        };
    }
}
```

It generates the following OpenAPI:

```yaml
#...
paths:
    /v1/messages/{message_id}:
        patch:
            tags:
                - Messaging
            summary: Update Message Summary # Look at this beautiful summary...
            description: This function updates a message.
#...
```

### Validation

We added partial support for `protoc-gen-validate` annotations

OpenAPI spec allows for a small handful of input validation configurations.
Proto has an awesome plugin called `protoc-gen-validate` for generating validation code in
Go, Java, C++, etc. We took those same annotations and added support in this project
for them.

Usage: add `validate=true` to protoc command.

`protoc sample.proto -I. --openapi_out=version=1.2.3,validate=true:.`

#### Example

```proto
message Message {
    string message_id = 1;
    string text = 2 [(validate.rules)= {
        string: {
            uri:true,
            max_len:45,
            min_len:1
        }
    }];
    int64 mynum = 3 [(validate.rules).int64 = {gte:1, lte:30}];
}

```

outputs:

```yaml
components:
    schemas:
        Message:
            properties:
                message_id:
                    type: string
                text:
                    maxLength: 45
                    minLength: 1
                    type: string
                    format: uri
                mynum:
                    maximum: !!float 30
                    minimum: !!float 1
                    type: integer
                    format: int64

```

#### Supported Validators

String
- uri
- uuid
- email
- ipv4
- ipv6
- max_len
- min_len

Int32
- gte
- lte

Int64
- gte
- lte

Adding more can easily be done in the function `addValidationRules` in `/generator/openapi-v3.yaml`

### Google Field Behavior Annotations

* `(google.api.field_behavior) = REQUIRED` will add the field to the required list in the openAPI schema
* `(google.api.field_behavior) = OUTPUT_ONLY` will add the `readOnly` property to the field
* `(google.api.field_behavior) = INPUT_ONLY` will add the `writeOnly` property to the field
* TODO: `(google.api.field_behavior) = IMMUTABLE` will add the `x-createOnly` property to the field (not supported by openapi yet)

### OAS3 header support


# https://github.com/kollalabs/protoc-gen-openapi?tab=readme-ov-file