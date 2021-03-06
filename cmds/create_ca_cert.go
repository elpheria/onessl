package cmds

import (
	"fmt"
	"os"

	"github.com/appscode/kutil/tools/certstore"
	"github.com/spf13/afero"
	"github.com/spf13/cobra"
)

func NewCmdCreateCA(certDir string) *cobra.Command {
	var (
		org       []string
		prefix    string
		overwrite bool
	)
	cmd := &cobra.Command{
		Use:               "ca-cert",
		Short:             "Create CA cert/key pair",
		DisableAutoGenTag: true,
		Run: func(cmd *cobra.Command, args []string) {
			store, err := certstore.NewCertStore(afero.NewOsFs(), certDir, org...)
			if err != nil {
				fmt.Printf("Failed to create certificate store. Reason: %v.", err)
				os.Exit(1)
			}

			var p []string
			if prefix != "" {
				p = append(p, prefix)
			}
			if store.IsExists("ca", p...) && !overwrite {
				fmt.Printf("CA certificate found at %s.", store.Location())
				os.Exit(1)
			}

			err = store.NewCA(p...)
			if err != nil {
				fmt.Printf("Failed to init ca. Reason: %v.", err)
				os.Exit(1)
			}
			fmt.Println("Wrote ca certificates in ", store.Location())
			os.Exit(0)
		},
	}

	cmd.Flags().StringVar(&certDir, "cert-dir", certDir, "Path to directory where pki files are stored.")
	cmd.Flags().StringSliceVarP(&org, "organization", "o", org, "Name of client organizations.")
	cmd.Flags().StringVarP(&prefix, "prefix", "p", prefix, "Prefix added to certificate files")
	cmd.Flags().BoolVar(&overwrite, "overwrite", overwrite, "Overwrite existing cert/key pair")
	return cmd
}
