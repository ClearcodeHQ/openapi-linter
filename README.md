# openapi-linter

![status](https://img.shields.io/badge/status-alpha-red.svg)
[![actions atatus](https://github.com/clearcodehq/openapi-linter/workflows/Pipeline/badge.svg)](https://github.com/clearcodehq/openapi-linter/actions)
[![go report](https://goreportcard.com/badge/github.com/clearcodehq/openapi-linter)](https://goreportcard.com/report/github.com/clearcodehq/openapi-linter)
[![godoc](https://godoc.org/github.com/clearcodehq/openapi-linter?status.svg)](http://godoc.org/github.com/clearcodehq/openapi-linter)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c)](https://pkg.go.dev/github.com/clearcodehq/openapi-linter)
[![license](https://img.shields.io/badge/License-MIT-blue.svg)](https://github.com/clearcodehq/openapi-linter/blob/master/LICENSE)


> :warning: `openapi-linter` is an Clearcode's internal tool that helps with the validation of some parts of the OpenAPI specification.
> It serves mostly as the first Golang project suffers the NIH syndrome.
> if you need a production ready solution, you should check projects from the [GopenAPI initiative](https://github.com/gopenapi).

## Getting started

This project requires Go to be installed. On OS X with Homebrew you can just run `brew install go`.

Running it then should be as simple as:

```console
$ make
$ ./bin/openapi-linter
```

### Testing

``make test``

# Contributing

Contributions are welcome, and they are greatly appreciated! Every little bit helps, and credit will always be given.

You can contribute in many ways:

## Types of Contributions

### Report Bugs

Report bugs at https://github.com/clearcodehq/openapi-linter/issues/

If you are reporting a bug, please include:

* Your operating system name and version.
* Any details about your local setup that might be helpful in troubleshooting.
* Detailed steps to reproduce the bug.


### Write Documentation

openapi-linter could always use more documentation, whether as part of the
official openapi-linter docs, in docstrings, or even on the web in blog posts,
articles, and such.

### Submit Feedback

The best way to send feedback is to file an issue at https://github.com/clearcodehq/openapi-linter/issues/

If you are proposing a feature:

* Explain in detail how it would work.
* Keep the scope as narrow as possible, to make it easier to implement.

## Get Started!

Ready to contribute? Here's how to set up `openapi-linter` for local development.

1. Fork the `openapi-linter` repo on GitHub.
2. Clone your fork locally::
```bash
    $ git clone git@github.com:clearcodehq/openapi-linter.git
```
3. Create a branch for local development::
```bash
    $ git checkout -b fix-<GITHUB_ISSUE_NUMBER>-helpful-keywords
```
   Now you can make your changes locally.

4. When you're done making changes, check that your changes pass the tests::
```bash
    $ make test
```
6. Commit your changes and push your branch to GitHub::
```bash
    $ git add .
    $ git commit -m "Your detailed description of your changes."
    $ git push origin fix-<GITHUB ISSUE NUMBER>-helpful-keywords
```
7. Submit a pull request through the GitHub website.

Pull Request Guidelines
-----------------------

Before you submit a pull request, check that it meets these guidelines:

1. The pull request should include tests.
2. If the pull request adds functionality, the docs should be updated. Put
   your new functionality into a function with a docstring, and add the
   feature to the list in README.md.
