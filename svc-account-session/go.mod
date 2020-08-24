module github.com/ODIM-Project/ODIM/svc-account-session

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2
	github.com/satori/go.uuid v1.2.0
	github.com/stretchr/testify v1.6.1
	golang.org/x/crypto v0.0.0-20200510223506-06a226fb4e37
	gopkg.in/go-playground/validator.v9 v9.30.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
replace github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2 => ../lib-utilities
replace github.com/ODIM-Project/ODIM/lib-persistence-manager => ../lib-persistence-manager
