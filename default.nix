{ pkgs,
  buildGoApplication,
  ...
}:

buildGoApplication {
  pname = "assume-role";
  version = "0.1";
  src = ./.;
  modules = ./gomod2nix.toml;
}
