{
  description = "Caddy S3 virtual file system";

  inputs = {
    nixpkgs.url = "nixpkgs/nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
    ci.url = "github:sagikazarmark/nix-ci-utils";
    ci.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs = { self, nixpkgs, flake-utils, ci, ... }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs { inherit system; };
      in
      rec {
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
        }
        //
        (ci.lib.genShellsFromList [ "1_19" "1_20" ] (goVersion:
          devShells.default.overrideAttrs (final: prev: {
            buildInputs = [ pkgs."go_${goVersion}" ] ++ prev.buildInputs;
          })
        ));
      }
    );
}
