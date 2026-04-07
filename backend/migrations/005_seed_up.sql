WITH new_shop AS (
    INSERT INTO shops (name, created_at, updated_at)
    VALUES ('Posiflora Demo Shop', NOW(), NOW())
    RETURNING id
)
INSERT INTO orders (shop_id, number, total, customer_name, created_at, updated_at)
SELECT
    new_shop.id,
    v.number,
    v.total,
    v.customer_name,
    NOW(),
    NOW()
FROM new_shop
CROSS JOIN (VALUES
    ('A-1001', 1290,  'Алексей'),
    ('A-1002', 3450,  'Мария'),
    ('A-1003', 990,   'Дмитрий'),
    ('A-1004', 2100,  'Елена'),
    ('A-1005', 4800,  'Сергей'),
    ('A-1006', 1750,  'Ольга'),
    ('A-1007', 3200,  'Николай'),
    ('A-1008', 890,   'Татьяна'),
    ('A-1009', 5600,  'Андрей'),
    ('A-1010', 2350,  'Юлия')
) AS v(number, total, customer_name);
