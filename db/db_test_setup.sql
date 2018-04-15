DROP DATABASE IF EXISTS otrmessengertest;
CREATE DATABASE IF NOT EXISTS otrmessengertest;
USE otrmessengertest;
CREATE TABLE IF NOT EXISTS sessions (
  SSID INT NOT NULL PRIMARY KEY,
  user_id INT NOT NULL,
  friend_id INT NOT NULL,
  private_key varchar(10000) NOT NULL,
  fingerprint varchar(10000) NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
  username varchar(1000) NOT NULL,
  password varchar(1000) NOT NULL,
  ipaddress varchar(18) NOT NULL
);
CREATE TABLE IF NOT EXISTS conversations (
  SSID INT NOT NULL,
  message varchar(10000) NOT NULL,
  timestamp varchar(30) NOT NULL,
  sent_or_received TINYINT NOT NULL
);