USE loyalty_program;

-- Inserting sample data into Stores
INSERT INTO stores (store_uuid, name, region) VALUES
(UUID(), 'Store A', 'North America'),
(UUID(), 'Store B', 'Europe'),
(UUID(), 'Store C', 'Europe');



-- Inserting sample data into Users
INSERT INTO users (user_uuid, name, email_address, phone_number, creation_date, invite_code) VALUES
(UUID(), 'John Doe', 'johndoe@example.com', '1234567890', NOW(), UUID()),
(UUID(), 'Jane Smith', 'janesmith@example.com', '0987654321', NOW(), UUID()),
(UUID(), 'Dave Smith', 'davesmith@example.com', '0987654322', NOW(), UUID()),
(UUID(), 'Mary Smith', 'marysmith@example.com', '0987654323', NOW(), UUID());


-- Inserting sample data into Accounts
INSERT INTO accounts (account_uuid, owner_id, creation_date, region) VALUES
(UUID(), (SELECT user_uuid FROM users WHERE name = 'John Doe'), NOW(), 'Europe'),
(UUID(), (SELECT user_uuid FROM users WHERE name = 'Jane Smith'), NOW(), 'Europe'),
(UUID(), (SELECT user_uuid FROM users WHERE name = 'Dave Smith'), NOW(), 'Europe'),
(UUID(), (SELECT user_uuid FROM users WHERE name = 'Mary Smith'), NOW(), 'Europe');



UPDATE users u
JOIN accounts a ON u.name = a.owner_id
SET u.account_uuid = a.account_uuid;



-- Inserting sample data into Transactions
INSERT INTO transactions (transaction_uuid, account_uuid, user_uuid, amount, date, store_uuid, points_earned) VALUES
(UUID(), (SELECT account_uuid FROM accounts WHERE owner_id = (SELECT user_uuid FROM users WHERE name = 'John Doe')), (SELECT user_uuid FROM users WHERE name = 'John Doe'), 100.00, NOW(), (SELECT store_uuid FROM stores WHERE name = 'Store B'), 1),
(UUID(), (SELECT account_uuid FROM accounts WHERE owner_id = (SELECT user_uuid FROM users WHERE name = 'Jane Smith')), (SELECT user_uuid FROM users WHERE name = 'Jane Smith'), 150.00, NOW(), (SELECT store_uuid FROM stores WHERE name = 'Store C'), 1);


UPDATE accounts
SET points_balance = points_balance + 5
WHERE account_uuid = (
    SELECT a.account_uuid FROM (SELECT account_uuid FROM accounts WHERE owner_id = (SELECT user_uuid FROM users WHERE name = 'John Doe')) AS a
);

UPDATE accounts
SET points_balance = points_balance + 1
WHERE account_uuid = (
    SELECT a.account_uuid FROM (SELECT account_uuid FROM accounts WHERE owner_id = (SELECT user_uuid FROM users WHERE name = 'Jane Smith')) AS a
);

-- Inserting sample data into Transactions
INSERT INTO transactions (transaction_uuid, account_uuid, user_uuid, amount, date, store_uuid, points_earned) VALUES
(UUID(), (SELECT account_uuid FROM accounts WHERE owner_id = (SELECT user_uuid FROM users WHERE name = 'John Doe')), (SELECT user_uuid FROM users WHERE name = 'John Doe'), 4.00, NOW(), (SELECT store_uuid FROM stores WHERE name = 'Store B'), 1);


UPDATE accounts
SET points_balance = points_balance -1
WHERE account_uuid = (
    SELECT a.account_uuid FROM (SELECT account_uuid FROM accounts WHERE owner_id = (SELECT user_uuid FROM users WHERE name = 'John Doe')) AS a
);
