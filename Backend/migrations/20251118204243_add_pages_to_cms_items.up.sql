-- Добавляем pages в blocks
ALTER TABLE blocks
    ADD COLUMN IF NOT EXISTS pages JSONB NOT NULL DEFAULT '[]'::jsonb;

-- Добавляем pages и text в gallery_items
ALTER TABLE gallery_items
    ADD COLUMN IF NOT EXISTS pages JSONB NOT NULL DEFAULT '[]'::jsonb,
    ADD COLUMN IF NOT EXISTS text TEXT;

-- Добавляем pages в slider_items
ALTER TABLE slider_items
    ADD COLUMN IF NOT EXISTS pages JSONB NOT NULL DEFAULT '[]'::jsonb;

