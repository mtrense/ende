package main

import (
	"encoding/base64"
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
		Use:     "url",
		Aliases: []string{"u"},
		Short:   "Handle URL encoding",
		Run:     executeURLEncoding,
	}
	cmdBase64 = &cobra.Command{
		Use:     "base64",
		Aliases: []string{"b64"},
		Short:   "Handle Base64 encoding",
		Run:     executeBase64,
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
	cmdURLEncoding.PersistentFlags().BoolP("path", "p", false, "Use path-safe encoding")

	cmdBase64.PersistentFlags().BoolP("decode", "d", false, "Decode given argument")

	app.AddCommand(cmdURLEncoding, cmdBase64, cmdVersion)
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
	path, _ := cmd.Flags().GetBool("path")
	convertAndPrintArgs(args, func(s string) (string, error) {
		if decode {
			if path {
				return url.PathUnescape(s)
			}
			return url.QueryUnescape(s)
		}
		if path {
			return url.PathEscape(s), nil
		}
		return url.QueryEscape(s), nil
	})
}

func executeBase64(cmd *cobra.Command, args []string) {
	decode, _ := cmd.Flags().GetBool("decode")
	convertAndPrintArgs(args, func(s string) (string, error) {
		if decode {
			res, err := base64.StdEncoding.DecodeString(s)
			return string(res), err
		}
		return base64.StdEncoding.EncodeToString([]byte(s)), nil
	})
}

func convertAndPrintArgs(args []string, fn func(s string) (string, error)) {
	for _, arg := range args {
		res, err := fn(arg)
		if err != nil {
			panic(err)
		}
		fmt.Println(res)
	}
}
