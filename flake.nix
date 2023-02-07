{
  inputs = {
    gomod2nix.url = "github:mexisme/gomod2nix";

    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, gomod2nix, flake-utils, ... }@inputs:
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
