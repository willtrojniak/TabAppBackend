CREATE TABLE IF NOT EXISTS item_addons (
  shop_id UUID NOT NULL,
  item_id UUID NOT NULL,
  addon_id UUID NOT NULL,
  index SMALLINT NOT NULL,

  PRIMARY KEY(shop_id, item_id, addon_id),
  FOREIGN KEY(shop_id, item_id) REFERENCES items(shop_id, id) ON DELETE CASCADE,
  FOREIGN KEY(shop_id, addon_id) REFERENCES items(shop_id, id) ON DELETE CASCADE
);
