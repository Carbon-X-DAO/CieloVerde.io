{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    fu.url = "github:numtide/flake-utils/bba5dcc8e0b20ab664967ad83d24d64cb64ec4f4";
  };

  outputs = self:
    self.fu.lib.eachDefaultSystem (system:
      let
        pkgs = self.nixpkgs.legacyPackages.${system};
      in
      rec {
        packages = self.fu.lib.flattenTree {
          frontend =
            let
              src = ./.;
              nodeDependencies = (pkgs.callPackage ./default.nix { }).shell.nodeDependencies;
            in
            pkgs.stdenv.mkDerivation {
              inherit src;

              version = "0.0.1";
              name = "qrinvite-frontend";

              buildPhase = ''
                mkdir $out
                cp -r ${nodeDependencies}/lib/node_modules ./node_modules;
                chmod -R u+w node_modules/;
                export PATH="${nodeDependencies}/bin:$PATH";

                ${pkgs.yarn}/bin/yarn run build;
                #${pkgs.nodePackages.npm}/bin/npm run build;
              '';


              installPhase = ''
                cp -r build/* $out/;
              '';
            };
        };

        defaultPackage = packages.frontend;
      }
    );
}
