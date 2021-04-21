package cmd

import (
	"fmt"
	"os"

	"github.com/softpuff/s3commander/helpers"
	"github.com/spf13/cobra"
)

var region string
var bucket string
var prefix string
var key string
var debug bool
var dest string

var (
	S3dirCMD = &cobra.Command{
		Use:   "s3commander [subcommand]",
		Short: "List s3 buckets",
	}
)

var (
	listCMD = &cobra.Command{
		Use:   "list-buckets",
		Short: "list all the buckets",
		Run: func(cmd *cobra.Command, args []string) {
			region := getRegion(region)
			if region == "" {
				fmt.Fprintln(os.Stderr, "No region")
				cmd.Help()
				os.Exit(1)
			}
			c := helpers.NewAWSConfig(region)
			buckets, err := c.ListS3()
			if err != nil {
				fmt.Fprintf(os.Stderr, "Listing buckets: %v\n", err)
				os.Exit(1)
			}

			for _, b := range buckets {
				fmt.Println(b)
			}
		},
	}
)

var (
	lsCMD = &cobra.Command{
		Use:   "ls",
		Short: "list objects in bucket",
		Args:  cobra.MaximumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			c := helpers.NewAWSConfig(region)
			if len(args) == 0 {

				result, _ := c.ListS3()
				return result, cobra.ShellCompDirectiveDefault
			}
			if len(args) == 1 {
				result, _ := c.ListS3Objects(args[0], "")
				return result, cobra.ShellCompDirectiveDefault
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
		Run: func(cmd *cobra.Command, args []string) {
			region := getRegion(region)
			if region == "" {
				fmt.Fprintln(os.Stderr, "No region")
				cmd.Help()
				os.Exit(1)
			}
			c := helpers.NewAWSConfig(region)
			if len(args) == 1 {
				prefix = ""
			} else {
				prefix = args[1]
			}
			objs, err := c.ListS3Objects(args[0], prefix)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Listing bucket objects for %s: %v\n", bucket, err)
				os.Exit(1)
			}
			for _, o := range objs {
				fmt.Println(o)
			}

		},
	}
)

var (
	cpCMD = &cobra.Command{
		Use:   "cp",
		Short: "cp filename location",
		Run: func(cmd *cobra.Command, args []string) {
			region := getRegion(region)
			if region == "" {
				fmt.Fprintln(os.Stderr, "No region")
				cmd.Help()
				os.Exit(1)
			}
			c := helpers.NewAWSConfig(region)

			if err := c.CpS3file(key, bucket, dest, debug); err != nil {
				fmt.Fprintf(os.Stderr, "Error copying %s from %s: %v\n", key, bucket, err)

			}

		},
	}
)

var (
	printCMD = &cobra.Command{
		Use:   "print",
		Short: "print s3 file",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			c := helpers.NewAWSConfig(region)
			if len(args) == 0 {

				result, _ := c.ListS3()
				return result, cobra.ShellCompDirectiveDefault
			}
			if len(args) == 1 {
				result, _ := c.ListS3Objects(args[0], "")
				return result, cobra.ShellCompDirectiveDefault
			}
			return nil, cobra.ShellCompDirectiveDefault
		},
		Run: func(cmd *cobra.Command, args []string) {
			c := helpers.NewAWSConfig(region)
			if err := c.PrintS3File(args[0], args[1]); err != nil {
				fmt.Printf("Error: %v", err)
				os.Exit(1)
			}
		},
	}
)

var (
	completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion script",
		Run: func(cmd *cobra.Command, args []string) {
			S3dirCMD.GenBashCompletion(os.Stdout)
		},
	}
)

func init() {
	S3dirCMD.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region")
	S3dirCMD.PersistentFlags().BoolVarP(&debug, "debug", "g", false, "debug output")
	S3dirCMD.AddCommand(completionCmd)
	S3dirCMD.AddCommand(listCMD)
	lsCMD.Flags().StringVarP(&bucket, "bucket", "b", "", "bucket name")
	lsCMD.Flags().StringVarP(&prefix, "prefix", "p", "", "prefix for bucket list")
	S3dirCMD.AddCommand(lsCMD)
	cpCMD.Flags().StringVarP(&key, "key", "f", "", "key to copy from s3")
	cpCMD.Flags().StringVarP(&bucket, "bucket", "b", "", "bucket name")
	cpCMD.Flags().StringVarP(&dest, "destination", "d", "", "destination folder for s3 file")
	S3dirCMD.AddCommand(cpCMD)
	S3dirCMD.AddCommand(printCMD)
	printCMD.Flags().StringVarP(&bucket, "bucket", "b", "", "bucket name")
	printCMD.Flags().StringVarP(&key, "key", "k", "", "bucket object key")
}

func getRegion(region string) string {
	reg := os.Getenv("AWS_REGION")
	if reg == "" && region == "" {
		return ""
	}
	return reg
}
