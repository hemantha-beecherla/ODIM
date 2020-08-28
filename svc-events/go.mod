module github.com/ODIM-Project/ODIM/svc-events

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-messagebus v0.0.0-20200727133207-df3dfb728bd1
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2
	github.com/ODIM-Project/ODIM/svc-plugin-rest-client v0.0.0-20200727133207-df3dfb728bd1
	github.com/google/uuid v1.1.2-0.20200519141726-cb32006e483f
	github.com/stretchr/testify v1.6.1
	gopkg.in/go-playground/validator.v9 v9.30.0
	gotest.tools v2.2.0+incompatible
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
replace github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2 => ../lib-utilities
replace github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
