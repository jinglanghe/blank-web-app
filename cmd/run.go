package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	config "gitlab.apulis.com.cn/hjl/blank-web-app/configs"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/cache"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/controllers"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/dao"
	"gitlab.apulis.com.cn/hjl/blank-web-app/internal/service"
	"gitlab.apulis.com.cn/hjl/blank-web-app/logging"
)

// alertCmd represents the alert command
var alertCmd = &cobra.Command{
	Use:   "run",
	Short: "Run blankWebApp service",
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
