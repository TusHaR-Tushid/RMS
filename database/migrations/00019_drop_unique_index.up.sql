DROP INDEX unique_dishes;
CREATE UNIQUE INDEX unique_dish ON dishes(restaurant_id, name) WHERE archived_at IS NULL ;