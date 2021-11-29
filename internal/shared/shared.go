package shared

import (
	"time"

	"github.com/yndd/ndd-runtime/pkg/logging"
	"github.com/yndd/ndd-yang/pkg/yentry"
)

type NddControllerOptions struct {
	Logger           logging.Logger
	Poll             time.Duration
	Namespace        string
	Yentry           *yentry.Entry
	GrpcQueryAddress string
}
