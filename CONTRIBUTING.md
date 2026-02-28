# Contributing

By participating in this project, you agree to abide our code of conduct, found
[here](https://github.com/jsnfwlr/filamate/blob/main/CODE_OF_CONDUCT.md).

**Note**: While we plan to allow for internationalization of the UI, we require
all non-internationalized code
(comments and table/column/file/variable/function/method names) to be written in
English.

## Set up your machine

The backend of `filamate` is written in [Go](https://go.dev/), while the
frontend is written in [VueJS](https://vuejs.org/).
[PostgreSQL](https://www.postgresql.org/) is used for persistent data storage,
and [Docker](https://www.docker.com/) test-containers are used when running
backend unit tests.

Prerequisites:

- [Go 1.26+](https://go.dev/doc/install)
- [VueJS](https://vuejs.org/guide/quick-start.html)
- [Docker](https://www.docker.com/)
- [just](https://github.com/casey/just?tab=readme-ov-file#installation):
    The command runner
- [pnpm](https://pnpm.io/installation):
    The performant node package manager

Other things you will need to install to run tests and build code (these are
installed by running `just tools`):

- [Air](https://github.com/air-verse/air):
    Go live rebuild tool
- [tparse](https://github.com/mfridman/tparse):
    Go test parser
- [sqlc](https://github.com/sqlc-dev/sqlc):
    Generates Go code from SQL queries
- [oapi-codegen](https://github.com/oapi-codegen/oapi-codegen):
    Generates Go code from Open API specifications

## Building

Fork `filamate` to your own GitHub account, then clone the fork anywhere
(replace `%%YOUR_ACCOUNT_NAME%%` with tour GitHub account name):

```shell
git clone git@github.com:%%YOUR_ACCOUNT_NAME%%/filamate.git
```

`cd` into the directory and install the dependencies:

```shell
just dep
```

You should then be able to build the binary:

```shell
just compile
./dist/filamate --version
```

## Running `filamate` with your changes

You can create a branch for your changes and try to build from the source as you
go, running filamate with the API live-reloading whenever a change is detected.

```shell
just run
```

To rebuild the UI on change, open a second terminal to the same folder:

```shell
just watch
```

## Testing and linting your changes

When you are satisfied with the changes, we suggest you run:

```shell
just test
just lint
```

**Note**: Test coverage on newly added functionality is required. While we don't
expect 100% coverage, we do need sufficient testing of the code to see the
expected behaviour of it in enough scenarios to deem it fit for purpose.

## Creating commits

Commits should be atomic - do not include multiple features or unrelated fixes
in a single commit, and ideally, do not spread a single feature or related fixes
out over multiple commits with other, unrelated changes between them.

Commit messages should be well formatted, and to make that "standardized", we
are using [Conventional Commits](https://www.conventionalcommits.org).

The _Type_, _Scope_, and _Description_ are all required.

Allowed values for _Type_ are `fix`, `feat`, `chore`, and `refactor`.
Allowed prefixes for _Scope_ are:
- `API` for changes to the backend code
- `CI` for changes to the continuous integration and/or build pipelines
- `CLI` for changes to the command-line
- `DB` for changes to the database made via the migrations
- `DOC` for changes to the documentation of the application or project.
- `UI` for changes to the frontend code

While breaking changes should be avoided, they are inevitable and commits that
include breaking changes should have a ! prefix for the first line of the commit
message.

### Examples

[This commit](https://github.com/jsnfwlr/filamate/commit/94b6fbe36f7eb8e0beda3086b8ebf7bd1203ebd8) should have had the following commit message:
```
feat: Added material chart
- (UI) Added pie chart showing the breakdown of filament spools by class of material, specific material, and brand
- (UI) Added material chart store to stats.ts to handle retrieving the chart data from the API
- (API) Added /chart/material endpoint to provide the data for the new chart
- (API) Added GetMaterialChartData dashboard query to retrieve the chart data from the database
```
Had the change been a breaking-change, the first line of the commit message would have been:

```
!feat: Added material chart
```

## Submitting a pull request

Push your branch to your `filamate` fork and open a pull request against the main branch.


