package main

import (
        "context"
        "fmt"
        "github.com/gin-gonic/gin"
        "github.com/redis/go-redis/v9"
        "go.uber.org/zap"
        "go.uber.org/zap/zapcore"
        "gopkg.in/yaml.v3"
        "net/http"
        "os"
        "os/signal"
        "syscall"
        "time"
)

// ========= 配置结构体 & 加载 =========
type Config struct {
        Server struct {
                Port         int `yaml:"port"`
                ReadTimeout  int `yaml:"read_timeout"`
                WriteTimeout int `yaml:"write_timeout"`
                IdleTimeout  int `yaml:"idle_timeout"`
        } `yaml:"server"`
        Redis struct {
                Addr         string `yaml:"addr"`
                Password     string `yaml:"password"`
                DB           int    `yaml:"db"`
                PoolSize     int    `yaml:"pool_size"`
                MinIdleConns int    `yaml:"min_idle_conns"`
                PoolTimeout  int    `yaml:"pool_timeout"`
                DialTimeout  int    `yaml:"dial_timeout"`
                ReadTimeout  int    `yaml:"read_timeout"`
                WriteTimeout int    `yaml:"write_timeout"`
        } `yaml:"redis"`
}

func LoadConfig(path string) (*Config, error) {
        f, err := os.Open(path)
        if err != nil {
                return nil, err
        }
        defer f.Close()
        decoder := yaml.NewDecoder(f)
        var c Config
        if err := decoder.Decode(&c); err != nil {
                return nil, err
        }
        return &c, nil
}

// ========= 全局 =========
var (
        rdb  *redis.Client
        ctx  = context.Background()
        conf *Config
        zlog *zap.Logger
)

// ========= gin zap 中间件 =========
func GinZapLogger(logger *zap.Logger) gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        path := c.Request.URL.Path
        raw := c.Request.URL.RawQuery
        c.Next()
        latency := time.Since(start)

        logger.Info("access",
            zap.String("host", c.Request.Host),
            zap.String("client_ip", c.ClientIP()),
            zap.String("method", c.Request.Method),
            zap.String("path", path),
            zap.String("query", raw),
            zap.Int("status", c.Writer.Status()),
            zap.Duration("latency", latency),
            zap.String("user_agent", c.Request.UserAgent()),
            zap.String("referer", c.Request.Referer()),
            zap.Int("body_size", c.Writer.Size()),
            zap.String("error", c.Errors.ByType(gin.ErrorTypePrivate).String()),
        )
    }
}
// ========= 路由处理 =========
func CheckHandler(c *gin.Context) {
        domain := c.Query("domain")
        if domain == "" {
                c.String(http.StatusBadRequest, "no domain")
                return
        }
        ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
        defer cancel()
        exists, err := rdb.HExists(ctxTimeout, "domain_map", domain).Result()
        if err != nil {
                zlog.Error("redis error", zap.Error(err))
                c.String(http.StatusInternalServerError, "service error")
                return
        }
        if exists {
                c.String(http.StatusOK, "yes")
        } else {
                c.String(http.StatusForbidden, "no")
        }
}

func HealthHandler(c *gin.Context) {
        if err := rdb.Ping(ctx).Err(); err != nil {
                zlog.Error("redis down", zap.Error(err))
                c.String(http.StatusServiceUnavailable, "redis down")
                return
        }
        c.String(http.StatusOK, "ok")
}

func AnyHandler(c *gin.Context) {
        host := c.Request.Host
        ctxTimeout, cancel := context.WithTimeout(ctx, 1*time.Second)
        defer cancel()
        url, err := rdb.HGet(ctxTimeout, "domain_map", host).Result()
        if err == redis.Nil {
                c.String(http.StatusNotFound, "<h1>Default for %s</h1>", host)
                return
        } else if err != nil {
                zlog.Error("redis error", zap.Error(err))
                c.String(http.StatusInternalServerError, "service error")
                return
        }
        html := fmt.Sprintf(`<meta http-equiv="refresh" content="0;url=%s">`, url)
        c.Data(http.StatusOK, "text/html", []byte(html))
}

// ========= zap logger =========
func InitZapLogger() *zap.Logger {
        // Console Encoder
        consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

        // JSON Encoder（写入 access.log 文件）
        jsonEncoderCfg := zap.NewProductionEncoderConfig()
        jsonEncoderCfg.TimeKey = "ts"
        jsonEncoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
        jsonEncoder := zapcore.NewJSONEncoder(jsonEncoderCfg)

        // 日志文件
        logFile, err := os.OpenFile("access.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
        if err != nil {
                panic(fmt.Sprintf("cannot open access.log: %v", err))
        }
        fileWriter := zapcore.AddSync(logFile)
        consoleWriter := zapcore.AddSync(os.Stdout)

        // 多个输出
        core := zapcore.NewTee(
                zapcore.NewCore(consoleEncoder, consoleWriter, zapcore.DebugLevel),
                zapcore.NewCore(jsonEncoder, fileWriter, zapcore.InfoLevel),
        )
        logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
        return logger
}

// ========= 主入口 =========
func main() {
        var err error
        conf, err = LoadConfig("config.yaml")
        if err != nil {
                fmt.Fprintf(os.Stderr, "load config error: %v\n", err)
                os.Exit(1)
        }
        zlog = InitZapLogger()
        defer zlog.Sync()

        // 初始化 redis
        rdb = redis.NewClient(&redis.Options{
                Addr:         conf.Redis.Addr,
                Password:     conf.Redis.Password,
                DB:           conf.Redis.DB,
                PoolSize:     conf.Redis.PoolSize,
                MinIdleConns: conf.Redis.MinIdleConns,
                PoolTimeout:  time.Duration(conf.Redis.PoolTimeout) * time.Second,
                DialTimeout:  time.Duration(conf.Redis.DialTimeout) * time.Second,
                ReadTimeout:  time.Duration(conf.Redis.ReadTimeout) * time.Second,
                WriteTimeout: time.Duration(conf.Redis.WriteTimeout) * time.Second,
        })

        gin.SetMode(gin.ReleaseMode)
        router := gin.New()
        router.Use(GinZapLogger(zlog), gin.Recovery())

        router.GET("/check", CheckHandler)
        router.GET("/health", HealthHandler)
        router.NoRoute(AnyHandler)

        srv := &http.Server{
                Addr:         fmt.Sprintf(":%d", conf.Server.Port),
                Handler:      router,
                ReadTimeout:  time.Duration(conf.Server.ReadTimeout) * time.Second,
                WriteTimeout: time.Duration(conf.Server.WriteTimeout) * time.Second,
                IdleTimeout:  time.Duration(conf.Server.IdleTimeout) * time.Second,
        }

        go func() {
                zlog.Info("server starting", zap.Int("port", conf.Server.Port))
                if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
                        zlog.Fatal("server error", zap.Error(err))
                }
        }()

        // 优雅关机
        quit := make(chan os.Signal, 1)
        signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
        <-quit
        ctxTimeout, cancel := context.WithTimeout(context.Background(), 10*time.Second)
        defer cancel()
        zlog.Info("shutting down...")
        if err := srv.Shutdown(ctxTimeout); err != nil {
                zlog.Fatal("server forced shutdown", zap.Error(err))
        }
        zlog.Info("server exiting")
}
