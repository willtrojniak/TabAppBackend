CREATE TABLE IF NOT EXISTS shop_slack_tokens(
  shop_id INT NOT NULL,
  slack_access_token VARCHAR(255),

  PRIMARY KEY(shop_id),
  FOREIGN KEY(shop_id) REFERENCES shops(id) ON DELETE CASCADE
);
