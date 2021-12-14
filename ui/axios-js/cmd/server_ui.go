package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Long",
	Run:   Run,
}

func Run(cmd *cobra.Command, args []string) {
	wd, err := os.Getwd()
	if err != nil {
		log.Panic(err)
	}
	httpStaticAssetsDir := http.Dir(fmt.Sprintf("%s/ui/web/static/", wd))

	router := chi.NewRouter()

	// Suggested basic middleware stack from chi's docs
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Handle("/*", http.FileServer(httpStaticAssetsDir))

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

const serverHostFlag = "server.host"
const serverPortFlag = "server.port"
const serverTimeoutReadFlag = "server.timeout.read"
const serverTimeoutWriteFlag = "server.timeout.write"
const serverTimeoutIdleFlag = "server.timeout.idle"

func init() {
	rootCmd.AddCommand(serverCmd)
	var err error

	// HTTP Server
	serverCmd.PersistentFlags().String(serverHostFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverHostFlag,
		serverCmd.PersistentFlags().Lookup(serverHostFlag),
	)
	serverCmd.PersistentFlags().String(serverPortFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverPortFlag,
		serverCmd.PersistentFlags().Lookup(serverPortFlag),
	)
	serverCmd.PersistentFlags().String(serverTimeoutReadFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutReadFlag,
		serverCmd.PersistentFlags().Lookup(serverTimeoutReadFlag),
	)
	serverCmd.PersistentFlags().String(serverTimeoutWriteFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutWriteFlag,
		serverCmd.PersistentFlags().Lookup(serverTimeoutWriteFlag),
	)
	serverCmd.PersistentFlags().String(serverTimeoutIdleFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverTimeoutIdleFlag,
		serverCmd.PersistentFlags().Lookup(serverTimeoutIdleFlag),
	)
	if err != nil {
		log.Panic(err)
	}
}
