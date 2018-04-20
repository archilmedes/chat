/* DB Setup */
DROP DATABASE IF EXISTS otrmessengertest;
CREATE DATABASE IF NOT EXISTS otrmessengertest;
USE otrmessengertest;

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

/* Inserting test data */
/* Users */
INSERT INTO users VALUES ("alice123", "alicepassword", "127.0.0.1");
INSERT INTO users VALUES ("bob", "Password", "123.456.789");
INSERT INTO users VALUES ("karateAMD", "pwd123", "192.168.10.123");
INSERT INTO users VALUES ("sameetandpotatoes", "iLuvMacs", "10.192.345.987");
INSERT INTO users VALUES ("archilmedes", "linuxFTW", "987.654.321");
INSERT INTO users VALUES ("andrew", "anotherPass", "888.888.888");

/* Sessions */
INSERT INTO sessions VALUES (12, "alice123", "123.456.789", "otr", "str1", "2018-04-20 14:18:05.283410");
INSERT INTO sessions VALUES (14, "alice123", "10.192.345.987", "plain", "protocol2", "2018-04-19 14:18:05.283410");
INSERT INTO sessions VALUES (34, "karateAMD", "10.192.345.987", "plain", "line3", "2018-04-20 10:18:05.283410");
INSERT INTO sessions VALUES (35, "karateAMD", "987.654.321", "otr", "serializedObject4", "2018-04-19 06:49:05.283410");
INSERT INTO sessions VALUES (64, "andrew", "10.192.345.987", "otr", "number5", "2018-04-20 21:04:16.283410");
INSERT INTO sessions VALUES (62, "andrew", "123.456.789", "otr", "part6", "2018-04-20 14:18:05.123456");
INSERT INTO sessions VALUES (32, "karateAMD", "123.456.789", "otr", "lastLine7", "2018-04-18 11:45:59.999999");

/* Messages */
INSERT INTO messages VALUES (12, "Hello World", "2017-02-01 08:20:19.123456", 1);
INSERT INTO messages VALUES (14, "Hey Sameet, its Alice <3", "2018-02-14 11:11:11.111111", 0);
INSERT INTO messages VALUES (34, "Hey Andrew, I need help with 511, when are you free?", "2018-04-10 12:30:08.222222", 1);
INSERT INTO messages VALUES (52, "lul", "2018-03-28 18:04:10.333333", 0);
INSERT INTO messages VALUES (34, "I almost made my Mac a brick", "2018-04-08 17:01:40.444444", 1);
INSERT INTO messages VALUES (42, "Why did the chicken cross the road?", "2018-04-12 07:56:00.555555", 1);
INSERT INTO messages VALUES (42, "To get to the other side?", "2018-04-12 07:59:13.666666", 0);
INSERT INTO messages VALUES (34, "When are we playing Fortnite?", "2018-04-08 17:59:02.777777", 0);