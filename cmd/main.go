package main

import (
	"log"

	"github.com/laupse/kubegraph/adapter/http"
	"github.com/laupse/kubegraph/adapter/k8s"
	"github.com/laupse/kubegraph/application/services"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func main() {
	pflag.String("Kubeconfig", "", "Kubeconfig")
	pflag.String("Kubecontext", "", "Kubecontext")
	pflag.Parse()
	viper.BindPFlags(pflag.CommandLine)

	repository, err := k8s.NewK8sRepository()
	if err != nil {
		log.Fatal(err.Error())
	}

	graphService := services.NewGraphService(repository)

	fiberHandler := http.NewFiberHandler(graphService)
	fiberHandler.SetupRoutes()
	fiberHandler.Run(":3000")

}
