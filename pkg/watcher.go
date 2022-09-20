package pkg

import (
	v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func WatchingSources() *source.Kind {
	return &source.Kind{Type: &v1.Service{}}
}
