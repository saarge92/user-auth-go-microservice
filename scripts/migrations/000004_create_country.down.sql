ALTER TABLE users DROP CONSTRAINT `player_country_id_fk`;
ALTER TABLE users DROP COLUMN country_id;

DROP TABLE countries;