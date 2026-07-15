{
  description = "flake.nix of sumika";
  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixos-unstable";
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
          config = {
            android_sdk.accept_license = true;
            allowUnfree = true;
          };
        };
      in
      {
        devShell =
          with pkgs;
          mkShell {
            buildInputs = with pkgs; [
              go
              gopls
              gotools
              delve
            ];

            hardeningDisable = [ "fortify" ];
          };
      }
    );
}