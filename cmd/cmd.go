package cmd

import (
	"fmt"
	"os"

	"github.com/softpuff/s3commander/helpers"
	"github.com/spf13/cobra"
)

var (
	region     string
	bucket     string
	prefix     string
	key        string
	debug      bool
	dest       string
	print      bool
	all        bool
	expression string
)

var (
	S3CommanderCMD = &cobra.Command{
		Use:   "s3commander [subcommand]",
		Short: "List s3 buckets",
	}
)

var (
	lsCMD = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list objects in bucket",
		Args:    cobra.MaximumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			result := helpers.CompleteArgs(args, region)
			return result, cobra.ShellCompDirectiveNoFileComp
		},
		Run: func(cmd *cobra.Command, args []string) {
			region, err := getRegion(region)
			helpers.BreakOnError(err)

			c := helpers.NewAWSConfig(helpers.WithRegion(region))

			if len(args) == 0 {
				buckets, err := c.ListS3()
				if err != nil {
					fmt.Printf("Error: %v\n", err)
					os.Exit(1)
				}
				for _, b := range buckets {
					fmt.Println(b)
				}
				return
			}
			if len(args) == 1 {
				prefix = ""
			} else {
				prefix = args[1]
			}
			objs, err := c.ListS3Objects(args[0], prefix, all)
			helpers.BreakOnError(err)

			for _, o := range objs {
				fmt.Println(o)
				if print {
					c.PrintS3File(args[0], o)
				}
			}

		},
	}
)

var (
	cpCMD = &cobra.Command{
		Use:   "cp",
		Short: "cp filename location",
		Run: func(cmd *cobra.Command, args []string) {
			region, err := getRegion(region)
			helpers.BreakOnError(err)

			c := helpers.NewAWSConfig(helpers.WithRegion(region))

			err = c.CpS3file(key, bucket, dest, debug)
			helpers.BreakOnError(err)

		},
	}
)

var (
	printCMD = &cobra.Command{
		Use:   "print",
		Short: "print s3 file",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			result := helpers.CompleteArgs(args, region)

			return result, cobra.ShellCompDirectiveDefault
		},
		Run: func(cmd *cobra.Command, args []string) {
			region, err := getRegion(region)
			helpers.BreakOnError(err)

			c := helpers.NewAWSConfig(helpers.WithRegion(region))
			err = c.PrintS3File(args[0], args[1])
			helpers.BreakOnError(err)
		},
	}
)

var (
	selectCMD = &cobra.Command{

		Use:   "select",
		Short: "invoke query on s3 object",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			result := helpers.CompleteArgs(args, region)

			return result, cobra.ShellCompDirectiveDefault
		},
		Run: func(cmd *cobra.Command, args []string) {
			region, err := getRegion(region)
			helpers.BreakOnError(err)

			c := helpers.NewAWSConfig(helpers.WithRegion(region))
			err = c.CountS3ObjectLines(args[0], args[1], expression)
			helpers.BreakOnError(err)

		},
	}
)

var (
	completionCmd = &cobra.Command{
		Use:   "completion",
		Short: "Generates bash completion script",
		Run: func(cmd *cobra.Command, args []string) {
			S3CommanderCMD.GenBashCompletion(os.Stdout)
		},
	}
)

func init() {
	S3CommanderCMD.PersistentFlags().StringVarP(&region, "region", "r", "", "AWS region")
	S3CommanderCMD.PersistentFlags().BoolVarP(&debug, "debug", "g", false, "debug output")
	S3CommanderCMD.AddCommand(completionCmd)
	S3CommanderCMD.AddCommand(lsCMD)
	cpCMD.Flags().StringVarP(&key, "key", "f", "", "key to copy from s3")
	cpCMD.Flags().StringVarP(&bucket, "bucket", "b", "", "bucket name")
	cpCMD.Flags().StringVarP(&dest, "destination", "d", "", "destination folder for s3 file")
	lsCMD.Flags().BoolVarP(&print, "print", "p", false, "print list results")
	lsCMD.Flags().BoolVarP(&all, "all", "A", false, "return all results, not just first 1000")
	selectCMD.Flags().StringVarP(&expression, "expression", "E", "", "SQL expression to invoke on s3 object")
	S3CommanderCMD.AddCommand(cpCMD)
	S3CommanderCMD.AddCommand(printCMD)
	S3CommanderCMD.AddCommand(selectCMD)

}

func getRegion(region string) (string, error) {
	reg := os.Getenv("AWS_REGION")
	if reg == "" && region == "" {
		return "", fmt.Errorf("no region")
	}
	return reg, nil
}
