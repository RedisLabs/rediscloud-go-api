{
  description = "Flake for Redis Cloud Go API";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs =
    {
      self,
      nixpkgs,
      flake-utils,
    }:
    flake-utils.lib.eachDefaultSystem (
      system:
      let
        pkgs = import nixpkgs {
          inherit system;
        };
        go = pkgs.go_1_26.overrideAttrs (finalAttrs: _: {
          version = "1.26.11";
          src = pkgs.fetchurl {
            url = "https://go.dev/dl/go${finalAttrs.version}.src.tar.gz";
            hash = "sha256-1EEBsP8n7SpMda1GwXucdDWSCeLrLKfJ32y0/baOGWk=";
          };
        });
      in
      {
        devShells.default = pkgs.mkShell {
          buildInputs = [
            pkgs.gnumake

            go
            pkgs.golangci-lint
            pkgs.gotools
            pkgs.govulncheck
          ];
        };
      }
    );
}
