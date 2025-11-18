ALTER TABLE users 
DROP COLUMN IF EXISTS email_verified,
DROP COLUMN IF EXISTS email_verification_token,
DROP COLUMN IF EXISTS email_verification_token_expires;

