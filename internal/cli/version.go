package cli

// Version specifies the build version of the CLI.
// It should be injected at build time via -ldflags="-X 'stackgenome/internal/cli.Version=vX.Y.Z'"
var Version = "dev"
