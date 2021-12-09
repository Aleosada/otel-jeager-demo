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

	"github.com/aleosada/otel-jaeger-demo/pkg/client"
	"github.com/spf13/cobra"
)

var interval, requests *int
var withError *string

// clientCmd represents the client command
var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Calls http endpoint for tracing demo",
	Run: func(cmd *cobra.Command, args []string) {
        tp := client.ConfigOtelJaeger()
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

        client.Run(requests, interval, *withError)
	},
}

func init() {
	rootCmd.AddCommand(clientCmd)

	interval = clientCmd.Flags().IntP("interval", "i", 0, "Interval for the requests in seconds")
	requests = clientCmd.Flags().IntP("requests", "r", 1, "Number of requests")
    withError = clientCmd.Flags().StringP("error", "e", "", "Error api path")
}
