
CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username VARCHAR(50) ,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash VARCHAR(255) NOT NULL,
    jwt_token TEXT,
    is_frozen BOOLEAN DEFAULT 0,
    attempt INTEGER DEFAULT 0,
    is_OAuth BOOLEAN DEFAULT 0,
    oauth_access_token TEXT,
    oauth_refresh_token TEXT,
    oauth_profile_data JSON,
    is_verified BOOLEAN DEFAULT 0,
    is_deleted BOOLEAN DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

SELECT * FROM sqlite_master WHERE type='table' AND name='users';

-- Add new columns to your existing users table
ALTER TABLE users 
ADD COLUMN otp_code VARCHAR(10),
ADD COLUMN otp_expires_at TIMESTAMP,
ADD COLUMN verified_at TIMESTAMP NULL;