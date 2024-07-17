CREATE TABLE IF NOT EXISTS tab_updates (
  shop_id UUID NOT NULL,
  tab_id INT NOT NULL,
  payment_method payment_method NOT NULL,
  organization VARCHAR(255) NOT NULL,
  display_name VARCHAR(255) NOT NULL,
  start_date DATE NOT NULL,
  end_date DATE NOT NULL,
  daily_start_time TIME(0) NOT NULL,
  daily_end_time TIME(0) NOT NULL,
  active_days_of_wk SMALLINT NOT NULL,
  dollar_limit_per_order REAL NOT NULL,
  verification_method tab_verification_method NOT NULL,
  payment_details VARCHAR(255) NOT NULL,
  billing_interval_days SMALLINT NOT NULL,

  PRIMARY KEY(shop_id, tab_id),
  FOREIGN KEY(shop_id, tab_id) REFERENCES tabs(shop_id, id) ON DELETE CASCADE,
  FOREIGN KEY(shop_id, payment_method) REFERENCES payment_methods(shop_id, method)
);
