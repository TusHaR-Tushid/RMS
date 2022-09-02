ALTER TABLE dishes DROP COLUMN image_id;
ALTER TABLE restaurants DROP COLUMN image_id;
ALTER TABLE images  ADD COLUMN image_id int ;