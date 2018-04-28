/* DB Setup */
CREATE DATABASE IF NOT EXISTS otrmessenger;
USE otrmessenger;

/* Table Setup*/
CREATE TABLE IF NOT EXISTS sessions (
  SSID BIGINT UNSIGNED NOT NULL PRIMARY KEY,
  username varchar(30) NOT NULL,
  friend_display_name varchar(30) NOT NULL,
  protocol_type varchar(8) NOT NULL,
  protocol varbinary(10000) NOT NULL,
  session_timestamp timestamp(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  username varchar(30) NOT NULL,
  password varchar(256) NOT NULL,
  ipaddress varchar(18) NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
  SSID BIGINT UNSIGNED NOT NULL,
  message varchar(10000) NOT NULL,
  message_timestamp timestamp(6) NOT NULL,
  sent_or_received TINYINT NOT NULL
);

CREATE TABLE IF NOT EXISTS friends (
  username varchar(30) NOT NULL,
  friend_display_name varchar(30) NOT NULL,
  friend_mac_address varchar(18) NOT NULL,
  friend_ip_address varchar(18) NOT NULL,
  friend_username varchar(30) NOT NULL
);