/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"log"
	"time"

	"github.com/aleosada/otel-jaeger-demo/pkg/server"
	"github.com/spf13/cobra"
)

var port *string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
        tp := server.ConfigOtelJaeger()
        ctx, cancel := context.WithCancel(context.Background())
        defer cancel()

        // Cleanly shutdown and flush telemetry when the application exits.
        defer func(ctx context.Context) {
            // Do not make the application hang when it is shutdown.
            ctx, cancel = context.WithTimeout(ctx, time.Second*5)
            defer cancel()
            if err := tp.Shutdown(ctx); err != nil {
                log.Fatal(err)
            }
        }(ctx)

        server.InitServer(port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

    port = serveCmd.Flags().StringP("port", "p", ":3000", "Port formar ':{number}'")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serveCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serveCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
