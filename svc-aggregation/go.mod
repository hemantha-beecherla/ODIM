module github.com/ODIM-Project/ODIM/svc-aggregation

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-dmtf v0.0.0-20200727115727-33d557ff397c
	github.com/ODIM-Project/ODIM/lib-messagebus v0.0.0-20200727103018-252e26a63065
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2
	github.com/ODIM-Project/ODIM/svc-plugin-rest-client v0.0.0-20200727110501-4599893e44fd
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.6.1
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0

replace github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2 => ../lib-utilities
replace github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager

