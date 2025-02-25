{
  description = "Go 1.23 Environment";

  inputs.nixpkgs.url = "github:NixOS/nixpkgs";

  outputs = { self, nixpkgs }: 
  let
    pkgs = import nixpkgs { system = "aarch64-darwin"; };
  in {
    devShells.aarch64-darwin.default = pkgs.mkShell {
      buildInputs = [
        pkgs.go_1_23
      ];
    };
  };
}