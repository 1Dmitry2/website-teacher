ALTER TABLE posts
    ADD COLUMN IF NOT EXISTS videos JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS alignment TEXT DEFAULT 'full-width',
    ADD COLUMN IF NOT EXISTS title_position TEXT DEFAULT 'top',
    ADD COLUMN IF NOT EXISTS content_position TEXT DEFAULT 'bottom';

