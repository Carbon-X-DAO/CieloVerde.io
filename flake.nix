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
    packages = {
      frontend = let
        src = ./frontend;
        #nodeDependencies = (pkgs.callPackage "${src}/default.nix" {}).shell.nodeDependencies;
        nodeDependencies = (pkgs.callPackage ./frontend/default.nix {}).shell.nodeDependencies;
      in
        pkgs.stdenv.mkDerivation {
        inherit src;

        version = "0.0.1";
        name = "qrinvite-frontend";

        sha256 = pkgs.lib.fakeSha256;

        buildPhase = ''
          mkdir $out
          cp -r ${nodeDependencies}/lib/node_modules ./node_modules;
          chmod -R u+w node_modules/;
          export PATH="${nodeDependencies}/bin:$PATH";

          ${pkgs.yarn}/bin/yarn run build;

        '';


        installPhase = ''
          cp -r build/* $out/;
        '';
      };

      webserver = pkgs.buildGoModule {
        name = "webserver";

        buildInputs = let
          np = pkgs.nodePackages_latest;
        in
        [
          np.create-react-app
        ];

        src = ./.;

        buildPhase = ''
          ${pkgs.yarn}/bin/yarn run build
        '';

        vendorSha256 = "pQpattmS9VmO3ZIQUFn66az8GSmB4IvYhTTCFn6SUmo=";

      };
    };

    defaultPackage.x86_64-linux = packages.webserver;
  }
  );
}
