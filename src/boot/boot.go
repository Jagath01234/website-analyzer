package boot

import (
	"context"
	"fmt"
	"github.com/gocolly/colly/v2"
	"github.com/google/uuid"
	"github.com/karlseguin/ccache/v2"
	"github.com/pickme-go/traceable-context"
	"log"
	netHttp "net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"
	"website-analyzer/src/configuration"
	"website-analyzer/src/controller"
	"website-analyzer/src/service"
	"website-analyzer/src/transport/http"
	"website-analyzer/src/worker"
)

const logPrefixBoot = "web-analyzer/src/boot/boot"

func Run() {
	ctx, stop := signal.NotifyContext(traceable_context.WithUUID(uuid.New()), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	configPath := filepath.Join(".", "config.json")

	configuration.LoadConfig(configPath)

	cache := ccache.New(ccache.Configure().MaxSize(configuration.AppConfig.Cache.MaxSize).ItemsToPrune(configuration.AppConfig.Cache.PruneSize))

	analyzerSvc := service.NewAnalyzerService(colly.NewCollector())
	worker := worker.NewAnalyzerWorker(analyzerSvc, cache, make(chan int64, configuration.AppConfig.Worker.BufferSize), make(chan os.Signal, 1))
	analyzerCont := controller.NewAnalyzerController(cache, worker)
	worker.InitAnalyzerWorkerPool()

	http.Start(analyzerCont)
	if configuration.AppConfig.Pprof.IsEnabled {
		go func() {
			err := netHttp.ListenAndServe(fmt.Sprintf("localhost:%v", configuration.AppConfig.Pprof.Port), nil)
			log.Printf("%v,%v,%v,%v", "ERROR", logPrefixBoot, "pprof 6060 port serve error", err.Error())
		}()
	}
	<-ctx.Done()
	stop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	http.Shutdown(shutdownCtx)

	cancel()
}
