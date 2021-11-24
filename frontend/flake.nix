{
  description = "A very basic flake";

  inputs = {
    flake-utils.url = "github:numtide/flake-utils";
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
  };

  outputs = { self, nixpkgs, flake-utils }:
    flake-utils.lib.eachDefaultSystem (system:
      let
        nodeDependencies = (import ./overrides.nix { inherit pkgs; }).shell.nodeDependencies;
        pkgs = nixpkgs.legacyPackages.${system};
      in
      rec {
        packages = flake-utils.lib.flattenTree {
          frontend = pkgs.stdenv.mkDerivation {
            name = "test";
            src = ./.;

            # nativeBuildInputs is used mostly for cross-compiling...
            # including the below may introduce problems
            buildInputs =
              let
                p = pkgs;
              in
              [
                p.nodePackages.npm
              ];

            configurePhase = ''
              ln -s ${nodeDependencies}/lib/node_modules ./node_modules;
              export PATH="${nodeDependencies}/bin:$PATH";
            '';

            buildPhase = ''
              npm run build
            '';

            installPhase = ''
              cp -r public $out/
            '';
          };
        };

        devShell = pkgs.mkShell {
          buildInputs =
            let
              p = pkgs;
            in
            [
              p.nodePackages.npm
              p.simple-http-server
            ];

          shellHook = ''
            export NODE_PATH="${nodeDependencies}/lib/node_modules";
            ln -s "${nodeDependencies}/lib/node_modules" ./node_modules;
            export PATH="${nodeDependencies}/bin:$PATH";
          '';
        };

        defaultPackage = packages.frontend;
      }
    );
}
