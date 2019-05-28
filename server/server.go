package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	drophere_go "github.com/bccfilkom/drophere-go"
	"github.com/bccfilkom/drophere-go/domain/user"
	"github.com/bccfilkom/drophere-go/infrastructure/auth"
	"github.com/bccfilkom/drophere-go/infrastructure/database/mysql"

	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

const defaultPort = "8080"

func main() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error config file: %s", err))
	}

	viper.SetEnvPrefix("DROPHERE")
	viper.AutomaticEnv()

	port := viper.GetString("PORT")
	if port == "" {
		port = defaultPort
	}

	// setup
	db, err := mysql.New(viper.GetString("db.dsn"))
	if err != nil {
		panic(err)
	}
	authenticator := auth.NewJWT(
		viper.GetString("jwt.secret"),
		time.Duration(viper.GetInt("jwt.duration")),
		viper.GetString("jwt.signingAlgorithm"),
	)
	userRepo := mysql.NewUserRepository(db)
	userSvc := user.NewService(userRepo, authenticator)

	resolver := drophere_go.NewResolver(userSvc)

	http.Handle("/", handler.Playground("GraphQL playground", "/query"))
	http.Handle("/query", handler.GraphQL(drophere_go.NewExecutableSchema(drophere_go.Config{Resolvers: resolver})))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
