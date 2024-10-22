CREATE TABLE IF NOT EXISTS shop_users(
  shop_id INT NOT NULL,
  user_id VARCHAR(255) NOT NULL,
  roles INT not NULL,
  confirmed BOOLEAN not NULL,
  updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),

  PRIMARY KEY(shop_id, user_id),
  FOREIGN KEY(shop_id) REFERENCES shops(id) ON DELETE CASCADE,
  FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);
