/*
Copyright © 2023 Rob Duarte <me@robduarte.com>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"fmt"
	"os"

	"github.com/hypebeast/go-osc/osc"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "oscxyz",
	Short: "An OSC-to-WebSockets bridge",
	Long: `
oscxyz is a simple OSC-to-WebSockets bridge that takes OSC messages from 
a client, like TouchOSC, and sends them to a WebSocket client, like a p5js sketch.

Note that this implementation is very basic and was created specifically to handle 
accelerometer data with an OSC type tag of ",fff" and an address pattern of "/accxyz"
(although the address pattern can be changed with the --pattern flag).
`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		oschostname, _ := cmd.Flags().GetString("oschost")
		oscport, _ := cmd.Flags().GetInt("oscport")
		//websocketsport, _ := cmd.Flags().GetInt("wsport")
		pattern, _ := cmd.Flags().GetString("pattern")
		quiet, _ := cmd.Flags().GetBool("quiet")

		startOSCServer(oschostname, oscport, pattern, quiet)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	rootCmd.Flags().String("oschost", "", "IP address to use when creating the OSC server (required)")
	rootCmd.MarkFlagRequired("oschost")
	rootCmd.Flags().Int("oscport", 0, "Port number to use when creating the OSC server (required)")
	rootCmd.MarkFlagRequired("oscport")
	rootCmd.Flags().Int("wsport", 0, "Port number to use when creating the WebSockets server (required)")
	rootCmd.MarkFlagRequired("wsport")

	rootCmd.Flags().String("pattern", "/accxyz", "OSC message pattern to listen for")
	rootCmd.Flags().BoolP("quiet", "q", false, "Don't show OSC messages on the console")
}

func startOSCServer(oschostname string, oscport int, pattern string, quiet bool) {
	addr := fmt.Sprintf("%s:%d", oschostname, oscport)

	d := osc.NewStandardDispatcher()
	d.AddMsgHandler(pattern, func(msg *osc.Message) {
		if !quiet {
			osc.PrintMessage(msg)
		}
	})

	server := &osc.Server{
		Addr:       addr,
		Dispatcher: d,
	}
	go server.ListenAndServe()
	fmt.Printf("OSC Server Listening on %s\n", addr)

	// keep the program running until a signal is received
	select {}
}
