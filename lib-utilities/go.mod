module github.com/ODIM-Project/ODIM/lib-utilities

go 1.13

require (
	github.com/Joker/hpp v1.0.0 // indirect
	github.com/ODIM-Project/ODIM/lib-persistence-manager v0.0.0-20200809094501-2cf1892f9897
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f
	github.com/kataras/iris/v12 v12.1.9-0.20200616210209-a85c83b70ad0
	github.com/micro/go-micro v1.13.2
//	google.golang.org/grpc/examples v0.0.0-20200807164945-d3e3e7a46f57 // indirect
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
replace github.com/ODIM-Project/ODIM/lib-persistence-manager v0.0.0-20200809094501-2cf1892f9897 => ../lib-persistence-manager
