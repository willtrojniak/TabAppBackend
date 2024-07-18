CREATE TABLE IF NOT EXISTS item_variants (
  shop_id INT NOT NULL,
  item_id INT NOT NULL,
  id SERIAL NOT NULL,
  name VARCHAR(255) NOT NULL,
  price REAL NOT NULL,
  index SMALLINT NOT NULL,

  PRIMARY KEY(shop_id, item_id, id),
  FOREIGN KEY(shop_id, item_id) REFERENCES items(shop_id, id) ON DELETE CASCADE,
  UNIQUE(shop_id, item_id, name)
)
