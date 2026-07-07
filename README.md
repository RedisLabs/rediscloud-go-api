# rediscloud-go-api

This repository is a Go SDK for the [Redis Cloud REST API](https://docs.redislabs.com/latest/rc/api/).

## Getting Started

### Installing
You can use this module by using `go get` to add it to either your `GOPATH` workspace or
the project's dependencies.
```shell script
go get github.com/RedisLabs/rediscloud-go-api
```

### Example
This is an example of using the SDK
```go
package main

import (
	"context"
	"fmt"

	rediscloud_api "github.com/RedisLabs/rediscloud-go-api"
	"github.com/RedisLabs/rediscloud-go-api/service/subscriptions"
)

func main() {
	// The client will use the credentials from `REDISCLOUD_ACCESS_KEY` and `REDISCLOUD_SECRET_KEY` by default
	client, err := rediscloud_api.NewClient()
	if err != nil {
		panic(err)
	}

	id, err := client.Subscription.Create(context.TODO(), subscriptions.CreateSubscription{
		// ...
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Created subscription: %d", id)
}
```

## Development

The repo ships a [Nix flake](flake.nix) that pins the toolchain (Go, `golangci-lint`, `gotools`, `govulncheck`, GNU `make`) so local development matches CI. Combined with [direnv](https://direnv.net/), your shell automatically enters that environment when you `cd` into the repo.

### Setup

1. Install Nix with flakes enabled — the [Determinate Nix Installer](https://determinate.systems/nix-installer/) is the easiest option; the [official Nix installer](https://nixos.org/download/) works too but requires [enabling flakes manually](https://nixos.wiki/wiki/Flakes#Enable_flakes_permanently_in_NixOS).
2. Install [direnv](https://direnv.net/docs/installation.html) and hook it into your shell.
3. Install [nix-direnv](https://github.com/nix-community/nix-direnv) so direnv can load the flake's dev shell.
4. From the repo root, activate the environment:
   ```shell
   direnv allow
   ```

Every subsequent `cd` into the repo brings up a shell with the pinned toolchain on `PATH`. Run `make` to execute the same checks CI runs.

## Releasing

Releases are published as Git tags of the form `vX.Y.Z`, so Go consumers pull a
specific release with `go get github.com/RedisLabs/rediscloud-go-api@vX.Y.Z`.

Tagging is automated on merge to `main`, using the [Tag Release workflow](.github/workflows/tag-release.yml). The `Version` constant in [`version.go`](version.go) is
the single source of truth for the release number.

To cut a release:
1. Bump `Version` in [`version.go`](version.go) (e.g. `0.51.0` → `0.52.0`).
2. Make sure there is a matching entry in [`CHANGELOG.md`](CHANGELOG.md).
3. Open a PR and merge it into `main`.
