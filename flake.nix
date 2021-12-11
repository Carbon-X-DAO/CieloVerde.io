{
  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs/nixos-unstable";
    fu.url = "github:numtide/flake-utils/bba5dcc8e0b20ab664967ad83d24d64cb64ec4f4";
    frontend.url = "git+ssh://git@github.com/Carbon-X-DAO/cieloverde.io-frontend?ref=main";
  };

  outputs = self:
    self.fu.lib.eachDefaultSystem (system:
      let
        pkgs = self.nixpkgs.legacyPackages.${system};
        webserver = pkgs.buildGoModule rec {
          pname = "webserver";
          version = "0.0.1";

          src = ./.;

          vendorSha256 = "ugVFrfT+7sHC+lajBfZfwzwRjliFQul7su9Zf08jNaI=";
        };
      in
      rec {
        packages = self.fu.lib.flattenTree {
          inherit webserver;

          site = pkgs.stdenv.mkDerivation {
            name = "site";
            version = "0.0.1";

            src = ./.;

            configurePhase = ''
              mkdir -p $out/static
              cp -r ${self.frontend.packages.${system}.frontend}/* $out/static/;
            '';

            installPhase = ''
              cp -r ${webserver}/bin/server $out/server
            '';
          };
        };

        defaultPackage = packages.site;
      }
    );
}
