{
  inputs = {
    gomod2nix.url = "github:mexisme/gomod2nix";

    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, gomod2nix, flake-utils, ... }@inputs:
    {
      overlays.default = (final: prev: {
        assume-role = final.callPackage ./. {};
      });
    } //
    flake-utils.lib.eachDefaultSystem (system:
      let
        pkgs = import nixpkgs {
          inherit system;
          overlays = [ gomod2nix.overlays.default ];
        };
      in
        {
          packages.default = pkgs.callPackage ./. { };
        });
}
