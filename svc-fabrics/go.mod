module github.com/ODIM-Project/ODIM/svc-fabrics

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-dmtf v0.0.0-20200727133207-df3dfb728bd1
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2
	github.com/ODIM-Project/ODIM/svc-plugin-rest-client v0.0.0-20200727133207-df3dfb728bd1
	github.com/stretchr/testify v1.6.1
	gotest.tools v2.2.0+incompatible
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
replace github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2 => ../lib-utilities
replace github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
