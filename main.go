package main

import (
	"context"
	"fmt"
	"log"

	"github.com/ilyakaznacheev/cleanenv"
	echo "github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/souviks72/tronicscorp/config"
	"github.com/souviks72/tronicscorp/handlers"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var (
	c   *mongo.Client
	db  *mongo.Database
	col *mongo.Collection
	cfg config.Properties
)

func init() {
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		log.Fatalf("Config not read %+v\n", err)
	}
	connectURI := fmt.Sprintf("mongodb://%s:%s", cfg.DBHost, cfg.DBPort)
	c, err := mongo.Connect(context.Background(), options.Client().ApplyURI(connectURI))
	if err != nil {
		log.Fatalf("Unable to connect to mongodb %+v\n", err)
	}

	db = c.Database(cfg.DBName)
	col = db.Collection(cfg.CollectionName)
}

func main() {
	e := echo.New()
	e.Pre(middleware.RemoveTrailingSlash())
	h := handlers.ProductHandler{Col: col}
	e.POST("/products", h.CreateProducts, middleware.BodyLimit("1M"))
	e.Logger.Infof("Listening n %s:%s", cfg.Host, cfg.Port)
	e.Logger.Fatal(e.Start(fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)))
}
