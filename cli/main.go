package main

import (
	"fmt"
	"net/url"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	version = "none"
	commit  = "none"
	app     = &cobra.Command{
		Use:   "ende",
		Short: "Handle common encodings on the commandline",
	}
	cmdURLEncoding = &cobra.Command{
		Use:   "url",
		Short: "Handle URL encoding",
		Run:   executeURLEncoding,
	}
	cmdVersion = &cobra.Command{
		Use:   "version",
		Short: "Show endes version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s (ref: %s)\n", version, commit)
		},
	}
)

func init() {
	app.PersistentFlags().StringP("loglevel", "l", "warn", "Minimum level for logmessages")
	viper.BindPFlag("loglevel", app.PersistentFlags().Lookup("loglevel"))
	app.PersistentFlags().String("logfile", "-", "Write logfiles to the given file (- for stderr)")
	viper.BindPFlag("logfile", app.PersistentFlags().Lookup("logfile"))

	cmdURLEncoding.PersistentFlags().BoolP("decode", "d", false, "Decode given argument")

	app.AddCommand(cmdURLEncoding, cmdVersion)
	viper.SetEnvPrefix("ENDE")
	viper.AutomaticEnv()
}

func main() {
	if err := app.Execute(); err != nil {
		panic(err)
	}
}

func executeURLEncoding(cmd *cobra.Command, args []string) {
	decode, _ := cmd.Flags().GetBool("decode")
	for _, arg := range args {
		if decode {
			fmt.Println(url.QueryUnescape(arg))
		} else {
			fmt.Println(url.QueryEscape(arg))
		}
	}
}
