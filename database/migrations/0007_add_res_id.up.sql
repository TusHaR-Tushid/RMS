ALTER TABLE dishes ADD COLUMN restaurant_id int REFERENCES restaurants(id);
CREATE UNIQUE INDEX idx_dishes ON dishes(restaurant_id, name, archived_at);

