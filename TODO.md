# TODO

- Need to fix `Keys` - think about it for a bit
- Need to fix the Item implementation, maybe just use the protbufs for everything?
- fix all the implementations
- fix all the tests
- decide where `multi` is going and how it should look
- Collapse all of the GetX functions into just a single Get
- Collapse all of the SetX functions into just a single Set
- Collapse all of the DeleteX functions into just a single Delete
- Collapse all of the IteratorX functions into just a single Iterator
- Could collapse Audit, QuickSync, and Sync into single function ...
- Experiment with websockets and grpc streams

- Using the protobufs as the end data structure is a bit weird because it artificially forces some attributes and more or less un-controllable features onto the rpc protos