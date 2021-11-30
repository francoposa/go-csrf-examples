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
)

const serverUIHostFlag = "serverUI.host"
const serverUIPortFlag = "serverUI.port"
const serverUITimeoutReadFlag = "serverUI.timeout.read"
const serverUITimeoutWriteFlag = "serverUI.timeout.write"
const serverUITimeoutIdleFlag = "serverUI.timeout.idle"

func init() {
	rootCmd.AddCommand(serverUICmd)
	var err error

	// HTTP Server
	serverUICmd.PersistentFlags().String(serverUIHostFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverUIHostFlag,
		serverUICmd.PersistentFlags().Lookup(serverUIHostFlag),
	)
	serverUICmd.PersistentFlags().String(serverUIPortFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverUIPortFlag,
		serverUICmd.PersistentFlags().Lookup(serverUIPortFlag),
	)
	serverUICmd.PersistentFlags().String(serverUITimeoutReadFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverUITimeoutReadFlag,
		serverUICmd.PersistentFlags().Lookup(serverUITimeoutReadFlag),
	)
	serverUICmd.PersistentFlags().String(serverUITimeoutWriteFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverUITimeoutWriteFlag,
		serverUICmd.PersistentFlags().Lookup(serverUITimeoutWriteFlag),
	)
	serverUICmd.PersistentFlags().String(serverUITimeoutIdleFlag, "", "")
	//nolint:ineffassign,staticcheck
	err = viper.BindPFlag(
		serverUITimeoutIdleFlag,
		serverUICmd.PersistentFlags().Lookup(serverUITimeoutIdleFlag),
	)
	if err != nil {
		panic(err)
	}
}

var serverUICmd = &cobra.Command{
	Use:   "server-ui",
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

	host := viper.GetString(serverUIHostFlag)
	port := viper.GetString(serverUIPortFlag)
	readTimeout := viper.GetInt(serverUITimeoutReadFlag)
	writeTimeout := viper.GetInt(serverUITimeoutWriteFlag)
	idleTimeout := viper.GetInt(serverUITimeoutIdleFlag)
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
