module attestor/host

go 1.25.7

require (
	attestor/core v0.0.0
	attestor/grpc v0.0.0
	attestor/types v0.0.0
	google.golang.org/grpc v1.82.1
)

require (
	golang.org/x/net v0.53.0 // indirect
	golang.org/x/sys v0.43.0 // indirect
	golang.org/x/text v0.36.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20260414002931-afd174a4e478 // indirect
	google.golang.org/protobuf v1.36.11 // indirect
)

replace attestor/core => ../core

replace attestor/grpc => ../grpc

replace attestor/types => ../types
