-- DROP INDEX idx_dishes;
CREATE UNIQUE INDEX unique_dishes ON dishes(restaurant_id, name, image_id) WHERE archived_at IS NULL ;
