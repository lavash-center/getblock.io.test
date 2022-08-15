package main

import (
	"context"
	"log"
	"net/http"
	"os/signal"
	"syscall"

	"github.com/lavash-center/getblock.io.test/src/managers"
	"github.com/lavash-center/getblock.io.test/src/resources"
	"github.com/ybbus/jsonrpc/v2"
	"golang.org/x/sync/errgroup"

	"github.com/lavash-center/getblock.io.test/src/config"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	var opts config.Configuration

	err := config.ParseConfig(&opts)
	if err != nil {
		log.Fatal("[ERROR]", err)
	}

	rpcCli := jsonrpc.NewClientWithOpts(opts.GetBlockRPCEndpoint, &jsonrpc.RPCClientOpts{
		CustomHeaders: map[string]string{
			"x-api-key": opts.ApiKey,
		},
	})

	blocksMan := managers.NewBlocksManagerImpl(rpcCli)

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		srv := http.Server{
			Addr:    opts.Listen,
			Handler: resources.NewResource("/", blocksMan).Routes(),
		}

		// Graceful shutdown
		go func() {
			<-gCtx.Done()
			log.Printf("[INFO] shutting down http server on %s\n", srv.Addr)
			err = srv.Shutdown(gCtx)
			if err != nil {
				log.Printf("[ERROR] shutting down http server: %s\n", err.Error())
			}

			log.Printf("[INFO] http server on %s processed all idle connections\n", srv.Addr)
		}()

		log.Printf("[INFO] serving HTTP on %s\n", srv.Addr)
		return srv.ListenAndServe()
	})

	err = g.Wait()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal("[ERROR]", err)
	}
}
