package main

import (
	"CoalCompany/delivery/http"
	"CoalCompany/domain"
	"context"
	"fmt"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())

	enterprise := domain.NewEnterprise(ctx, cancel)
	go enterprise.StartPassiveIncome()

	httpHandlers := http.NewHTTPHandlers(enterprise)
	httpServer := http.NewHTTPServer(httpHandlers)

	go func() {
		if err := httpServer.StartServer(); err != nil {
			fmt.Println("failed to start sever")
			cancel()
		}
	}()

	<-ctx.Done()
}
