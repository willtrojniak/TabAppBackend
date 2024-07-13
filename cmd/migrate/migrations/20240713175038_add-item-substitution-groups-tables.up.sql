CREATE TABLE IF NOT EXISTS item_substitution_groups (
  shop_id UUID NOT NULL,
  id UUID NOT NULL,
  name VARCHAR(255) NOT NULL,

  PRIMARY KEY(shop_id, id),
  FOREIGN KEY(shop_id) REFERENCES shops(id) ON DELETE CASCADE,
  UNIQUE(shop_id, name)
);

CREATE TABLE IF NOT EXISTS item_substitution_groups_to_items (
  shop_id UUID NOT NULL,
  substitution_group_id UUID NOT NULL,
  item_id UUID NOT NULL,
  index SMALLINT NOT NULL,

  PRIMARY KEY(shop_id, substitution_group_id, item_id),
  FOREIGN KEY(shop_id, substitution_group_id) REFERENCES item_substitution_groups(shop_id, id) ON DELETE CASCADE,
  FOREIGN KEY(shop_id, item_id) REFERENCES items(shop_id, id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS items_to_item_substitution_groups (
  shop_id UUID NOT NULL,
  substitution_group_id UUID NOT NULL,
  item_id UUID NOT NULL,
  index SMALLINT NOT NULL,

  PRIMARY KEY(shop_id, substitution_group_id, item_id),
  FOREIGN KEY(shop_id, substitution_group_id) REFERENCES item_substitution_groups(shop_id, id) ON DELETE CASCADE,
  FOREIGN KEY(shop_id, item_id) REFERENCES items(shop_id, id) ON DELETE CASCADE
);
