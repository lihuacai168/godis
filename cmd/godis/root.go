/*
Copyright © 2021 lihuacai

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
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/lihuacai168/godis/cmd/config"
	"github.com/mattn/go-colorable"
	"github.com/spf13/cobra"
	"io"
	"log"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var (
	outWriter io.Writer = os.Stdout
	errWriter io.Writer = os.Stderr
	inReader  io.Reader = os.Stdin

	colorableOut io.Writer = colorable.NewColorableStdout()
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "godis",
	Short: "A brief description of your application",

	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		outWriter = cmd.OutOrStdout()
		errWriter = cmd.ErrOrStderr()
		inReader = cmd.InOrStdin()

		if outWriter != os.Stdout {
			colorableOut = outWriter
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

var cfg config.Config

var currentCluster *config.Cluster
var rdb *redis.ClusterClient
var ctx = context.Background()
var (
	password           string
	clusterDescription string
	addrs              []string
	clusterOverride    string
)

func init() {
	cobra.OnInitialize(onInit)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".godis" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".godis")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}
}
func onInit() {
	var err error
	cfg, err = config.ReadConfig(cfgFile)
	if err != nil {
		errorExit("Invalid config: %v", err)
	}

	cfg.ClusterOverride = clusterOverride

	cluster := cfg.ActiveCluster()
	if cluster != nil {
		// Use active cluster from config
		currentCluster = cluster
	} else {
		// Create sane default if not configured
		errorExit(`not configured,please use "godis config add" to set a cluster configuration`)
	}

	initClient()
}
func initClient() {
	a := currentCluster.Addrs
	p := currentCluster.Password
	rdb = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    a,
		Password: p,
	})
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		errorExit("use %s conf, addr: %s, password: %s\nPing server: %v", currentCluster.Name, currentCluster.Addrs, currentCluster.Password, err)
	}
	log.Printf("use %s conf, addr is %s connected success\n", currentCluster.Name, a)
}
