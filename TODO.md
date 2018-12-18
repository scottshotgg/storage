# TODO

- Move grpc-gateway stuff except for protobufs folder to a new repo since the server has to vendor the protobufs anyways
- Implement proto.Marshaler and json.Marshaler
- Change MarshalBinary to use proto.Marshaler
- Make an object.Encode/DecodeValue
- Implement a GobEncoder/GobDecoder
- Experiment with websockets and grpc streams