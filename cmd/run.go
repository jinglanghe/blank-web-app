package cmd

import (
	config "github.com/apulis/bmod/aistudio-aom/configs"
	"github.com/apulis/bmod/aistudio-aom/internal/cache"
	"github.com/apulis/bmod/aistudio-aom/internal/controllers"
	"github.com/apulis/bmod/aistudio-aom/internal/dao"
	"github.com/apulis/bmod/aistudio-aom/internal/service"
	"github.com/apulis/sdk/go-utils/logging"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"os"
	"os/signal"
	"syscall"
)

// alertCmd represents the alert command
var alertCmd = &cobra.Command{
	Use:   "run",
	Short: "Run aom service",
	Run: func(cmd *cobra.Command, args []string) {
		config.Init()
		cache.Init()
		cache.InitRedLock()
		dao.Init()
		service.Init()

		e := gin.Default()
		controllers.RegisterRoutes(e)

		go func() {
			if err := e.Run(":" + config.Config.Port); err != nil {
				logging.Fatal().Err(err).Msg("failed to start alert service")
			}

		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 5 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		logging.Info().Msg("Shutting down server...")
	},
}

func init() {
	rootCmd.AddCommand(alertCmd)
}
