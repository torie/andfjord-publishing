// Copyright 2022 Searis AS
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/clarify/clarify-go"
	"github.com/clarify/clarify-go/automation"
	"github.com/clarify/clarify-go/jsonrpc"
	"github.com/peterbourgon/ff/v3"
	"github.com/peterbourgon/ff/v3/ffcli"
	"golang.org/x/exp/maps"
	"golang.org/x/exp/slices"
)

func init() {
	log.SetFlags(0)
	log.SetPrefix(progName + ": ")
}

func rootCommand() *ffcli.Command {
	var opts rootOptions
	fs := flag.NewFlagSet(progName, flag.ExitOnError)
	fs.StringVar(&opts.username, "username", "", "Clarify integration ID to use as username.")
	fs.StringVar(&opts.password, "password", "", "Clarify integration password (required when username is set, ignored otherwise).")
	fs.StringVar(&opts.credentialsFile, "credentials", "clarify-credentials.json", "Clarify credentials file location (ignored if username is set).")
	fs.BoolVar(&opts.log.Verbose, "v", false, "Set the program to be extra verbose.")
	fs.BoolVar(&opts.quiet, "q", false, "Set the program to be extra quiet (mostly only log errors).")

	return &ffcli.Command{
		ShortUsage: progName + " [flags] <subcommand>",
		ShortHelp:  "Clarify Automation CLI (from template)",
		FlagSet:    fs,
		Exec: func(ctx context.Context, args []string) error {
			return fmt.Errorf("subcommand required; try -help")
		},
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithAllowMissingConfigFile(true),
			ff.WithEnvVarPrefix("CLARIFY"),
		},
		Subcommands: []*ffcli.Command{
			publishCommand(&opts),
		},
	}
}

type rootOptions struct {
	// raw parameters.
	username        string
	password        string
	credentialsFile string
	quiet           bool

	// runtime variables.
	client *clarify.Client
	log    automation.LogOptions
}

func (p *rootOptions) init(ctx context.Context) {
	var creds *clarify.Credentials
	if p.username != "" {
		if p.password == "" {
			log.Println("fatal: password is required when a username is provided.")
			os.Exit(2)
		}
		creds = clarify.BasicAuthCredentials(p.username, p.password)
	} else {
		var err error
		creds, err = clarify.CredentialsFromFile(p.credentialsFile)
		if err != nil {
			log.Println("fatal:", err)
			os.Exit(2)
		}
	}

	h, err := creds.HTTPHandler(ctx)
	if err != nil {
		log.Println("fatal:", err)
		os.Exit(2)
	}
	if !p.quiet {
		p.log.Out = os.Stderr
	}
	if p.log.Verbose {
		h.RequestLogger = func(req jsonrpc.Request, trace string, latency time.Duration, err error) {
			p.log.Printf("JSONRPC request: %s, trace: %s, latency: %s, error: %v\n", req.Method, trace, latency, err)
		}
	}
	p.client = clarify.NewClient(creds.Integration, h)

	p.log.Printf("Configured using %s for integration %s.\n", creds.Credentials.Type, creds.Integration)
}

const publishLongHelp = `
(Re-)publish signals as items using a pre-defined set of rules. Each rule define
a list of integration IDs to publish signals from, a filter of for which signals
to select, and a list of transforms to apply.
`

func publishCommand(root *rootOptions) *ffcli.Command {
	var opts publish
	opts.availableRules = maps.Keys(publishRules)
	slices.Sort(opts.availableRules)

	fs := flag.NewFlagSet(progName, flag.ExitOnError)
	fs.BoolVar(&opts.DryRun, "dry-run", false, "Set to not persist changes to Clarify.")
	fs.Var(stringSlice{target: &opts.rules}, "rules", "Comma separated list rules to run (default is all).")

	return &ffcli.Command{
		Name:       "publish",
		ShortUsage: progName + " [flags] publish [command flags]",
		ShortHelp:  "(Re-)publish signals as items using a set of rules.",
		LongHelp: fmt.Sprintf(
			"%s\n\nRULES:\n- %s",
			strings.TrimSpace(publishLongHelp),
			strings.Join(opts.availableRules, "\n- "),
		),
		FlagSet: fs,
		Exec: func(ctx context.Context, args []string) error {
			opts.init(ctx, root)
			return opts.do(ctx)
		},
		Options: []ff.Option{
			ff.WithConfigFileFlag("config"),
			ff.WithConfigFileParser(ff.PlainParser),
			ff.WithAllowMissingConfigFile(true),
			ff.WithEnvVarPrefix("CLARIFY"),
		},
	}
}

type publish struct {
	automation.PublishOptions
	rules []string

	availableRules []string
	client         *clarify.Client
}

func (opts *publish) init(ctx context.Context, root *rootOptions) {
	root.init(ctx)
	opts.client = root.client
	opts.LogOptions = root.log

	if len(opts.rules) == 0 {
		opts.rules = maps.Keys(publishRules)
	} else {
		for _, k := range opts.rules {
			_, ok := publishRules[k]
			if !ok {
				log.Printf("no publish rule named %q; see -help", k)
				os.Exit(2)
			}
		}
	}
}

func (opts *publish) do(ctx context.Context) error {
	for _, k := range opts.rules {
		if err := ctx.Err(); err != nil {
			return err
		}
		opts.Printf("== Running publish rule: %s\n", k)
		if err := publishRules[k].Do(ctx, opts.client, opts.PublishOptions); err != nil {
			return err
		}
	}
	return nil
}
