CREATE TABLE IF NOT EXISTS order_items (
  shop_id INT NOT NULL,
  tab_id INT NOT NULL,
  bill_id INT NOT NULL,
  item_id INT NOT NULL,
  quantity INT NOT NULL DEFAULT 0,

  PRIMARY KEY(shop_id, tab_id, bill_id, item_id),
  FOREIGN KEY(shop_id, tab_id, bill_id) REFERENCES tab_bills(shop_id, tab_id, id),
  FOREIGN KEY(shop_id, item_id) REFERENCES items(shop_id, id),
  CHECK ( quantity >= 0 )
);

CREATE TABLE IF NOT EXISTS order_variants (
  shop_id INT NOT NULL,
  tab_id INT NOT NULL,
  bill_id INT NOT NULL,
  item_id INT NOT NULL,
  variant_id INT NOT NULL,
  quantity INT NOT NULL DEFAULT 0,

  PRIMARY KEY(shop_id, tab_id, bill_id, item_id, variant_id),
  FOREIGN KEY(shop_id, tab_id, bill_id) REFERENCES tab_bills(shop_id, tab_id, id),
  FOREIGN KEY(shop_id, item_id, variant_id) REFERENCES item_variants(shop_id, item_id, id),
  CHECK ( quantity >= 0 )
);
