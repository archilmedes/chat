/* DB Setup */
CREATE DATABASE IF NOT EXISTS otrmessenger;
USE otrmessenger;

/* Table Setup*/
CREATE TABLE IF NOT EXISTS sessions (
  SSID INT NOT NULL PRIMARY KEY,
  user_id INT DEFAULT -1,
  friend_id INT DEFAULT -1,
  private_key varchar(10000) DEFAULT "",
  fingerprint varchar(10000) DEFAULT ""
);

CREATE TABLE IF NOT EXISTS users (
  username varchar(1000) NOT NULL,
  password varchar(1000) NOT NULL,
  ipaddress varchar(18) NOT NULL
);

CREATE TABLE IF NOT EXISTS messages (
  SSID INT NOT NULL,
  message varchar(10000) NOT NULL,
  timestamp varchar(30) NOT NULL,
  sent_or_received TINYINT NOT NULL
);