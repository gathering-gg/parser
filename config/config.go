package config

// Version of the Gathering.gg Client
var Version string

// Root uri for the API. If you want to change this for local development,
// override it with a build `-ldflag`. The Makefile does this automatically for
// development builds.
// go build -ldflags "-X 'gitlab.com/gathering-gg/gathering/config.Root=http://localhost"
var Root = "https://api.gathering.gg"
