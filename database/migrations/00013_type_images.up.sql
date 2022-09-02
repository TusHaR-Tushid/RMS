create type image_type as enum('restaurant', 'dish');

ALTER TABLE images ADD COLUMN image_type image_type;