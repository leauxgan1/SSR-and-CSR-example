{
  description = "A development environment with Golang";

  inputs = {
    nixpkgs.url = "github:nixos/nixpkgs?ref=nixos-unstable";
    flake-utils.url = "github:numtide/flake-utils";
  };

  outputs = { self, nixpkgs, flake-utils, ...}:  
    flake-utils.lib.eachDefaultSystem (system: 
			let
				pkgs = import nixpkgs { inherit system; };
				go = pkgs.go_1_24;
				buildGoApp = {name, dir} : pkgs.buildGoModule {
					inherit name;
					src = ./.;
					subPackages = [dir];
					vendorHash = null;
					CGO_ENABLED = 0;
					buildPhase = ''
            cd ${dir}
            go build -o $out/bin/${name}
          '';
				};
			in {
				packages = {
					csr = buildGoApp {name = "csr"; dir = "vanilla_csr";};
					ssr = buildGoApp {name = "ssr"; dir = "htmx_ssr";};
				};
				apps = {
					csr = { type = "app"; program = "${self.packages.${system}.csr}/bin/csr";};
					ssr = { type = "app"; program = "${self.packages.${system}.ssr}/bin/ssr";};
				};
				devShells.default = pkgs.mkShell {
					packages = [ go pkgs.gopls];
					shellHook = ''
						echo "Run nix run .#csr or nix run .#ssr in their respective directories to try out each server. "
					'';
				};
  });
}
