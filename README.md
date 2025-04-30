# Go HTMX App

A quick-start app template using Go, Templ, Tailwind, DaisyUI and HTMX

## Included features

- authentication using a [SAML Identity Provider](https://github.com/silinternational/ssp-base)
- basic homepage starting point built with the Templ templating engine for Go, styled using Tailwind CSS and DaisyUI, and enhanced with HTMX for interactivity.
- email notification using [MailGun](https://www.mailgun.com/) or [AWS SES](https://aws.amazon.com/ses/)
- database migration using [Goose](https://github.com/pressly/goose)
- database connection using [sqlc](https://github.com/sqlc-dev/sqlc)
- error logging using [logrus](https://github.com/sirupsen/logrus) and [Sentry](https://sentry.io/welcome/) remote option

## Packages

### action

HTTP handlers

General principles:

- Restrict use of `echo.Context` to `action` package

For instance, use `toCtx(c)` to convert `echo.Context` to `context.Context`.

- No `data` types

Convert all `data` types to `app` types before rendering http response. Ideally, the `core` package would do this before returning to `action`.

- Simple http handlers

Handler functions should be simple, containing only authentication, authorization, and data type conversion. All processing, data manipulation, and logic should be in the `core` package.

### api

Placeholder for REST API structures. It currently contains only an error type.

### app

Application initialization and configuration.

Startup sequence

```
app.init() -> app.readEnv()
main.main() -> log.ErrLogger.Init()
            -> sql.Open
            -> action.NewApp()
            -> App.Start()
```

### cmd

The `main` package for the server and a command utility for scheduling with a cron service.

### core

Business logic

Should have no direct access to `sqlc` package. Use or create a wrapper in `data` package if necessary.

### data

Data access -- a thin wrapper around the `sqlc` generated package.

### email

Email notification service

### log

Logging service

### public

Assets for the user interface and email messages

### saml

SAML authentication

# Getting started

- optional: create a local.env file in the project root and add variables as described in local-example.env
- run `make` and wait for the build to complete
- run `docker compose logs -f app` and wait for the app to build and show "http server started on [::]:80"
- open a browser to http://localhost:8100
- login with username "john_doe" and password "boot promote elegant bottle"
