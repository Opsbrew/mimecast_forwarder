   
/*
Copyright © 2019 Adron Hall <adron@thrashingcode.com>
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
	"github.com/spf13/cobra"
	"github.com/Opsbrew/mimecast_forwarder/helper"
	"fmt"
)

// addCmd represents the add command
var addCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the forwarder",
	Long: `Start the forwarder`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Checking remote syslog connection")
		connection := helper.Raw_connect()
		if(connection){
			fmt.Println("Started")
			base_url := helper.Get_base_url("email")
			if(base_url=="no_url"){
				fmt.Println("Error discovering base url for the given email ID")
			}else{
				fmt.Println("Base url found,",base_url)
			}

			ok := true

			for(ok){
				ok = helper.Get_mta_siem_logs(base_url)
			}
		}else{
			fmt.Println("Checking remote syslog connection timedout")
		}
		
		
		
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}