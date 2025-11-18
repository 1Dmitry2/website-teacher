ALTER TABLE posts
    DROP COLUMN IF EXISTS videos,
    DROP COLUMN IF EXISTS alignment,
    DROP COLUMN IF EXISTS title_position,
    DROP COLUMN IF EXISTS content_position;

