CREATE TABLE IF NOT EXISTS shop_slack_connections (
  shop_id INT NOT NULL,
  slack_access_token VARCHAR(255),
  daily_update_slack_channel VARCHAR(255) NOT NULL DEFAULT '',
  tab_request_slack_channel VARCHAR(255) NOT NULL DEFAULT '',
  tab_bill_receipt_slack_channel VARCHAR(255) NOT NULL DEFAULT '',


  PRIMARY KEY(shop_id),
  FOREIGN KEY(shop_id) REFERENCES shops(id) ON DELETE CASCADE
);
