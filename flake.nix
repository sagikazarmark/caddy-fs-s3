{
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
    flake-parts.url = "github:hercules-ci/flake-parts";
    devenv.url = "github:cachix/devenv";
    dagger.url = "github:dagger/nix";
    dagger.inputs.nixpkgs.follows = "nixpkgs";
  };

  outputs =
    inputs@{ flake-parts, ... }:
    flake-parts.lib.mkFlake { inherit inputs; } {
      imports = [
        inputs.devenv.flakeModule
      ];

      systems = [
        "x86_64-linux"
        "x86_64-darwin"
        "aarch64-darwin"
      ];

      perSystem =
        {
          inputs',
          pkgs,
          system,
          ...
        }:
        rec {
          _module.args.pkgs = import inputs.nixpkgs {
            inherit system;

            overlays = [
              (final: prev: {
                dagger = inputs'.dagger.packages.dagger;
              })
            ];
          };

          devenv.shells = {
            default = {
              languages = {
                go.enable = true;
              };

              packages = with pkgs; [
                git
                dagger
                just
                xcaddy
                golangci-lint
              ];

              # https://github.com/cachix/devenv/issues/528#issuecomment-1556108767
              containers = pkgs.lib.mkForce { };
            };

            ci = devenv.shells.default;
          };
        };
    };
}
