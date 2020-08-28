module github.com/ODIM-Project/ODIM/lib-persistence-manager

go 1.13

require (
	github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2
	github.com/go-redis/redis/v8 v8.0.0-beta.7.0.20200807064721-0ddc3abd36b0
	github.com/gomodule/redigo v2.0.0+incompatible
)

replace github.com/ODIM-Project/ODIM/lib-utilities v0.0.0-20200809093149-80a1a9247bb2 => ../lib-utilities
