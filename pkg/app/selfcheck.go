package app

// RunSelfCheck 在启动时进行关键配置的轻量自检并输出日志提示
// 注意：这里只做非侵入式检查与提示，不阻断服务启动
func (app *App) RunSelfCheck() {
	if app == nil || app.Config == nil || app.Logger == nil {
		return
	}

	cfg := app.Config

	// PasswdKey 警告
	if cfg.PasswdKey == "" {
		app.Logger.Warn("security: PasswdKey is empty, please set PASSWD_KEY for production")
	}

	// CORS 组合检查
	if len(cfg.CORS.AllowedOrigins) > 0 {
		anyStar := false
		for _, o := range cfg.CORS.AllowedOrigins {
			if o == "*" {
				anyStar = true
				break
			}
		}
		if anyStar && cfg.CORS.AllowCredentials {
			app.Logger.Warn("security: CORS AllowCredentials=true with '*' origin; browsers may block; restrict origins in production")
		}
	}

	// Upload 目录提示
	if cfg.Upload.Dir == "" {
		app.Logger.Info("upload: Upload.Dir is empty; file upload static serving may be disabled")
	}

	// SQL 打印提示
	if cfg.MySQL.ShowSQL && cfg.LogConfig.Level != "debug" {
		app.Logger.Warn("security: MySQL.ShowSQL is enabled while log level is not 'debug'; may leak sensitive data")
	}

	// 采样与堆栈提示
	if !cfg.LogConfig.EnableSampling && cfg.LogConfig.Level != "debug" {
		app.Logger.Info("log: sampling disabled on non-debug level; consider enabling LOG_SAMPLING_ENABLE to reduce volume")
	}
}
