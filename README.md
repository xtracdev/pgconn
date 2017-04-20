## PGConn - Postgres Connection Utilities


This package provide a simple mechanism to extend the Postgres sql.DB
implementation with some additional capabilities:

* The ability to retry initial connection attempts using a simple backoff
 mechanism. This is useful in scenarios such as starting containerized
 applications without wanting to worry about start up order.
 * Detecting certain classes of connection related errors.
 * Reconnecting to the database using retrys and backoff.
 
## Usage
 
 * Use the `OpenAndConnect` method to instantiate a `*sql.DB` instance.
 * Use `IsConnectionError` to determine if the error indicates the connection
 to the database is lost or hopeless, in which `Reconnect` should be used
 to reconnect to the database.
 
## Dependencies
 
<pre>
go get github.com/lib/pq
go get github.com/gucumber/gucumber/cmd/gucumber
go get github.com/stretchr/testify/assert
go get github.com/Sirupsen/logrus
</pre>

## Contributing

To contribute, you must certify you agree with the [Developer Certificate of Origin](http://developercertificate.org/)
by signing your commits via `git -s`. To create a signature, configure your user name and email address in git.
Sign with your real name, do not use pseudonyms or submit anonymous commits.


In terms of workflow:

0. For significant changes or improvement, create an issue before commencing work.
1. Fork the respository, and create a branch for your edits.
2. Add tests that cover your changes, unit tests for smaller changes, acceptance test
for more significant functionality.
3. Run gofmt on each file you change before committing your changes.
4. Run golint on each file you change before committing your changes.
5. Make sure all the tests pass before committing your changes.
6. Commit your changes and issue a pull request.

## License

(c) 2017 Fidelity Investments
Licensed under the Apache License, Version 2.0
