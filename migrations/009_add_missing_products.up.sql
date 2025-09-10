-- Add the 2 missing products to complete the set of 10 new products

-- Add Smart Tracking Card
INSERT INTO products (id, vendor_id, name, description, price, compare_price, sku, category, tags, images, stock, status, featured, weight, dimensions, seo, numeric_id) VALUES
('550e8400-e29b-41d4-a716-446655440028', '550e8400-e29b-41d4-a716-446655440002', 'Smart Tracking Card', 
'Ultra-slim Bluetooth tracking card that fits perfectly in your wallet. Works with Apple Find My network for easy location tracking of your valuables.', 
299.00, 349.00, 'SMART-TRACK-CARD', 'electronics', 
ARRAY['tracking', 'bluetooth', 'wallet', 'find-my', 'smart'], 
ARRAY['https://d10qehs4k3bdf9.cloudfront.net/smart-tracking-card.jpg'], 
35, 'active', false, 0.020,
JSON_BUILD_OBJECT('length', 8.5, 'width', 5.4, 'height', 0.2),
JSON_BUILD_OBJECT('title', 'Smart Tracking Card - Bluetooth Wallet Tracker', 'description', 'Ultra-slim tracking card for wallets with Apple Find My support', 'keywords', ARRAY['tracking', 'card', 'wallet', 'bluetooth', 'find-my']), 20),

-- Add MacBook Air M3 Power Adapter and Cable
('550e8400-e29b-41d4-a716-446655440029', '550e8400-e29b-41d4-a716-446655440002', 'MacBook Air M3 Power Adapter and Cable', 
'Official Apple 70W USB-C Power Adapter with MagSafe 3 charging cable for MacBook Air M3. Fast charging capability and magnetic connection for safety.', 
890.00, 990.00, 'MBA-M3-ADAPTER-CABLE', 'computers', 
ARRAY['apple', 'macbook', 'charger', 'magsafe', 'adapter', 'cable'], 
ARRAY['https://d10qehs4k3bdf9.cloudfront.net/macbookair-adaptor-and-cable.png'], 
28, 'active', true, 0.285,
JSON_BUILD_OBJECT('length', 10.5, 'width', 7.5, 'height', 3.2),
JSON_BUILD_OBJECT('title', 'MacBook Air M3 70W Power Adapter with MagSafe 3 Cable', 'description', 'Official Apple power adapter and MagSafe 3 cable for MacBook Air M3', 'keywords', ARRAY['macbook', 'air', 'm3', 'charger', 'adapter', 'magsafe', 'cable']), 21);