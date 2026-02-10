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
        default = "localhost";
        description = "PostgreSQL database host";
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
        ExecStart = "${pkgs.go-migrate}/bin/migrate -path ${cfg.package}/share/imphnen-backend-qr/migrations -database \"postgres://${cfg.database.user}@${cfg.database.host}:${toString cfg.database.port}/${cfg.database.name}?sslmode=disable\" up";

        # Run as postgres user for local database
        User = lib.mkIf cfg.database.createLocally "postgres";

        # Hardening
        NoNewPrivileges = true;
        ProtectSystem = "strict";
        ProtectHome = true;
        PrivateTmp = true;
      } // lib.optionalAttrs (cfg.environmentFile != null) {
        EnvironmentFile = cfg.environmentFile;
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
        # Database URL will be set via environmentFile for security
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
