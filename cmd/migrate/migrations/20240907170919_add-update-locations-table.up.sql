CREATE TABLE IF NOT EXISTS tab_update_locations (
  shop_id INT NOT NULL,
  tab_id INT NOT NULL,
  location_id INT NOT NULL,

  PRIMARY KEY(shop_id, tab_id, location_id),
  FOREIGN KEY(shop_id, location_id) REFERENCES locations(shop_id, id) ON DELETE CASCADE,
  FOREIGN KEY(shop_id, tab_id) REFERENCES tabs(shop_id, id) ON DELETE CASCADE
);
