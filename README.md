# KintoHub Go Utils

A project that holds common logic used across all go services @ KintoHub.

## Logger

Basic logger using zero log. This is a dependency across all of utils and expected to call `InitLogger` to set it up.

Logger expects an environment variable `LOG_LEVEL` of value `VERBOSE, DEBUG, INFO, WARN, ERROR, PANIC, FATAL`


## Server

Server currently only has an implementation of grpc but may have future server implementations such as fasthttp, websockets
etc.  There are common utilities within server so that it can be abstracted.  Most importantly, `server/utils/errors.go` can
be used to create standard errors across different server implementations so that your business logic can return errors
such as `NotFound` or `Internal` and depending on the implementation, it will handle the error code and message gracefully.

Additionally, any errors that are returned and sent out will be automatically logged through middleware. So logging errors
is unnecessary when using this package.

## Config

Config package has utility functions to load configuration from environment variables with ease.

## Utils

Utils package has basic logic that is reused across projects.
