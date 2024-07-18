CREATE TYPE payment_method AS ENUM ('in person', 'chartstring');

CREATE TABLE IF NOT EXISTS payment_methods (
  shop_id INT NOT NULL,
  method  payment_method NOT NULL,

  PRIMARY KEY(shop_id, method),
  FOREIGN KEY(shop_id) REFERENCES shops(id) ON DELETE CASCADE
);
