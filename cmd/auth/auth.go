package main

import (
	"auth-service/internal/core"
	"auth-service/internal/middlewares"
	"auth-service/internal/services/auth"
	"auth-service/sdk/auth/proto"
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.opencensus.io/plugin/ocgrpc"
	"go.opencensus.io/trace"
	"google.golang.org/api/pagespeedonline/v5"
	_ "google.golang.org/api/pagespeedonline/v5"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/encoding/protojson"
	"os"
	"strconv"

	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
)

var (
	BuildAt         string
	Commit          string
	Version         string
	jwtSK           string
	tokenExp        int64
	refreshTokenExp int64
)

func init() {
	tokenExpString := os.Getenv("TOKEN_EXP")
	refreshTokenExpString := os.Getenv("REFRESH_TOKEN_EXP")
	jwtSK = os.Getenv("JWT_SECRET")

	if tokenExpString == "" || refreshTokenExpString == "" || jwtSK == "" {
		log.Fatalln("not found TOKEN_EXP || REFRESH_TOKEN_EXP || JWT_SECRET in env")
	}
	if tokenExpString == "" || refreshTokenExpString == "" {
		log.Fatalln("not found TOKEN_EXP || REFRESH_TOKEN_EXP in env")
	}

	tokenExpInt, err := strconv.Atoi(tokenExpString)
	if err != nil {
		log.Fatalln(err)
	}
	tokenExp = int64(tokenExpInt)
	refreshTokenExpInt, err := strconv.Atoi(refreshTokenExpString)
	if err != nil {
		log.Fatalln(err)
	}
	refreshTokenExp = int64(refreshTokenExpInt)
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	pagespeedonlineService, err := pagespeedonline.NewService(ctx)

	service := auth.NewService(Commit, BuildAt, Version, jwtSK, tokenExp, refreshTokenExp)

	//db, err := mysql.NewAuthStore(true)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//service.SetStore(db)

	listener, err := net.Listen("tcp", "0.0.0.0:3009")
	if err != nil {
		log.Fatalln("Error create listener:", err)
	}

	defer cancel()

	server := grpc.NewServer(grpc.ChainUnaryInterceptor(core.ServiceOptions),
		grpc.StatsHandler(&ocgrpc.ServerHandler{}))
	proto.RegisterAuthServer(server, service)

	go func() {
		// GRPC server
		if err := server.Serve(listener); err != nil {
			log.Fatalln("Error start service", err)
		}
	}()

	go func() {
		mux := runtime.NewServeMux(
			runtime.WithErrorHandler(func(ctx context.Context, serveMux *runtime.ServeMux, marshaler runtime.Marshaler, writer http.ResponseWriter, request *http.Request, err error) {

				newError := runtime.HTTPStatusError{
					HTTPStatus: 200,
					Err:        err,
				}
				// using default handler to do the rest of heavy lifting of marshaling error and adding headers
				runtime.DefaultHTTPErrorHandler(ctx, serveMux, marshaler, writer, request, &newError)
			}),
			runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.HTTPBodyMarshaler{
				Marshaler: &runtime.JSONPb{
					MarshalOptions: protojson.MarshalOptions{
						UseProtoNames:   false,
						EmitUnpopulated: true,
					},
					UnmarshalOptions: protojson.UnmarshalOptions{
						DiscardUnknown: true,
					},
				},
			}))
		connectOpts := []grpc.DialOption{
			grpc.WithInsecure(),
		}

		if err := proto.RegisterAuthHandlerFromEndpoint(ctx, mux, "0.0.0.0:3009", connectOpts); err != nil {
			log.Fatalln("Could not connect to GRPC server on port 3009", err)
		}

		if err := http.ListenAndServe("0.0.0.0:3030", middlewares.AllowCORS(mux)); err != nil {
			log.Println("Error listen http:", err)
		}

	}()
	trace.RegisterExporter(core.NewExporter("auth-service"))
	log.Println("Success started auth-service")

	//log.Println(pagespeedonlineService.UserAgent)
	log.Println(pagespeedonlineService)

	<-ctx.Done()
}
