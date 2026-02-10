# NixOS Module for IMPHNEN Backend QR
{ self }:

{
  config,
  lib,
  pkgs,
  ...
}:

let
  cfg = config.services.imphnen-backend-qr;

  # Build a minimal go-migrate without problematic drivers
  goMigrate = pkgs.buildGoModule rec {
    pname = "migrate";
    version = "4.18.1";

    src = pkgs.fetchFromGitHub {
      owner = "golang-migrate";
      repo = "migrate";
      rev = "v${version}";
      hash = "sha256-bLo3kihkPpuU+gzWNFN/bxOLe6z+ILEHxkyZ9XB3eek=";
    };

    vendorHash = "sha256-Wu3if5gNAEuD4YwaZfjC+YQK2lPyb1FMuaoWRlKyJYo=";

    subPackages = [ "cmd/migrate" ];

    tags = [ "postgres" ];

    meta = {
      description = "Database migrations";
      mainProgram = "migrate";
    };
  };
in
{
  options.services.imphnen-backend-qr = {
    enable = lib.mkEnableOption "IMPHNEN Backend QR Service";

    port = lib.mkOption {
      type = lib.types.port;
      default = 8080;
      description = "Port to run the backend server on";
    };

    package = lib.mkOption {
      type = lib.types.package;
      default = self.packages.${pkgs.system}.imphnen-backend-qr;
      description = "The backend-qr package to use";
    };

    openFirewall = lib.mkOption {
      type = lib.types.bool;
      default = false;
      description = "Open firewall for the backend port";
    };

    environmentFile = lib.mkOption {
      type = lib.types.nullOr lib.types.path;
      default = null;
      description = "Environment file containing secrets (DATABASE_URL, JWT_SECRET, etc.)";
    };

    database = {
      host = lib.mkOption {
        type = lib.types.str;
        default = "/run/postgresql";
        description = "PostgreSQL database host (use socket path for local)";
      };

      port = lib.mkOption {
        type = lib.types.port;
        default = 5432;
        description = "PostgreSQL database port";
      };

      name = lib.mkOption {
        type = lib.types.str;
        default = "imphnen_qr";
        description = "PostgreSQL database name";
      };

      user = lib.mkOption {
        type = lib.types.str;
        default = "imphnen_qr";
        description = "PostgreSQL database user";
      };

      createLocally = lib.mkOption {
        type = lib.types.bool;
        default = true;
        description = "Create the database locally using PostgreSQL service";
      };
    };

    runMigrations = lib.mkOption {
      type = lib.types.bool;
      default = true;
      description = "Run database migrations on startup";
    };
  };

  config = lib.mkIf cfg.enable {
    # PostgreSQL configuration if createLocally is enabled
    services.postgresql = lib.mkIf cfg.database.createLocally {
      enable = true;
      ensureDatabases = [ cfg.database.name ];
      ensureUsers = [
        {
          name = cfg.database.user;
          ensureDBOwnership = true;
        }
      ];
      # Enable TCP connections with md5 auth
      enableTCPIP = true;
      authentication = lib.mkForce ''
        # Local connections
        local all all trust
        # IPv4 local connections
        host all all 127.0.0.1/32 md5
        # IPv6 local connections
        host all all ::1/128 md5
      '';
      # Set password for the user after creation
      initialScript = pkgs.writeText "init-sql" ''
        ALTER USER ${cfg.database.user} WITH PASSWORD '${cfg.database.user}';
        GRANT ALL PRIVILEGES ON DATABASE ${cfg.database.name} TO ${cfg.database.user};
      '';
    };

    # Migration service (runs before main service)
    systemd.services.imphnen-backend-qr-migrate = lib.mkIf cfg.runMigrations {
      description = "IMPHNEN Backend QR Database Migrations";
      wantedBy = [ "multi-user.target" ];
      after = [ "network.target" ] ++ lib.optional cfg.database.createLocally "postgresql.service";
      requires = lib.optional cfg.database.createLocally "postgresql.service";
      before = [ "imphnen-backend-qr.service" ];

      serviceConfig = {
        Type = "oneshot";
        RemainAfterExit = true;
        ExecStart = "${goMigrate}/bin/migrate -path ${cfg.package}/share/imphnen-backend-qr/migrations -database \"postgres://${cfg.database.user}:${cfg.database.user}@127.0.0.1:${toString cfg.database.port}/${cfg.database.name}?sslmode=disable\" up";

        # Hardening
        NoNewPrivileges = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        PrivateTmp = true;
      };
    };

    # Main backend service
    systemd.services.imphnen-backend-qr = {
      description = "IMPHNEN Backend QR Service";
      wantedBy = [ "multi-user.target" ];
      after =
        [ "network.target" ]
        ++ lib.optional cfg.database.createLocally "postgresql.service"
        ++ lib.optional cfg.runMigrations "imphnen-backend-qr-migrate.service";
      requires = lib.optional cfg.database.createLocally "postgresql.service";

      environment = {
        PORT = toString cfg.port;
      };

      serviceConfig = {
        Type = "simple";
        ExecStart = "${cfg.package}/bin/server";
        Restart = "on-failure";
        RestartSec = "5s";

        # Hardening
        DynamicUser = true;
        NoNewPrivileges = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        PrivateTmp = true;
        ProtectKernelTunables = true;
        ProtectKernelModules = true;
        ProtectControlGroups = true;
        RestrictAddressFamilies = [
          "AF_INET"
          "AF_INET6"
          "AF_UNIX"
        ];
        RestrictNamespaces = true;
        LockPersonality = true;
        RestrictRealtime = true;
        RestrictSUIDSGID = true;
      } // lib.optionalAttrs (cfg.environmentFile != null) {
        EnvironmentFile = cfg.environmentFile;
      };
    };

    networking.firewall.allowedTCPPorts = lib.mkIf cfg.openFirewall [ cfg.port ];
  };
}
