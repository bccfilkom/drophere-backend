package main

import (
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	htmlTemplate "html/template"
	textTemplate "text/template"

	drophere_go "github.com/bccfilkom/drophere-go"
	"github.com/bccfilkom/drophere-go/domain"
	"github.com/bccfilkom/drophere-go/domain/link"
	"github.com/bccfilkom/drophere-go/domain/user"
	"github.com/bccfilkom/drophere-go/infrastructure/auth"
	"github.com/bccfilkom/drophere-go/infrastructure/database/mysql"
	"github.com/bccfilkom/drophere-go/infrastructure/hasher"
	"github.com/bccfilkom/drophere-go/infrastructure/mailer"
	"github.com/bccfilkom/drophere-go/infrastructure/storageprovider"
	"github.com/bccfilkom/drophere-go/infrastructure/stringgenerator"

	"github.com/99designs/gqlgen/handler"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/cors"
	"github.com/spf13/viper"

	_ "github.com/go-sql-driver/mysql"
)

const defaultPort = "8080"

var debug bool

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

	// set debug mode
	debug = viper.GetBool("app.debug")

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
	userStorageCredRepo := mysql.NewUserStorageCredentialRepository(db)

	// initialize infrastructures
	authenticator := auth.NewJWT(
		viper.GetString("jwt.secret"),
		time.Duration(viper.GetInt("jwt.duration"))*time.Hour,
		viper.GetString("jwt.signingAlgorithm"),
		userRepo,
	)
	bcryptHasher := hasher.NewBcryptHasher()
	// mailtrap := mailer.NewMailtrap(
	// 	viper.GetString("mailer.mailtrap.username"),
	// 	viper.GetString("mailer.mailtrap.password"),
	// )
	sendgridMailer := mailer.NewSendgrid(
		viper.GetString("mailer.sendgrid.apiKey"),
		debug,
	)
	uuidGenerator := stringgenerator.NewUUID()

	remoteDirectory := "drophere"
	if remoteDirCfg := viper.GetString("app.storageRootDirectoryName"); remoteDirCfg != "" {
		remoteDirectory = remoteDirCfg
	}

	dropboxService := storageprovider.NewDropboxStorageProvider(remoteDirectory)
	storageProviderPool := domain.StorageProviderPool{}
	storageProviderPool.Register(dropboxService)

	basePath := viper.GetString("app.templatePath")
	htmlTemplates, err := htmlTemplate.ParseGlob(filepath.Join(basePath, "html", "*.html"))
	if err != nil {
		panic(err)
	}

	textTemplates, err := textTemplate.ParseGlob(filepath.Join(basePath, "text", "*.txt"))
	if err != nil {
		panic(err)
	}

	// initialize services
	userSvc := user.NewService(
		userRepo,
		userStorageCredRepo,
		authenticator,
		sendgridMailer,
		bcryptHasher,
		uuidGenerator,
		storageProviderPool,
		htmlTemplates,
		textTemplates,
		user.Config{
			PasswordRecoveryTokenExpiryDuration: viper.GetInt("app.passwordRecovery.tokenExpiryDuration"),
			RecoverPasswordWebURL:               viper.GetString("app.passwordRecovery.webURL"),
			MailerEmail:                         viper.GetString("app.passwordRecovery.mailer.email"),
			MailerName:                          viper.GetString("app.passwordRecovery.mailer.name"),
		},
	)
	linkSvc := link.NewService(linkRepo, userStorageCredRepo, bcryptHasher)

	resolver := drophere_go.NewResolver(userSvc, authenticator, linkSvc)

	// setup router
	router := chi.NewRouter()

	// A good base middleware stack
	router.Use(cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		Debug:            debug,
	}).Handler)
	router.Use(authenticator.Middleware())
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Handle("/", handler.Playground("GraphQL playground", "/query"))
	router.Handle("/query", handler.GraphQL(drophere_go.NewExecutableSchema(drophere_go.Config{Resolvers: resolver})))
	router.Post("/uploadfile", fileUploadHandler(userSvc, linkSvc, storageProviderPool))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	err = http.ListenAndServe(":"+port, router)
	if err != nil {
		log.Fatal(err)
	}
}
