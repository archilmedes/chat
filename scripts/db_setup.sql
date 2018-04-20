/* DB Setup */
CREATE DATABASE IF NOT EXISTS otrmessenger;
USE otrmessenger;

/* Table Setup*/
CREATE TABLE IF NOT EXISTS sessions (
  SSID BIGINT NOT NULL PRIMARY KEY,
  username varchar(100) NOT NULL,
  friend_mac varchar(18) NOT NULL,
  protocol_type varchar(8) NOT NULL,
  protocol varchar(1000) NOT NULL,
  session_timestamp timestamp(6) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
  username varchar(100) NOT NULL,
  password varchar(1000) NOT NULL,
  ipaddress varchar(18) NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
  SSID INT NOT NULL,
  message varchar(10000) NOT NULL,
  message_timestamp timestamp(6) NOT NULL,
  sent_or_received TINYINT NOT NULL
);