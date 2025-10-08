-- Seed data awal untuk produk
INSERT INTO products (id, name, stock, price, created_at, updated_at)
VALUES
  ('product-1', 'Starter Pack', 100, 150000, NOW(), NOW())
ON CONFLICT (id) DO NOTHING;

-- Seed dummy transaksi (optional untuk settlement test)
INSERT INTO orders (id, product_id, buyer_id, quantity, total_price, created_at, updated_at)
SELECT
  'order-' || i,
  'product-1',
  'buyer-' || i,
  1,
  150000,
  NOW(),
  NOW()
FROM generate_series(1, 1000) AS s(i)
ON CONFLICT (id) DO NOTHING;

-- Misal generate 10 merchants * 100 transaksi per merchant
DO $$
DECLARE
  i INT;
  m_id TEXT;
  txn_id TEXT;
BEGIN
  FOR i IN 1..1000 LOOP
    m_id := 'merchant-' || (1 + floor(random() * 10));
    txn_id := 'txn-' || i;

    INSERT INTO transactions (
      id, order_id, merchant_id, amount, fee, status, paid_at, created_at, updated_at
    ) VALUES (
      txn_id,
      'order-' || i,
      m_id,
      (10000 + floor(random() * 50000)),
      500,
      'PAID',
      date '2025-01-01' + (floor(random() * 30)) * interval '1 day',
      NOW(),
      NOW()
    );
  END LOOP;
END $$;
