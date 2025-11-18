-- Удаляем pages из blocks
ALTER TABLE blocks
    DROP COLUMN IF EXISTS pages;

-- Удаляем pages и text из gallery_items
ALTER TABLE gallery_items
    DROP COLUMN IF EXISTS pages,
    DROP COLUMN IF EXISTS text;

-- Удаляем pages из slider_items
ALTER TABLE slider_items
    DROP COLUMN IF EXISTS pages;

