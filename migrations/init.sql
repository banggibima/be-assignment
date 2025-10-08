-- Membuat tabel products untuk menyimpan data produk
CREATE TABLE IF NOT EXISTS products (
  id TEXT PRIMARY KEY,
  name TEXT NOT NULL,
  stock INTEGER NOT NULL,
  price INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membuat tabel orders untuk menyimpan data pesanan
CREATE TABLE IF NOT EXISTS orders (
  id TEXT PRIMARY KEY,
  product_id TEXT NOT NULL REFERENCES products(id) ON DELETE CASCADE,
  buyer_id TEXT NOT NULL,
  quantity INTEGER NOT NULL,
  total_price INTEGER NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membuat tabel transactions untuk menyimpan data transaksi
CREATE TABLE IF NOT EXISTS transactions (
  id TEXT PRIMARY KEY,
  order_id TEXT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
  merchant_id TEXT NOT NULL,
  amount INTEGER NOT NULL,
  fee INTEGER NOT NULL,
  status TEXT NOT NULL,
  paid_at DATE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Membuat tabel settlements untuk menyimpan data settlement harian per merchant
CREATE TABLE IF NOT EXISTS settlements (
  id TEXT PRIMARY KEY,
  merchant_id TEXT NOT NULL,
  date DATE NOT NULL,
  gross_amount INTEGER NOT NULL,
  fee_amount INTEGER NOT NULL,
  net_amount INTEGER NOT NULL,
  txn_count INTEGER NOT NULL,
  unique_run_id TEXT NOT NULL,
  generated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  UNIQUE (merchant_id, date)
);

-- Membuat tabel jobs untuk menyimpan data job processing
CREATE TABLE IF NOT EXISTS jobs (
  id TEXT PRIMARY KEY,
  job_id TEXT NOT NULL,
  status TEXT NOT NULL,
  processed INTEGER NOT NULL,
  total INTEGER NOT NULL,
  progress INTEGER NOT NULL,
  result_path TEXT,
  from_date DATE NOT NULL,    
  to_date DATE NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);