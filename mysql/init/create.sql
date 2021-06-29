CREATE DATABASE rest_crud;
USE rest_crud;
DROP TABLE IF EXISTS recipes;
CREATE TABLE IF NOT EXISTS recipes (
  id integer PRIMARY KEY AUTO_INCREMENT,
  title varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  making_time varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  serves varchar(100) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  ingredients varchar(300) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL,
  cost integer NOT NULL,
  created_at datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at datetime on update CURRENT_TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT INTO recipes (
    id,
    title,
    making_time,
    serves,
    ingredients,
    cost,
    created_at,
    updated_at
  )
VALUES (
    1,
    'チキンカレー',
    '45分',
    '4人',
    '玉ねぎ,肉,スパイス',
    1000,
    '2016-01-10 12:10:12',
    '2016-01-10 12:10:12'
  );
INSERT INTO recipes (
    id,
    title,
    making_time,
    serves,
    ingredients,
    cost,
    created_at,
    updated_at
  )
VALUES (
    2,
    'オムライス',
    '30分',
    '2人',
    '玉ねぎ,卵,スパイス,醤油',
    700,
    '2016-01-11 13:10:12',
    '2016-01-11 13:10:12'
  );