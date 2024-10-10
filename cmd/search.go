package cmd

import (
	"fmt"
	"net/url"
	"os"

	"github.com/alex27riva/crtsh/fetcher"
	"github.com/olekukonko/tablewriter"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// searchCmd represents the search command
var searchCmd = &cobra.Command{
	Use:   "search",
	Short: "search",
	Long:  `search`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return search()
	},
	SilenceErrors: true,
	SilenceUsage:  true,
}

func init() {
	rootCmd.AddCommand(searchCmd)

	searchCmd.Flags().StringP("query", "q", "", "query (e.g. Facebook)")
	searchCmd.Flags().StringP("domain", "d", "", "Domain Name (e.g. %.exmaple.com)")
	searchCmd.Flags().Bool("plain", false, "plain text mode")
	viper.BindPFlag("query", searchCmd.Flags().Lookup("query"))
	viper.BindPFlag("domain", searchCmd.Flags().Lookup("domain"))
	viper.BindPFlag("plain", searchCmd.Flags().Lookup("plain"))
}

func search() (err error) {
	query := viper.GetString("query")
	domain := viper.GetString("domain")
	if query == "" && domain == "" {
		return errors.New("--query or --domain must be specified")
	}


	values := url.Values{}
	values.Add("output", "json")
	if query != "" {
		values.Add("q", query)
	} else if domain != "" {
		values.Add("q", domain)
	}

	// if query != "" {
	// 	for _, result := range cer {
	// 		url := fmt.Sprintf("%s?output=json&id=%d", crtURL, result.MinCertID)
	// 		urls = append(urls, url)
	// 	}
	// 	certs := fetcher.FetchConcurrently(urls, 5, 0)
	// 	if err != nil {
	// 		return errors.Wrap(err, "Failed to fetch concurrently")
	// 	}

	// 	table := tablewriter.NewWriter(os.Stdout)
	// 	table.SetHeader([]string{"Common Name", "Organization", "Locality", "Country", "Not After"})
	// 	table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
	// 		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
	// 		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
	// 		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
	// 		tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
	// 	)

	// 	for _, cert := range certs {
	// 		if !viper.GetBool("plain") {
	// 			table.Append([]string{cert.CommonName, cert.OrganizationName, cert.LocalityName, cert.CountryName, cert.NotAfter})
	// 		} else {
	// 			fmt.Println(cert.CommonName)
	// 		}
	// 	}
	// 	if !viper.GetBool("plain") {
	// 		table.Render()
	// 	}
	// } else
	
if domain != "" {
		certs := fetcher.fetchCertificates(domain)
		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Name", "Issuer", "Not Before", "Not After"})
		table.SetColumnColor(tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiRedColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
			tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiBlackColor},
		)
		for _, result := range certs {
			if !viper.GetBool("plain") {
				table.Append([]string{result.NameValue, result.IssuerName, result.NotBefore, result.NotAfter})
			} else {
				fmt.Println(result.NameValue)
			}
		}
		if !viper.GetBool("plain") {
			table.Render()
		}
	}

	return nil
}
