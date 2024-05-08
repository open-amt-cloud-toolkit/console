CREATE TABLE IF NOT EXISTS devices
(
    guid uuid NOT NULL,
    tags text COLLATE pg_catalog."default",
    hostname character varying(256) COLLATE pg_catalog."default",
    mpsinstance text COLLATE pg_catalog."default",
    connectionstatus boolean,
    mpsusername text COLLATE pg_catalog."default",
    tenantid character varying(36) COLLATE pg_catalog."default" NOT NULL,
    friendlyname character varying(256) COLLATE pg_catalog."default",
    dnssuffix character varying(256) COLLATE pg_catalog."default",
    lastconnected timestamp with time zone,
    lastseen timestamp with time zone,
    lastdisconnected timestamp with time zone,
    deviceinfo text,
    username varchar(256),
    password varchar(256),
    usetls boolean,
    allowselfsigned boolean,
    CONSTRAINT devices_pkey PRIMARY KEY (guid, tenantid),
    CONSTRAINT device_guid UNIQUE (guid)
);

CREATE TABLE IF NOT EXISTS ciraconfigs(
  cira_config_name TEXT NOT NULL,
  mps_server_address TEXT,
  mps_port INTEGER,
  user_name TEXT,
  password TEXT,
  common_name TEXT,
  server_address_format INTEGER,
  auth_method INTEGER,
  mps_root_certificate TEXT,
  proxydetails TEXT,
  tenant_id TEXT NOT NULL,
  PRIMARY KEY (cira_config_name, tenant_id)
);

CREATE TABLE IF NOT EXISTS ieee8021xconfigs(
  profile_name TEXT,
  auth_protocol INTEGER,
  servername TEXT,
  domain TEXT,
  username TEXT,
  password TEXT,
  roaming_identity TEXT,
  active_in_s0 INTEGER, -- BOOLEAN converted to INTEGER
  pxe_timeout INTEGER,
  wired_interface INTEGER NOT NULL, -- BOOLEAN converted to INTEGER
  tenant_id TEXT NOT NULL,
  PRIMARY KEY (profile_name, tenant_id)
);

CREATE TABLE IF NOT EXISTS wirelessconfigs(
  wireless_profile_name TEXT NOT NULL,
  authentication_method INTEGER,
  encryption_method INTEGER,
  ssid TEXT,
  psk_value INTEGER,
  psk_passphrase TEXT,
  link_policy TEXT,
  creation_date TEXT, -- TIMESTAMP is usually represented as TEXT in SQLite
  created_by TEXT,
  tenant_id TEXT NOT NULL,
  ieee8021x_profile_name TEXT,
  FOREIGN KEY (ieee8021x_profile_name, tenant_id) REFERENCES ieee8021xconfigs(profile_name, tenant_id),
  PRIMARY KEY (wireless_profile_name, tenant_id)
);

CREATE TABLE IF NOT EXISTS profiles(
  profile_name TEXT NOT NULL,
  activation TEXT NOT NULL,
  amt_password TEXT,
  generate_random_password INTEGER NOT NULL, -- BOOLEAN converted to INTEGER
  cira_config_name TEXT,
  FOREIGN KEY (cira_config_name, tenant_id) REFERENCES ciraconfigs(cira_config_name, tenant_id),
  creation_date TEXT, -- TIMESTAMP as TEXT
  created_by TEXT,
  mebx_password TEXT,
  generate_random_mebx_password INTEGER NOT NULL, -- BOOLEAN converted to INTEGER
  tags TEXT,
  dhcp_enabled INTEGER, -- BOOLEAN converted to INTEGER
  tenant_id TEXT NOT NULL,
  tls_mode INTEGER,
  user_consent TEXT,
  ider_enabled INTEGER, -- BOOLEAN converted to INTEGER
  kvm_enabled INTEGER, -- BOOLEAN converted to INTEGER
  sol_enabled INTEGER, -- BOOLEAN converted to INTEGER
  tls_signing_authority TEXT,
  ip_sync_enabled INTEGER, -- BOOLEAN converted to INTEGER
  local_wifi_sync_enabled INTEGER, -- BOOLEAN converted to INTEGER
  ieee8021x_profile_name TEXT,
  FOREIGN KEY (ieee8021x_profile_name, tenant_id) REFERENCES ieee8021xconfigs(profile_name, tenant_id),
  PRIMARY KEY (profile_name, tenant_id)
);

CREATE TABLE IF NOT EXISTS profiles_wirelessconfigs(
  wireless_profile_name TEXT,
  profile_name TEXT,
  FOREIGN KEY (wireless_profile_name, tenant_id) REFERENCES wirelessconfigs(wireless_profile_name, tenant_id),
  FOREIGN KEY (profile_name, tenant_id) REFERENCES profiles(profile_name, tenant_id),
  priority INTEGER,
  creation_date TEXT, -- TIMESTAMP as TEXT
  created_by TEXT,
  tenant_id TEXT NOT NULL,
  PRIMARY KEY (wireless_profile_name, profile_name, priority, tenant_id)
);

CREATE TABLE IF NOT EXISTS domains(
  name TEXT NOT NULL,
  domain_suffix TEXT NOT NULL,
  provisioning_cert TEXT,
  provisioning_cert_storage_format TEXT,
  provisioning_cert_key TEXT,
  creation_date TEXT, -- TIMESTAMP as TEXT
  created_by TEXT,
  tenant_id TEXT NOT NULL,
  CONSTRAINT domainsuffix UNIQUE (domain_suffix, tenant_id),
  PRIMARY KEY (name, tenant_id)
);

CREATE UNIQUE INDEX lower_name_suffix_idx ON domains (LOWER(name), LOWER(domain_suffix));
