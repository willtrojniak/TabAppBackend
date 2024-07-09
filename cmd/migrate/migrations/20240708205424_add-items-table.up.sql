CREATE TABLE IF NOT EXISTS items (
  id UUID NOT NULL,
  shop_id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,
  base_price REAL NOT NULL,

  PRIMARY KEY(id),
  FOREIGN KEY(shop_id) REFERENCES shops(id) ON DELETE CASCADE,
  UNIQUE(shop_id, name)
)
