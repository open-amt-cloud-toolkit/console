/*********************************************************************
* Copyright (c) Intel Corporation 2023
* SPDX-License-Identifier: Apache-2.0
**********************************************************************/
CREATE EXTENSION IF NOT EXISTS citext;
CREATE USER rpsdb;
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
  cira_config_name citext NOT NULL,
  mps_server_address varchar(256),
  mps_port integer,
  user_name varchar(40),
  password varchar(63),
  common_name varchar(256),
  server_address_format integer,
  auth_method integer,
  mps_root_certificate text,
  proxydetails text,
  tenant_id varchar(36) NOT NULL,
  PRIMARY KEY (cira_config_name, tenant_id)
);
CREATE TABLE IF NOT EXISTS ieee8021xconfigs(
    profile_name citext,
    auth_protocol integer,
    servername VARCHAR(255),
    domain VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255),
    roaming_identity VARCHAR(255),
    active_in_s0 BOOLEAN,
    pxe_timeout integer,
    wired_interface BOOLEAN NOT NULL,
    tenant_id varchar(36) NOT NULL,
    PRIMARY KEY (profile_name, tenant_id)
);
CREATE TABLE IF NOT EXISTS wirelessconfigs(
  wireless_profile_name citext NOT NULL,
  authentication_method integer,
  encryption_method integer,
  ssid varchar(32),
  psk_value integer,
  psk_passphrase varchar(63),
  link_policy text,
  creation_date timestamp,
  created_by varchar(40),
  tenant_id varchar(36) NOT NULL, 
  ieee8021x_profile_name citext,
  FOREIGN KEY (ieee8021x_profile_name,tenant_id)  REFERENCES ieee8021xconfigs(profile_name,tenant_id),
  PRIMARY KEY (wireless_profile_name, tenant_id)
);
CREATE TABLE IF NOT EXISTS ieee8021xconfigs(
    profile_name citext,
    auth_protocol integer,
    servername VARCHAR(255),
    domain VARCHAR(255),
    username VARCHAR(255),
    password VARCHAR(255),
    roaming_identity VARCHAR(255),
    active_in_s0 BOOLEAN,
    pxe_timeout integer,
    wired_interface BOOLEAN NOT NULL,
    tenant_id varchar(36) NOT NULL,
    PRIMARY KEY (profile_name, tenant_id)
);
CREATE TABLE IF NOT EXISTS profiles(
  profile_name citext NOT NULL,
  activation varchar(20) NOT NULL,
  amt_password varchar(40),
  generate_random_password BOOLEAN NOT NULL,
  cira_config_name citext,
  FOREIGN KEY (cira_config_name,tenant_id)  REFERENCES ciraconfigs(cira_config_name,tenant_id),
  creation_date timestamp,
  created_by varchar(40),
  mebx_password varchar(40),
  generate_random_mebx_password BOOLEAN NOT NULL,
  tags text,
  dhcp_enabled BOOLEAN,
  tenant_id varchar(36) NOT NULL,
  tls_mode integer NULL,
  user_consent varchar(7) NULL,
  ider_enabled BOOLEAN NULL,
  kvm_enabled BOOLEAN NULL,
  sol_enabled BOOLEAN NULL,
  tls_signing_authority varchar(40) NULL,
  ip_sync_enabled BOOLEAN NULL,
  local_wifi_sync_enabled BOOLEAN NULL,
  ieee8021x_profile_name citext,
  FOREIGN KEY (ieee8021x_profile_name,tenant_id)  REFERENCES ieee8021xconfigs(profile_name,tenant_id),
  PRIMARY KEY (profile_name, tenant_id)
);
CREATE TABLE IF NOT EXISTS profiles_wirelessconfigs(
  wireless_profile_name citext,
  profile_name citext,
  FOREIGN KEY (wireless_profile_name,tenant_id)  REFERENCES wirelessconfigs(wireless_profile_name,tenant_id),
  FOREIGN KEY (profile_name,tenant_id)  REFERENCES profiles(profile_name,tenant_id),
  priority integer,
  creation_date timestamp,
  created_by varchar(40),
  tenant_id varchar(36) NOT NULL,
  PRIMARY KEY (wireless_profile_name, profile_name, priority, tenant_id)
);
CREATE TABLE IF NOT EXISTS domains(
  name citext NOT NULL,
  domain_suffix citext NOT NULL,
  provisioning_cert text,
  provisioning_cert_storage_format varchar(40),
  provisioning_cert_key text,
  creation_date timestamp,
  created_by varchar(40),
  tenant_id varchar(36) NOT NULL,
  CONSTRAINT domainsuffix UNIQUE(domain_suffix,tenant_id),
  PRIMARY KEY (name, tenant_id)
);
CREATE UNIQUE INDEX lower_name_suffix_idx ON domains ((lower(name)), (lower(domain_suffix)));
