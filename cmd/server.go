// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"net/http"
	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/guenter/network-quality/utils"
	"time"
)

var serverAddress string
var pingInternal time.Duration

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("server called")

		var pingTimes = prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: "network_quality",
			Name: "ping_time_seconds",
			Help: "Ping response time from a host in seconds.",
		}, []string{"ip"})

		prometheus.MustRegister(pingTimes)

		go ping(args, pingTimes)

		http.Handle("/metrics", promhttp.Handler())
		log.Printf("Listening on %s", serverAddress)
		log.Fatal(http.ListenAndServe(serverAddress, nil))
	},
}

func ping(ips []string, pingTimes *prometheus.HistogramVec) {
	for true {
		for _, ip := range ips {
			duration, err := utils.Ping(networkInterface, ip)
			if err != nil {
				log.Fatal(err)
			}
			log.Printf("Got reply from %s in %s", ip, duration)
			pingTimes.With(prometheus.Labels{"ip": ip}).Observe(duration.Seconds())
		}
		time.Sleep(pingInternal)
	}
}

func init() {
	RootCmd.AddCommand(serverCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// serverCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// serverCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	serverCmd.Flags().StringVar(&serverAddress, "serverAddress", ":8100", "Address to listen on")
	serverCmd.Flags().DurationVar(&pingInternal, "pingInterval", 10 * time.Second, "Time between pings")
}
