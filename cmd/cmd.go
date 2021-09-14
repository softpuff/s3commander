package cmd

import (
	"fmt"
	"os"

	"github.com/softpuff/s3commander/helpers"
	"github.com/spf13/cobra"
)

var (
	Version     string = "v0.0.0"
	region      string
	prefix      string
	debug       bool
	print       bool
	all         bool
	expression  string
	showVersion bool
	version     string
	verbose     bool
)

var (
	S3CommanderCMD = &cobra.Command{
		Use:   "s3commander [subcommand]",
		Short: "List s3 buckets",
	}
)

var (
	versionCMD = &cobra.Command{
		Use:   "version",
		Short: "show version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", Version)
		},
	}
)

var (
	lsCMD = &cobra.Command{
		Use:     "list",
		Aliases: []string{"ls"},
		Short:   "list objects in bucket",
		Args:    cobra.MaximumNArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			result, direct := CompleteArgs(args, region)
			return result, direct
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

			if showVersion {
				for _, o := range objs {

					versions, err := c.GetVersions(args[0], o)
					helpers.BreakOnError(err)

					if verbose {
						for _, v := range versions {
							fmt.Println(v)
						}
					} else {
						for _, v := range versions {
							fmt.Println(*v.VersionId)
						}
					}
				}
				return
			}
			for _, o := range objs {
				fmt.Println(o)
				if print {
					c.PrintS3File(args[0], o, version)
				}

			}

		},
	}
)

var (
	cpCMD = &cobra.Command{
		Use:   "cp [BUCKET] [KEY] [DESTINATION]",
		Short: "cp filename location",
		Args:  cobra.ExactArgs(3),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			result, direct := CompleteArgs(args, region)

			return result, direct
		},
		Run: func(cmd *cobra.Command, args []string) {
			region, err := getRegion(region)
			helpers.BreakOnError(err)

			c := helpers.NewAWSConfig(helpers.WithRegion(region))

			err = c.CpS3file(args[1], args[0], args[2], version, debug)
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
			result, direct := CompleteArgs(args, region)

			return result, direct
		},
		Run: func(cmd *cobra.Command, args []string) {
			region, err := getRegion(region)
			helpers.BreakOnError(err)

			c := helpers.NewAWSConfig(helpers.WithRegion(region))
			err = c.PrintS3File(args[0], args[1], version)
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
			result, direct := CompleteArgs(args, region)

			return result, direct
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
	listVersionsCMD = &cobra.Command{
		Use:   "list-versions",
		Short: "list versions of object",
		Args:  cobra.ExactArgs(2),
		ValidArgsFunction: func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
			return CompleteArgs(args, region)
		},
		Run: func(cmd *cobra.Command, args []string) {
			region, err := getRegion(region)
			helpers.BreakOnError(err)

			c := helpers.NewAWSConfig(helpers.WithRegion(region))
			result, err := c.GetVersions(args[0], args[1])
			helpers.BreakOnError(err)

			vers := helpers.GetObjectVersions(result)

			for _, v := range vers {
				fmt.Println(v)
			}
		},
	}
)

var (
	credsCMD = &cobra.Command{
		Use:   "creds",
		Short: "show creds that are being used",
		Run: func(cmd *cobra.Command, args []string) {
			region, err := getRegion(region)
			helpers.BreakOnError(err)

			c := helpers.NewAWSConfig(helpers.WithRegion(region))
			creds, err := c.Session.Config.Credentials.Get()
			helpers.BreakOnError(err)

			fmt.Printf("AccesKeyId: %s\n", creds.AccessKeyID)
			fmt.Printf("SecretAccessKey: %s\n", creds.SecretAccessKey)
			fmt.Printf("Provider: %s\n", creds.ProviderName)
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
	S3CommanderCMD.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	S3CommanderCMD.AddCommand(completionCmd)
	S3CommanderCMD.AddCommand(lsCMD)
	lsCMD.Flags().BoolVarP(&print, "print", "p", false, "print list results")
	lsCMD.Flags().BoolVarP(&all, "all", "A", false, "return all results, not just first 1000")
	lsCMD.Flags().BoolVarP(&showVersion, "show-version", "V", false, "show versions of s3 objects")
	cpCMD.Flags().StringVarP(&version, "version", "", "", "version of s3 object to print")
	printCMD.Flags().StringVarP(&version, "version", "", "", "version of s3 object to print")
	selectCMD.Flags().StringVarP(&expression, "expression", "E", "", "SQL expression to invoke on s3 object")
	S3CommanderCMD.AddCommand(cpCMD)
	S3CommanderCMD.AddCommand(printCMD)
	S3CommanderCMD.AddCommand(selectCMD)
	S3CommanderCMD.AddCommand(listVersionsCMD)
	S3CommanderCMD.AddCommand(versionCMD)
	S3CommanderCMD.AddCommand(credsCMD)

	cpCMD.RegisterFlagCompletionFunc("version", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 2 {
			bucket := args[0]
			key := args[1]

			region, err := getRegion(region)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			c := helpers.NewAWSConfig(helpers.WithRegion(region))
			vers, err := c.GetVersions(bucket, key)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return helpers.GetObjectVersionsIDs(vers), cobra.ShellCompDirectiveDefault

		} else {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

	})
	printCMD.RegisterFlagCompletionFunc("version", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		if len(args) == 2 {
			bucket := args[0]
			key := args[1]

			region, err := getRegion(region)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			c := helpers.NewAWSConfig(helpers.WithRegion(region))
			vers, err := c.GetVersions(bucket, key)
			if err != nil {
				return nil, cobra.ShellCompDirectiveNoFileComp
			}
			return helpers.GetObjectVersionsIDs(vers), cobra.ShellCompDirectiveDefault

		} else {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}

	})
}

func getRegion(region string) (string, error) {
	reg := os.Getenv("AWS_REGION")
	if reg == "" && region == "" {
		return "", fmt.Errorf("no region")
	}
	return reg, nil
}

func CompleteArgs(args []string, region string) (result []string, direct cobra.ShellCompDirective) {
	c := helpers.NewAWSConfig(helpers.WithRegion(region))
	direct = cobra.ShellCompDirectiveNoFileComp
	switch len(args) {
	case 0:
		result, _ = c.ListS3()

		return

	case 1:
		result, _ = c.ListS3Objects(args[0], "", false)
		return

	default:
		direct = cobra.ShellCompDirectiveDefault
		return
	}
}
