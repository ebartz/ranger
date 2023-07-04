package clusterrouter

import (
	"net/http"

	v3 "github.com/ranger/ranger/pkg/generated/norman/management.cattle.io/v3"
)

type server interface {
	Close()
	Handler() http.Handler
	Cluster() *v3.Cluster
}
