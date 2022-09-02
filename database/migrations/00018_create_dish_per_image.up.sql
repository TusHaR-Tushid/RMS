CREATE TABLE IF NOT EXISTS dish_per_image(
                                                  id SERIAL PRIMARY KEY ,
                                                  dish_id INTEGER REFERENCES dishes(id),
                                                  image_id       INTEGER REFERENCES images(id)
);

ALTER TABLE dishes DROP COLUMN image_id;