CREATE TYPE address_type AS ENUM ('billing', 'shipping', 'both');
CREATE TABLE IF NOT EXISTS addresses (
  id SERIAL PRIMARY KEY,
  street_line_1 TEXT NOT NULL,
  street_line_2 TEXT,
  city TEXT NOT NULL,
  state TEXT NOT NULL,
  zip_code TEXT NOT NULL,
  type address_type NOT NULL,
  is_default BOOLEAN DEFAULT false,
  customer_id INTEGER NOT NULL REFERENCES customers(id)
);
