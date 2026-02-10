{
  description = "IMPHNEN Backend QR - Go backend for QR code generation";

  inputs = {
    nixpkgs.url = "github:NixOS/nixpkgs/nixpkgs-unstable";
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
        pkgs = import nixpkgs { inherit system; };

        # Development services script
        devServices = pkgs.writeShellScriptBin "dev-services" ''
          set -e

          PGDATA="$PWD/.dev/postgres"
          REDIS_DIR="$PWD/.dev/redis"

          mkdir -p "$PGDATA" "$REDIS_DIR"

          cleanup() {
            echo "Stopping services..."
            ${pkgs.postgresql_16}/bin/pg_ctl -D "$PGDATA" stop -m fast 2>/dev/null || true
            kill $REDIS_PID 2>/dev/null || true
            exit 0
          }
          trap cleanup EXIT INT TERM

          # Initialize PostgreSQL if needed
          if [ ! -f "$PGDATA/PG_VERSION" ]; then
            echo "Initializing PostgreSQL..."
            ${pkgs.postgresql_16}/bin/initdb -D "$PGDATA" --auth=trust --no-locale --encoding=UTF8
            echo "listen_addresses = '127.0.0.1'" >> "$PGDATA/postgresql.conf"
            echo "port = 5432" >> "$PGDATA/postgresql.conf"
            echo "unix_socket_directories = '$PGDATA'" >> "$PGDATA/postgresql.conf"
          fi

          # Start PostgreSQL
          echo "Starting PostgreSQL..."
          ${pkgs.postgresql_16}/bin/pg_ctl -D "$PGDATA" -l "$PGDATA/postgres.log" start
          sleep 2

          # Create database and user if not exists
          ${pkgs.postgresql_16}/bin/psql -h 127.0.0.1 -p 5432 -d postgres -c "SELECT 1 FROM pg_database WHERE datname = 'imphnen_qr'" | grep -q 1 || \
            ${pkgs.postgresql_16}/bin/psql -h 127.0.0.1 -p 5432 -d postgres -c "CREATE DATABASE imphnen_qr"

          ${pkgs.postgresql_16}/bin/psql -h 127.0.0.1 -p 5432 -d postgres -c "SELECT 1 FROM pg_roles WHERE rolname = 'imphnen_qr'" | grep -q 1 || \
            ${pkgs.postgresql_16}/bin/psql -h 127.0.0.1 -p 5432 -d postgres -c "CREATE USER imphnen_qr WITH PASSWORD 'imphnen_qr'"

          ${pkgs.postgresql_16}/bin/psql -h 127.0.0.1 -p 5432 -d postgres -c "GRANT ALL PRIVILEGES ON DATABASE imphnen_qr TO imphnen_qr"
          ${pkgs.postgresql_16}/bin/psql -h 127.0.0.1 -p 5432 -d imphnen_qr -c "GRANT ALL ON SCHEMA public TO imphnen_qr"

          # Start Redis
          echo "Starting Redis..."
          ${pkgs.redis}/bin/redis-server --bind 127.0.0.1 --port 6379 --dir "$REDIS_DIR" --daemonize yes --pidfile "$REDIS_DIR/redis.pid"
          REDIS_PID=$(cat "$REDIS_DIR/redis.pid")

          echo ""
          echo "Services running:"
          echo "  PostgreSQL: 127.0.0.1:5432 (database: imphnen_qr, user: imphnen_qr, password: imphnen_qr)"
          echo "  Redis:      127.0.0.1:6379"
          echo ""
          echo "Press Ctrl+C to stop all services"

          # Wait forever
          while true; do sleep 3600; done
        '';
      in
      {
        packages = {
          default = self.packages.${system}.imphnen-backend-qr;

          imphnen-backend-qr = pkgs.buildGoModule {
            pname = "imphnen-backend-qr";
            version = "0.1.0";
            src = ./.;

            vendorHash = "sha256-A8uCcZ/diCLzV/Ebk7DMAa389d2zAfBfKTa3dLFO5Oc=";

            # Build the server binary
            subPackages = [ "cmd/server" ];

            # Include migrations in the output
            postInstall = ''
              mkdir -p $out/share/imphnen-backend-qr
              cp -r db/migrations $out/share/imphnen-backend-qr/
            '';

            meta = with pkgs.lib; {
              description = "IMPHNEN Backend QR - QR code generation and image overlay service";
              homepage = "https://github.com/IMPHNEN/imphnen-backend-qr";
              license = licenses.mit;
              maintainers = [ ];
              mainProgram = "server";
            };
          };

          # Development services runner
          dev-services = devServices;
        };

        # Development shell with all tools
        devShells.default = pkgs.mkShell {
          buildInputs = with pkgs; [
            # Go toolchain
            go_1_24
            gopls
            gotools
            go-tools
            delve
            golangci-lint

            # Database tools
            postgresql_16
            redis

            # Migration tool
            go-migrate

            # HTTP testing
            curl
            jq
            httpie

            # Development tools
            git
            air # Live reload for Go

            # Dev services script
            devServices
          ];

          # Environment variables for local development
          env = {
            PORT = "8080";
            DATABASE_URL = "postgres://imphnen_qr:imphnen_qr@127.0.0.1:5432/imphnen_qr?sslmode=disable";
            JWT_SECRET = "dev-secret-change-in-production";
            REDIS_URL = "redis://127.0.0.1:6379";
          };

          shellHook = ''
            echo "IMPHNEN Backend QR Development Shell"
            echo "Go: $(go version)"
            echo ""
            echo "Development Services:"
            echo "  dev-services              - Start PostgreSQL + Redis"
            echo ""
            echo "Database Commands:"
            echo "  migrate -path db/migrations -database \$DATABASE_URL up    - Run migrations"
            echo "  migrate -path db/migrations -database \$DATABASE_URL down  - Rollback migrations"
            echo "  psql \$DATABASE_URL        - Connect to database"
            echo ""
            echo "Server Commands:"
            echo "  go run cmd/server/main.go  - Start the server"
            echo "  go run cmd/seeder/main.go  - Run the seeder"
            echo "  air                        - Start with live reload"
            echo "  go test ./...              - Run tests"
            echo ""
            echo "Build:"
            echo "  nix build .#imphnen-backend-qr"
            echo ""
            echo "Environment:"
            echo "  DATABASE_URL=\$DATABASE_URL"
            echo "  PORT=\$PORT"
          '';
        };
      }
    )
    // {
      # NixOS modules for deployment
      nixosModules = {
        default = self.nixosModules.backend-qr;
        backend-qr = import ./nix/modules/backend-qr.nix { inherit self; };
      };

      # Overlay for easy integration
      overlays.default = final: prev: {
        imphnen = (prev.imphnen or { }) // {
          backend-qr = self.packages.${final.system}.imphnen-backend-qr;
        };
      };
    };
}
