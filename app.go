package main

import (
	"amumu-wms-backend/controller/dispatch"
	"amumu-wms-backend/controller/goods"
	"amumu-wms-backend/controller/inventory"
	"amumu-wms-backend/controller/logs"
	"amumu-wms-backend/controller/users"
	"amumu-wms-backend/controller/warehouse"
	"amumu-wms-backend/models"
	"amumu-wms-backend/utils"
	"flag"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/pelletier/go-toml/v2"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Backend *BackendConfig `toml:"backend"`
	DB      *DBConfig      `toml:"db"`
	JWT     *JWTConfig     `toml:"jwt"`
	Crypto  *CryptoConfig  `toml:"crypto"`
	Debug   *DebugConfig   `toml:"debug"`
}

type BackendConfig struct {
	ListenAddr *string `toml:"listen-addr"`
	Prefix     *string `toml:"prefix"`
}

type DBConfig struct {
	Username *string `toml:"username"`
	Password *string `toml:"password"`
	Host     *string `toml:"host"`
	Port     *int    `toml:"port"`
	DB       *string `toml:"db"`
}

type JWTConfig struct {
	SigningKey *string `toml:"signing_key"`
}

type CryptoConfig struct {
	Salt *string `toml:"salt"`
}

type DebugConfig struct {
	Enabled *bool `toml:"enabled"`
}

var (
	config     Config
	configPath = flag.String("c", "config.toml", "config file")
)

func loadConfig() {
	configBytes, err := os.ReadFile(*configPath)
	if err != nil {
		log.Fatalf("load config error: %s.\n", err.Error())
	}

	err = toml.Unmarshal(configBytes, &config)
	if err != nil {
		log.Fatalf("parse config error: %s.\n", err.Error())
	}

	if config.Backend == nil {
		config.Backend = &BackendConfig{}
	}

	if config.DB == nil {
		config.DB = &DBConfig{}
	}

	if config.JWT == nil {
		config.JWT = &JWTConfig{}
	}

	if config.Crypto == nil {
		config.Crypto = &CryptoConfig{}
	}

	if config.Debug == nil {
		config.Debug = &DebugConfig{}
	}

	if config.Backend.ListenAddr == nil {
		listenAddr := ":8080"
		config.Backend.ListenAddr = &listenAddr
	}

	if config.Backend.Prefix == nil {
		prefix := "/"
		config.Backend.Prefix = &prefix
	}

	if config.DB.Username == nil || config.DB.Password == nil || config.DB.Host == nil || config.DB.Port == nil || config.DB.DB == nil {
		log.Fatalf("username, password, host, port or db in database config should not be null.\n")
	}

	if config.JWT.SigningKey == nil || *config.JWT.SigningKey == "" {
		log.Fatalf("jwt signing key should not be null.\n")
	}

	if config.Crypto.Salt == nil || *config.Crypto.Salt == "" {
		log.Fatalf("crypto salt should not be null.\n")
	}

	if config.Debug.Enabled == nil {
		enabled := false
		config.Debug.Enabled = &enabled
	}
}

func loadDB() {
	models.Init(*config.DB.Username, *config.DB.Password, *config.DB.Host, *config.DB.Port, *config.DB.DB)
}

func startBackend() {
	if !*config.Debug.Enabled {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()
	router.Use(cors.Default())
	router.Use(utils.GinRequestLogMiddleware)
	router.Use(utils.GinBodyLogMiddleware)

	prefix := *config.Backend.Prefix
	prefix += "api"
	group := router.Group(prefix)

	users.Init(group)
	warehouse.Init(group)
	goods.Init(group)
	inventory.Init(group)
	dispatch.Init(group)
	logs.Init(group)

	s := &http.Server{
		Addr:           *config.Backend.ListenAddr,
		Handler:        router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	s.ListenAndServe()
}

func main() {
	flag.Parse()

	if _, err := os.Stat(*configPath); err != nil {
		log.Fatalf("config file not exists.\n")
	}

	loadConfig()
	loadDB()
	utils.InitJWT(*config.JWT.SigningKey)
	utils.InitCrypto(*config.Crypto.Salt)
	utils.InitLog(*config.Debug.Enabled)

	log.Printf("Starting amumu-wms-backend at %s\n", *config.Backend.ListenAddr)

	startBackend()
}
