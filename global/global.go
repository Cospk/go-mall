package global

import (
	"github.com/Cospk/go-mall/config"
	"go.uber.org/zap"
)

var (
	Config config.Config
	Logger *zap.Logger
)
