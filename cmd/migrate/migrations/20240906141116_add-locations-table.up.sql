CREATE TABLE IF NOT EXISTS locations (
  id SERIAL NOT NULL,
  shop_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,

  PRIMARY KEY(shop_id, id),
  FOREIGN KEY(shop_id) REFERENCES shops(id) ON DELETE CASCADE,
  UNIQUE(shop_id, name)
);

CREATE TABLE IF NOT EXISTS tab_locations (
  shop_id INT NOT NULL,
  tab_id INT NOT NULL,
  location_id INT NOT NULL,

  PRIMARY KEY(shop_id, tab_id, location_id),
  FOREIGN KEY(shop_id, location_id) REFERENCES locations(shop_id, id) ON DELETE CASCADE,
  FOREIGN KEY(shop_id, tab_id) REFERENCES tabs(shop_id, id) ON DELETE CASCADE
);
