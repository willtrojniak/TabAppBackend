CREATE TABLE IF NOT EXISTS items (
  id SERIAL NOT NULL,
  shop_id INT NOT NULL,
  name VARCHAR(255) NOT NULL,
  base_price REAL NOT NULL,

  PRIMARY KEY(shop_id, id),
  FOREIGN KEY(shop_id) REFERENCES shops(id) ON DELETE CASCADE,
  UNIQUE(shop_id, name)
)
