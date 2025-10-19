package main

import (
	"github.com/mSulimenko/dev-blog-platform/internal/articles/config"
	"github.com/mSulimenko/dev-blog-platform/internal/shared/logger"
)

func main() {

	cfg := config.Load()

	log := logger.New(cfg.Env)
	defer log.Sync()

}
