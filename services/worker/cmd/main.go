package main

import "github.com/empaid/estateedge/pkg/env"

func main() {
	NewGrpcServer(env.GetString("WORKER_SERVICE_ADDR", ""))
}
