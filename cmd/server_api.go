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

const serverAPIHostFlag = "serverAPI.host"
const serverAPIPortFlag = "serverAPI.port"
const serverAPITimeoutReadFlag = "serverAPI.timeout.read"
const serverAPITimeoutWriteFlag = "serverAPI.timeout.write"
const serverAPITimeoutIdleFlag = "serverAPI.timeout.idle"
const serverAPICORSAllowCredentialsFlag = "serverAPI.cors.allowCredentials"
const serverAPICORSAllowedHeadersFlag = "serverAPI.cors.allowedHeaders"
const serverAPICORSExposedHeadersFlag = "serverAPI.cors.exposedHeaders"
const serverAPICORSAllowedOriginsFlag = "serverAPI.cors.allowedOrigins"
const serverAPICORSAllowedMethodsFlag = "serverAPI.cors.allowedMethods"
const serverAPICORSDebugFlag = "serverAPI.cors.debug"
const serverAPICSRFSecureFlag = "serverAPI.csrf.secure"
const serverAPICSRFKeyFlag = "serverAPI.csrf.key"
const serverAPICSRFCookieName = "serverAPI.csrf.cookieName"
const serverAPICSRFHeader = "serverAPI.csrf.header"

func init() {
	rootCmd.AddCommand(serverAPICmd)

	cmdFlags := serverAPICmd.Flags()
	var err error

	// HTTP Server
	cmdFlags.String(serverAPIHostFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPIHostFlag, cmdFlags.Lookup(serverAPIHostFlag),
	)
	serverAPICmd.PersistentFlags().String(serverAPIPortFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPIPortFlag, cmdFlags.Lookup(serverAPIPortFlag),
	)
	cmdFlags.String(serverAPITimeoutReadFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPITimeoutReadFlag, cmdFlags.Lookup(serverAPITimeoutReadFlag),
	)
	cmdFlags.String(serverAPITimeoutWriteFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPITimeoutWriteFlag, cmdFlags.Lookup(serverAPITimeoutWriteFlag),
	)
	cmdFlags.String(serverAPITimeoutIdleFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPITimeoutIdleFlag, cmdFlags.Lookup(serverAPITimeoutIdleFlag),
	)
	// HTTP server CORS
	cmdFlags.String(serverAPICORSAllowCredentialsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICORSAllowCredentialsFlag, cmdFlags.Lookup(serverAPICORSAllowCredentialsFlag),
	)
	cmdFlags.String(serverAPICORSAllowedHeadersFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICORSAllowedHeadersFlag, cmdFlags.Lookup(serverAPICORSAllowedHeadersFlag),
	)
	cmdFlags.String(serverAPICORSExposedHeadersFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICORSExposedHeadersFlag, cmdFlags.Lookup(serverAPICORSExposedHeadersFlag),
	)
	cmdFlags.String(serverAPICORSAllowedMethodsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICORSAllowedMethodsFlag, cmdFlags.Lookup(serverAPICORSAllowedMethodsFlag),
	)
	cmdFlags.String(serverAPICORSAllowedOriginsFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICORSAllowedOriginsFlag, cmdFlags.Lookup(serverAPICORSAllowedOriginsFlag),
	)
	cmdFlags.String(serverAPICORSDebugFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICORSDebugFlag, cmdFlags.Lookup(serverAPICORSDebugFlag),
	)
	// HTTP Server CSRF
	cmdFlags.String(serverAPICSRFSecureFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICSRFSecureFlag, cmdFlags.Lookup(serverAPICSRFSecureFlag),
	)
	cmdFlags.String(serverAPICSRFKeyFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICSRFKeyFlag, cmdFlags.Lookup(serverAPICSRFKeyFlag),
	)
	cmdFlags.String(serverAPICSRFCookieName, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverAPICSRFCookieName, cmdFlags.Lookup(serverAPICSRFCookieName),
	)
	cmdFlags.String(serverAPICSRFHeader, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(serverAPICSRFHeader, cmdFlags.Lookup(serverAPICSRFHeader))
	if err != nil {
		panic(err)
	}
}

var serverAPICmd = &cobra.Command{
	Use:   "server-api",
	Short: "Long",
	Run:   run,
}

func run(cmd *cobra.Command, args []string) {
	router := chi.NewRouter()

	// Suggested basic middleware stack from chi's docs
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	CORSMiddleware := cors.New(cors.Options{
		AllowedOrigins:   viper.GetStringSlice(serverAPICORSAllowedOriginsFlag),
		AllowCredentials: viper.GetBool(serverAPICORSAllowCredentialsFlag),
		AllowedHeaders:   viper.GetStringSlice(serverAPICORSAllowedHeadersFlag),
		ExposedHeaders:   viper.GetStringSlice(serverAPICORSExposedHeadersFlag),
		Debug:            viper.GetBool(serverAPICORSDebugFlag),
	}).Handler
	router.Use(CORSMiddleware)

	CSRFMiddleware := csrf.Protect(
		[]byte(viper.GetString(serverAPICSRFKeyFlag)),
		csrf.Secure(viper.GetBool(serverAPICSRFSecureFlag)),
		csrf.CookieName(viper.GetString(serverAPICSRFCookieName)),
		csrf.RequestHeader(viper.GetString(serverAPICSRFHeader)),
	)

	router.Route("/api", func(router chi.Router) {
		router.With(CSRFMiddleware).Get("/", Get)
		router.With(CSRFMiddleware).Post("/", Post)
	})

	host := viper.GetString(serverAPIHostFlag)
	port := viper.GetString(serverAPIPortFlag)
	readTimeout := viper.GetInt(serverAPITimeoutReadFlag)
	writeTimeout := viper.GetInt(serverAPITimeoutWriteFlag)
	idleTimeout := viper.GetInt(serverAPITimeoutIdleFlag)
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
