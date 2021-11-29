package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"
	"github.com/rs/cors"
)

const serverHostFlag = "serverBackend.host"
const serverPortFlag = "serverBackend.port"
const serverTimeoutReadFlag = "serverBackend.timeout.read"
const serverTimeoutWriteFlag = "serverBackend.timeout.write"
const serverTimeoutIdleFlag = "serverBackend.timeout.idle"
const serverCORSAllowCredentialsFlag = "serverBackend.cors.allowCredentials"
const serverCORSAllowedHeadersFlag = "serverBackend.cors.allowedHeaders"
const serverCORSExposedHeadersFlag = "serverBackend.cors.exposedHeaders"
const serverCORSAllowedOriginsFlag = "serverBackend.cors.allowedOrigins"
const serverCORSAllowedMethodsFlag = "serverBackend.cors.allowedMethods"
const serverCORSDebugFlag = "serverBackend.cors.debug"
const serverCSRFSecureFlag = "serverBackend.csrf.secure"
const serverCSRFKeyFlag = "serverBackend.csrf.key"
const serverCSRFCookieName = "serverBackend.csrf.cookieName"
const serverCSRFHeader = "serverBackend.csrf.header"

func init() {
	rootCmd.AddCommand(serverBackendCmd)
	var err error

	// HTTP Server
	serverBackendCmd.PersistentFlags().String(serverHostFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverHostFlag, serverBackendCmd.PersistentFlags().Lookup(serverHostFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverPortFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverPortFlag, serverBackendCmd.PersistentFlags().Lookup(serverPortFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverTimeoutReadFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutReadFlag, serverBackendCmd.PersistentFlags().Lookup(serverTimeoutReadFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverTimeoutWriteFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutWriteFlag, serverBackendCmd.PersistentFlags().Lookup(serverTimeoutWriteFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverTimeoutIdleFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutIdleFlag, serverBackendCmd.PersistentFlags().Lookup(serverTimeoutIdleFlag),
	)
	// HTTP server CORS
	serverBackendCmd.PersistentFlags().String(serverCORSAllowCredentialsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSAllowCredentialsFlag,
		serverBackendCmd.PersistentFlags().Lookup(serverCORSAllowCredentialsFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverCORSAllowedHeadersFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSAllowedHeadersFlag,
		serverBackendCmd.PersistentFlags().Lookup(serverCORSAllowedHeadersFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverCORSExposedHeadersFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSExposedHeadersFlag,
		serverBackendCmd.PersistentFlags().Lookup(serverCORSExposedHeadersFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverCORSAllowedMethodsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSAllowedMethodsFlag,
		serverBackendCmd.PersistentFlags().Lookup(serverCORSAllowedMethodsFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverCORSAllowedOriginsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSAllowedOriginsFlag,
		serverBackendCmd.PersistentFlags().Lookup(serverCORSAllowedOriginsFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverCORSDebugFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCORSDebugFlag,
		serverBackendCmd.PersistentFlags().Lookup(serverCORSDebugFlag),
	)
	// HTTP Server CSRF
	serverBackendCmd.PersistentFlags().String(serverCSRFSecureFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCSRFSecureFlag, serverBackendCmd.PersistentFlags().Lookup(serverCSRFSecureFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverCSRFKeyFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCSRFKeyFlag, serverBackendCmd.PersistentFlags().Lookup(serverCSRFKeyFlag),
	)
	serverBackendCmd.PersistentFlags().String(serverCSRFCookieName, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCSRFCookieName, serverBackendCmd.PersistentFlags().Lookup(serverCSRFCookieName),
	)
	serverBackendCmd.PersistentFlags().String(serverCSRFHeader, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverCSRFHeader, serverBackendCmd.PersistentFlags().Lookup(serverCSRFHeader),
	)
	if err != nil {
		panic(err)
	}
}

var serverBackendCmd = &cobra.Command{
	Use:   "server-backend",
	Short: "Long",
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	router := chi.NewRouter()

	// Suggested basic middleware stack from chi's docs
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	CORSMiddleware := cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice(serverCORSAllowedOriginsFlag),
		AllowCredentials: viper.GetBool(serverCORSAllowCredentialsFlag),
		AllowedHeaders:   viper.GetStringSlice(serverCORSAllowedHeadersFlag),
		ExposedHeaders:   viper.GetStringSlice(serverCORSExposedHeadersFlag),
		Debug:            viper.GetBool(serverCORSDebugFlag),
	}).Handler
	router.Use(CORSMiddleware)

	CSRFMiddleware := csrf.Protect(
		[]byte(viper.GetString(serverCSRFKeyFlag)),
		csrf.Secure(viper.GetBool(serverCSRFSecureFlag)),
		csrf.CookieName(viper.GetString(serverCSRFCookieName)),
		csrf.RequestHeader(viper.GetString(serverCSRFHeader)),
	)

	router.Route("/api", func(router chi.Router) {
		router.With(CSRFMiddleware).Get("/", Get)
		router.With(CSRFMiddleware).Post("/", Post)
	})

	host := viper.GetString(serverHostFlag)
	port := viper.GetString(serverPortFlag)
	readTimeout := viper.GetInt(serverTimeoutReadFlag)
	writeTimeout := viper.GetInt(serverTimeoutWriteFlag)
	idleTimeout := viper.GetInt(serverTimeoutIdleFlag)
	server := &http.Server{
		Handler:      router,
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
