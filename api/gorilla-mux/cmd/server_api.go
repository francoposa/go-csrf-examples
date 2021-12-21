package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/csrf"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Long",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	router := mux.NewRouter()

	loggingMiddleware := func(h http.Handler) http.Handler {
		return handlers.LoggingHandler(os.Stdout, h)
	}
	router.Use(loggingMiddleware)

	// This has to get applied later because Gorilla Mux is a pain
	CORSMiddleware := cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice(serverCORSAllowedOriginsFlag),
		AllowCredentials: viper.GetBool(serverCORSAllowCredentialsFlag),
		AllowedHeaders:   viper.GetStringSlice(serverCORSAllowedHeadersFlag),
		ExposedHeaders:   viper.GetStringSlice(serverCORSExposedHeadersFlag),
		Debug:            viper.GetBool(serverCORSDebugFlag),
	}).Handler

	CSRFMiddleware := csrf.Protect(
		[]byte(viper.GetString(serverCSRFKeyFlag)),
		csrf.Secure(viper.GetBool(serverCSRFSecureFlag)),
		csrf.CookieName(viper.GetString(serverCSRFCookieName)),
		csrf.RequestHeader(viper.GetString(serverCSRFHeader)),
	)

	APIRouter := router.PathPrefix("/api").Subrouter()
	APIRouter.Use(CSRFMiddleware)
	APIRouter.HandleFunc("", Get).Methods(http.MethodGet)
	APIRouter.HandleFunc("", Post).Methods(http.MethodPost)

	host := viper.GetString(serverHostFlag)
	port := viper.GetString(serverPortFlag)
	readTimeout := viper.GetInt(serverTimeoutReadFlag)
	writeTimeout := viper.GetInt(serverTimeoutWriteFlag)
	idleTimeout := viper.GetInt(serverTimeoutIdleFlag)
	server := &http.Server{
		Handler:      CORSMiddleware(router),
		Addr:         host + ":" + port,
		ReadTimeout:  time.Duration(readTimeout) * time.Second,
		WriteTimeout: time.Duration(writeTimeout) * time.Second,
		IdleTimeout:  time.Duration(idleTimeout) * time.Second,
	}

	fmt.Printf("starting http server on port %s...\n", port)
	log.Fatal(server.ListenAndServe())
}

func Get(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("X-CSRF-Token", csrf.Token(r))
	w.WriteHeader(http.StatusOK)
}

func Post(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

const serverHostFlag = "server.host"
const serverPortFlag = "server.port"
const serverTimeoutReadFlag = "server.timeout.read"
const serverTimeoutWriteFlag = "server.timeout.write"
const serverTimeoutIdleFlag = "server.timeout.idle"
const serverCORSAllowCredentialsFlag = "server.cors.allowCredentials"
const serverCORSAllowedHeadersFlag = "server.cors.allowedHeaders"
const serverCORSExposedHeadersFlag = "server.cors.exposedHeaders"
const serverCORSAllowedOriginsFlag = "server.cors.allowedOrigins"
const serverCORSAllowedMethodsFlag = "server.cors.allowedMethods"
const serverCORSDebugFlag = "server.cors.debug"
const serverCSRFSecureFlag = "server.csrf.secure"
const serverCSRFKeyFlag = "server.csrf.key"
const serverCSRFCookieName = "server.csrf.cookieName"
const serverCSRFHeader = "server.csrf.header"

func init() {
	rootCmd.AddCommand(serverCmd)

	cmdFlags := serverCmd.Flags()
	var err error

	// HTTP Server
	cmdFlags.String(serverHostFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverHostFlag, cmdFlags.Lookup(serverHostFlag),
	)
	serverCmd.PersistentFlags().String(serverPortFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverPortFlag, cmdFlags.Lookup(serverPortFlag),
	)
	cmdFlags.String(serverTimeoutReadFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutReadFlag, cmdFlags.Lookup(serverTimeoutReadFlag),
	)
	cmdFlags.String(serverTimeoutWriteFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutWriteFlag, cmdFlags.Lookup(serverTimeoutWriteFlag),
	)
	cmdFlags.String(serverTimeoutIdleFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutIdleFlag, cmdFlags.Lookup(serverTimeoutIdleFlag),
	)
	// HTTP server CORS
	cmdFlags.String(serverCORSAllowCredentialsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSAllowCredentialsFlag, cmdFlags.Lookup(serverCORSAllowCredentialsFlag),
	)
	cmdFlags.String(serverCORSAllowedHeadersFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSAllowedHeadersFlag, cmdFlags.Lookup(serverCORSAllowedHeadersFlag),
	)
	cmdFlags.String(serverCORSExposedHeadersFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSExposedHeadersFlag, cmdFlags.Lookup(serverCORSExposedHeadersFlag),
	)
	cmdFlags.String(serverCORSAllowedMethodsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSAllowedMethodsFlag, cmdFlags.Lookup(serverCORSAllowedMethodsFlag),
	)
	cmdFlags.String(serverCORSAllowedOriginsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSAllowedOriginsFlag, cmdFlags.Lookup(serverCORSAllowedOriginsFlag),
	)
	cmdFlags.String(serverCORSDebugFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSDebugFlag, cmdFlags.Lookup(serverCORSDebugFlag),
	)
	// HTTP Server CSRF
	cmdFlags.String(serverCSRFSecureFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCSRFSecureFlag, cmdFlags.Lookup(serverCSRFSecureFlag),
	)
	cmdFlags.String(serverCSRFKeyFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCSRFKeyFlag, cmdFlags.Lookup(serverCSRFKeyFlag),
	)
	cmdFlags.String(serverCSRFCookieName, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCSRFCookieName, cmdFlags.Lookup(serverCSRFCookieName),
	)
	cmdFlags.String(serverCSRFHeader, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(serverCSRFHeader, cmdFlags.Lookup(serverCSRFHeader))
	if err != nil {
		log.Panic(err)
	}
}
