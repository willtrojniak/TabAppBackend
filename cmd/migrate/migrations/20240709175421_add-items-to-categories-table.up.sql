CREATE TABLE IF NOT EXISTS items_to_categories(
  shop_id INT NOT NULL,
  item_id INT NOT NULL,
  item_category_id INT NOT NULL,
  index SMALLINT NOT NULL,

  PRIMARY KEY(shop_id, item_id, item_category_id),
  FOREIGN KEY(shop_id, item_id) REFERENCES items(shop_id, id) ON DELETE CASCADE,
  FOREIGN KEY(shop_id, item_category_id) REFERENCES item_categories(shop_id, id) ON DELETE CASCADE
)
