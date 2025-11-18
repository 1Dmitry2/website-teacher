ALTER TABLE users 
ADD COLUMN email_verified BOOLEAN NOT NULL DEFAULT false,
ADD COLUMN email_verification_token VARCHAR,
ADD COLUMN email_verification_token_expires TIMESTAMP;

