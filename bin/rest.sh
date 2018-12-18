ENV=dev \
go build \
&& \
./storage rest \
--servicename=storage \
--server-ip="localhost" \
--rest-port="6000" \
--rpc-port="6001"