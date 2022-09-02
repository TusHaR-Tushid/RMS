ALTER TABLE dishes DROP COLUMN dish_image;

ALTER TABLE dishes ADD COLUMN image_id int;

ALTER TABLE restaurants ADD COLUMN image_id int;

CREATE TABLE IF NOT EXISTS images(
                                           id SERIAL PRIMARY KEY ,
                                           url TEXT,
                                           uploaded_by int REFERENCES users(id),
                                           created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
                                           updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
                                           archived_at TIMESTAMP WITH TIME ZONE
);