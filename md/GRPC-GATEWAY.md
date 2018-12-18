# GRPC-GATEWAY

- grpc-gateway is an implementation itself, however - like the `store` implementaton - it takes `storage` object and uses that under the hood

- this provides both `HTTP` and `RPC` access points, but also has an abstraction over the grpc stuff so that it can satisfy the interface

- doing so allows a `client` to use the rpc server normally through the `storage` interface while the `server` can serve a storage object through `HTTP` and `RPC` protocols