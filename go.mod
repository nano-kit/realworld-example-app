module realworld-example-app

go 1.13

// This can be removed once etcd becomes go gettable, version 3.4 and 3.5 is not,
// see https://github.com/etcd-io/etcd/issues/11154 and https://github.com/etcd-io/etcd/issues/11931.
replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

// Use a production ready go-micro/v2 stable version maintained by nano-kit.
replace github.com/micro/go-micro/v2 => ../go-micro

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gobwas/ws v1.0.4
	github.com/golang/protobuf v1.4.3
	github.com/google/uuid v1.2.0
	github.com/gorilla/mux v1.7.3
	github.com/json-iterator/go v1.1.10
	github.com/micro/go-micro/v2 v2.0.0-00010101000000-000000000000
)
