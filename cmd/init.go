// Copyright 2017 The kubecfg authors
//
//
//    Licensed under the Apache License, Version 2.0 (the "License");
//    you may not use this file except in compliance with the License.
//    You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//    Unless required by applicable law or agreed to in writing, software
//    distributed under the License is distributed on an "AS IS" BASIS,
//    WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//    See the License for the specific language governing permissions and
//    limitations under the License.

package cmd

import (
	"fmt"
	"os"
	"path"

	"k8s.io/client-go/tools/clientcmd/api"

	"github.com/ksonnet/ksonnet/metadata"
	"github.com/ksonnet/ksonnet/pkg/kubecfg"
	"github.com/spf13/cobra"
)

func init() {
	RootCmd.AddCommand(initCmd)
	// TODO: We need to make this default to checking the `kubeconfig` file.
	initCmd.PersistentFlags().String(flagAPISpec, "version:v1.7.0",
		"Manually specify API version from OpenAPI schema, cluster, or Kubernetes version")

	bindClientGoFlags(initCmd)
}

var initCmd = &cobra.Command{
	Use:   "init <app-name>",
	Short: "Initialize a ksonnet project",
	RunE: func(cmd *cobra.Command, args []string) error {
		flags := cmd.Flags()
		if len(args) != 1 {
			return fmt.Errorf("'init' takes a single argument that names the application we're initializing")
		}

		appName := args[0]
		appDir, err := os.Getwd()
		if err != nil {
			return err
		}
		appRoot := metadata.AbsPath(path.Join(appDir, appName))

		specFlag, err := flags.GetString(flagAPISpec)
		if err != nil {
			return err
		}

		//
		// Find the URI of the current cluster, if it exists.
		//

		rawConfig, err := clientConfig.RawConfig()
		if err != nil {
			return err
		}

		var currCtx *api.Context
		for name, ctx := range rawConfig.Contexts {
			if name == rawConfig.CurrentContext {
				currCtx = ctx
			}
		}

		var currClusterURI *string
		if rawConfig.CurrentContext != "" && currCtx != nil {
			for name, cluster := range rawConfig.Clusters {
				if currCtx.Cluster == name {
					currClusterURI = &cluster.Server
				}
			}
		}

		c, err := kubecfg.NewInitCmd(appRoot, specFlag, currClusterURI, &currCtx.Namespace)
		if err != nil {
			return err
		}

		return c.Run()
	},
	Long: `Initialize a ksonnet project in a new directory, 'app-name'. This process
consists of two steps:

1. Generating ksonnet-lib. Users can set flags to generate the library based on
   a variety of data, including server configuration and an OpenAPI
   specification of a Kubernetes build. By default, this is generated from the
   capabilities of the cluster specified in the cluster of the current context
   specified in $KUBECONFIG.
2. Generating the following tree in the current directory.

   app-name/
     .gitignore     Default .gitignore; can customize VCS
     .ksonnet/      Metadata for ksonnet
     envs/
       default/     Default generated environment]
         k.libsonnet
         k8s.libsonnet
         swagger.json
         spec.json
     components/    Top-level Kubernetes objects defining application
     lib/           user-written .libsonnet files
     vendor/        mixin libraries, prototypes
`,
	Example: `  # Initialize ksonnet application, using the capabilities of live cluster
  # specified in the $KUBECONFIG environment variable (specifically: the
  # current context) to generate 'ksonnet-lib'.
  ks init app-name

  # Initialize ksonnet application, using the OpenAPI specification generated
  # in the Kubenetes v1.7.1 build to generate 'ksonnet-lib'.
  ks init app-name --api-spec=version:v1.7.1

  # Initialize ksonnet application, using an OpenAPI specification file
  # generated by a build of Kubernetes to generate 'ksonnet-lib'.
  ks init app-name --api-spec=file:swagger.json`,
}
