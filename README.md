# template-clarify-automation

This is a template repo for running [Clarify.io](https://clarify.io) automation routines using the [clarify-go](https://github.com/clarify/clarify-go) SDK. It currently includes examples for the following one routine:

- `publish`: Automatically bulk publish items from signals by filtering which signals to publish and apply your own _transforms_ to improve the meta-data.

## Who should use the publish routine?

If you have been using Clarify for some time, and you are starting to accumulate a lot of [signals](https://docs.clarify.io/developers/quickstart/create-signals), you may need an automated rule-based process for you to publish signals as items. You may want to script how to clean up the raw meta-data from your [integrations](https://docs.clarify.io/users/admin/integrations/) and produce nicely looking [items](https://docs.clarify.io/users/admin/items/) with _normalized_ labels, naming, engineering unit and other properties. Maybe your data sources encode specific information into the name? Maybe they have their own way of encoding engineering units that just don't make sense inside Clarify?

This automation allows you to write a series of _transform functions_ using the Go programming language. Transforms are again organized into rules. Each rule declare a list of integrations to publish from and an optional signal filter. Be sure that the same signal is only handled by one rule! _Transforms_ are functions that apply changes to the raw signal meta-data. We recommend that you write short and simple transforms that _gradually_ improve the meta-data based on previous steps, rather than one large and overly complex transform that does everything in one operation. This also allows you to more easily compose and reuse transforms across multiple rules. However, it's your choice!

See Video for a hands-on introduction:

[![Video of how to bulk publish routine](https://img.youtube.com/vi/wrDtwy9SyUY/0.jpg)](https://www.youtube.com/watch?v=wrDtwy9SyUY)

## Getting started

This is a _template repository_. To use this code and host it on GitHub, just click the "Use this template" button near the top of the page. Next up, you can clone the repository to your local computer and make your own changes. If your organization does not use GitHub for hosting code, don't worry. Simply create a local clone of this repo and change the origin host to point to your code management system of choice.

To _customize_ this template files, it's recommended that you [configure your editor](https://go.dev/doc/editors) to work with the [Go](https://go.dev/) programming language. When everything is configured correctly, your editor will automatically remove and add imports as you make changes to files, as well as auto-format your code to follow Go conventions. Extensive knowledge of the Go language is not required. One thing you should be aware of, is that file-names doesn't matter to the go tool-chain; content from all files in the same folder, share the same namespace, and the file names matter only to humans.

### Customize publish rules

1. Navigate to Clarify, and copy down the integration ID(s) that you want to publish signals from.
2. Edit the file `publish_rules.go`

You are now ready tur run your automation; to see what it's planning to do, run:

```sh
go run . -credentials clarify-credentials.json -v publish -dry-run
```

To run a sub-set of your publish rules, you can specify the rule names.

```sh
go run . -credentials clarify-credentials.json -v publish -dry-run -rules my-rules
```

When you are happy with the results, you can run it without the `-dry-run` flag. You can also skip the `-v` flag if you want less details.

### Run routines locally

To run this code locally, you need a pair of credentials with access to the _admin_ namespace. For security reasons, do not grant wider permissions than you need. See [our docs](https://docs.clarify.io/developers/quickstart/create-integration) for how to set-up an integration in Clarify and generate new credentials. We recommend that you create a separate integration for using with this repo. For running this locally, we are going to use a _credentials file_. Once you have set-up the credentials, download it and either remember where you placed it, or move it to the root of this repository. You are now ready to start playing with automation routines.

For information on all available commands and options, you can use the `-help` flag:

```sh
- go run . -help          # Help on global flags and available sub-commands.
- go run . publish -help  # Help regarding the publish command.
```

### Run routines in GitHub actions

This repository comes with templates that allows it to run directly in [GitHub Actions][ga]. GitHub Actions is just a very convenient way to run your code without setting up a separate infrastructure to run it in. If you prefer to run your code somewhere else, you can always disable it.

The template is set-up to run the automation routines on:

- New pull-request (using dry-run mode)
- Merge to main (not using dry-run)
- On a fixed schedule (every 24 hours by default)
- On manual trigger

In order to do this, you must:

- In your cloned repository settings, ensure GitHub Actions are enabled (usually on by default).
- Generate username/password credentials for an "automation" integration in Clarify.
- Copy the values and add [secrets][ga-secrets] `CLARIFY_USERNAME` and `CLARIFY_PASSWORD` to the repository.

[ga-secrets]: https://docs.github.com/en/rest/actions/secrets
[ga]: https://github.com/features/actions

### Run routines elsewhere

If you want to run your routines elsewhere, you can build a static binary and deploy it to your preferred destination along with your credentials.
