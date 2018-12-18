package server

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	pb "github.com/scottshotgg/storage/protobufs"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
)

type errorBody struct {
	Err string `json:"error,omitempty"`
}

// CustomHTTPError is used for implementing a custom error return
func CustomHTTPError(ctx context.Context, _ *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, _ *http.Request, err error) {
	const (
		fallback = `{"error": "failed to marshal error message"}`
	)

	w.Header().Set("Content-type", marshaler.ContentType())
	w.WriteHeader(runtime.HTTPStatusFromCode(grpc.Code(err)))

	var jErr = json.NewEncoder(w).Encode(errorBody{
		Err: grpc.ErrorDesc(err),
	})

	if jErr != nil {
		var _, err = w.Write([]byte(fallback))
		if err != nil {
			// log
		}
	}
}

// Move these errors to be declared above

// RunREST Start rest server for geosearch
func RunREST() error {
	var servicename = viper.GetString("servicename")
	if servicename == "" {
		return errors.New("You must supply a valid servicename for logging using the `servicename` flag")
	}

	var serverIP = viper.GetString("server-ip")
	if serverIP == "" {
		return errors.New("You must supply a valid server-ip using the 'server-ip' argument")
	}

	var restPort = ":" + viper.GetString("rest-port")
	if restPort == ":" {
		return errors.New("You must supply a valid port using the 'rest-port' argument")
	}

	var rpcAddr = viper.GetString("rpc-addr")
	if rpcAddr == "" {
		return errors.New("You must supply a valid address using the 'rpc-addr' argument")
	}

	var rpcPort = viper.GetString("rpc-port")
	if rpcPort == "" {
		return errors.New("You must supply a valid port using the 'rpc-port' argument")
	}

	var (
		mux = runtime.NewServeMux()
		err = pb.RegisterStorageHandlerFromEndpoint(
			context.Background(),
			mux,
			rpcAddr+":"+rpcPort,
			[]grpc.DialOption{
				grpc.WithInsecure(),
				grpc.WithBlock(),
			})
	)

	if err != nil {
		return errors.Wrap(err, "failed to start HTTP server: %v")
	}

	// TODO: Make middleware stuff here

	var scheme = viper.GetString("scheme")

	if scheme == "https" {
		var crt = viper.GetString("tls-certificate")
		if crt == "" {
			return errors.New("You must supply a valid crt using the 'tls-certificate' argument")
		}

		var key = viper.GetString("tls-key")
		if key == "" {
			return errors.New("You must supply a valid key using the 'tls-key' argument")
		}

		return http.ListenAndServeTLS(restPort, crt, key, mux)
	}

	return http.ListenAndServe(restPort, mux)
}
