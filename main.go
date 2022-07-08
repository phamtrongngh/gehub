package main

import (
	"fmt"
)

func main() {
	wsServer := StartWsServer()
	defer wsServer.Close()

	httpServer := SetupHttpServer(wsServer)
	if err := httpServer.Run(
		fmt.Sprintf(":%s", Port),
	); err != nil {
		Logger.Fatalln(err)
	}
}
