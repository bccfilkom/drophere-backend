package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/99designs/gqlgen/handler"
	drophere_go "github.com/bccfilkom/drophere-go"
	"github.com/bccfilkom/drophere-go/domain/link"
	"github.com/bccfilkom/drophere-go/domain/user"
	"github.com/bccfilkom/drophere-go/infrastructure/auth"
	"github.com/bccfilkom/drophere-go/infrastructure/database/mysql"
	"github.com/bccfilkom/drophere-go/infrastructure/hasher"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
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
		panic(fmt.Errorf("config: %s", err))
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

	// initialize repositories
	userRepo := mysql.NewUserRepository(db)
	linkRepo := mysql.NewLinkRepository(db)

	// initialize infrastructures
	authenticator := auth.NewJWT(
		viper.GetString("jwt.secret"),
		time.Duration(viper.GetInt("jwt.duration"))*time.Hour,
		viper.GetString("jwt.signingAlgorithm"),
		userRepo,
	)
	bcryptHasher := hasher.NewBcryptHasher()

	// initialize services
	userSvc := user.NewService(userRepo, authenticator, bcryptHasher)
	linkSvc := link.NewService(linkRepo, bcryptHasher)

	resolver := drophere_go.NewResolver(userSvc, authenticator, linkSvc)

	// setup router
	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(authenticator.Middleware())
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:8080"},
		AllowCredentials: true,
		Debug:            false,
	}).Handler)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(drophere_go.NewExecutableSchema(drophere_go.Config{Resolvers: resolver})))
	router.Post("/uploadfile", fileUploadHandler(userSvc, linkSvc))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
