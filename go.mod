module github.com/thepetk/snippetbox

go 1.19

replace github.com/thepetk/snippetbox/cmd/web/config => ./cmd/web/config.go

require (
	github.com/gofrs/uuid v4.4.0+incompatible
	github.com/lib/pq v1.10.9
)
