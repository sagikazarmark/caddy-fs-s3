{
  description = "Caddy S3 virtual file system";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec
      {
        devShells = {
          default = pkgs.mkShell {
            buildInputs = with pkgs; [
              git
              go_1_19
              gnumake
              xcaddy
              go-task
              golangci-lint
            ];
          };

          ci = devShells.default;
        };
      });
}
