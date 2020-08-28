module github.com/ODIM-Project/ODIM/svc-task

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-messagebus v0.0.0-20200727133207-df3dfb728bd1
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2
	github.com/RediSearch/redisearch-go v1.0.1
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f
	github.com/micro/go-micro v1.13.2
	github.com/satori/go.uuid v1.2.0
	github.com/satori/uuid v1.2.0
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
replace github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2 => ../lib-utilities
replace github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
