-- Create the database
CREATE DATABASE IF NOT EXISTS loyalty_program;
USE loyalty_program;

-- Create the accounts table
CREATE TABLE IF NOT EXISTS accounts (
    account_uuid CHAR(36) PRIMARY KEY,
    owner_id CHAR(36),
    creation_date DATETIME,
    region VARCHAR(255)
) ENGINE=NDBCLUSTER;

-- Create the users table
CREATE TABLE IF NOT EXISTS users (
    user_uuid CHAR(36) PRIMARY KEY,
    account_uuid CHAR(36),
    name VARCHAR(255),
    email_address VARCHAR(255) UNIQUE,
    phone_number VARCHAR(20),
    creation_date DATETIME,
    invite_code CHAR(36) UNIQUE NULL
) ENGINE=NDBCLUSTER;

-- Create the transactions table
CREATE TABLE IF NOT EXISTS transactions (
    transaction_uuid CHAR(36) PRIMARY KEY,
    account_uuid CHAR(36),
    user_uuid CHAR(36),
    amount DECIMAL(10,2),
    date DATETIME,
    store_uuid CHAR(36),
    points_earned INT
) ENGINE=NDBCLUSTER;

-- Create the transaction_items table
CREATE TABLE IF NOT EXISTS transaction_items (
    transaction_item_uuid CHAR(36) PRIMARY KEY,
    transaction_uuid CHAR(36),
    item_number INT,
    item VARCHAR(255),
    amount DECIMAL(10,2)
) ENGINE=NDBCLUSTER;

-- Create the points_redemption table
CREATE TABLE IF NOT EXISTS points_redemption (
    redemption_uuid CHAR(36) PRIMARY KEY,
    user_uuid CHAR(36),
    redemption_date DATETIME,
    points_used INT,
    reward_description VARCHAR(255)
) ENGINE=NDBCLUSTER;

-- Create the stores table
CREATE TABLE IF NOT EXISTS stores (
    store_uuid CHAR(36) PRIMARY KEY,
    name VARCHAR(255),
    region VARCHAR(255)
) ENGINE=NDBCLUSTER;

-- Create the invitations table
CREATE TABLE IF NOT EXISTS invitations (
    invitation_uuid CHAR(36) PRIMARY KEY,
    email VARCHAR(255) NOT NULL,
    account_uuid CHAR(36),
    inviter_uuid CHAR(36),
    token CHAR(36) UNIQUE NOT NULL,
    creation_date DATETIME NOT NULL,
    expiration_date DATETIME NOT NULL,
    status VARCHAR(20) NOT NULL,
    CONSTRAINT fk_invitations_accounts FOREIGN KEY (account_uuid) REFERENCES accounts(account_uuid),
    CONSTRAINT fk_invitations_inviter FOREIGN KEY (inviter_uuid) REFERENCES users(user_uuid)
) ENGINE=NDBCLUSTER;


-- Add accounts references
ALTER TABLE accounts
ADD FOREIGN KEY (owner_id) REFERENCES users(user_uuid);

-- Add users references
ALTER TABLE users
ADD FOREIGN KEY (account_uuid) REFERENCES accounts(account_uuid);

-- Add transactions references
ALTER TABLE transactions
ADD FOREIGN KEY (account_uuid) REFERENCES accounts(account_uuid),
ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid),
ADD FOREIGN KEY (store_uuid) REFERENCES stores(store_uuid);

-- Add transaction_items references
ALTER TABLE transaction_items
ADD FOREIGN KEY (transaction_uuid) REFERENCES transactions(transaction_uuid);

-- Add points_redemption references
ALTER TABLE points_redemption
ADD FOREIGN KEY (user_uuid) REFERENCES users(user_uuid);

-- Add points_balance to accounts 
ALTER TABLE accounts
ADD COLUMN points_balance INT DEFAULT 0;

-- Add password to user
ALTER TABLE users
ADD COLUMN password VARCHAR(255);