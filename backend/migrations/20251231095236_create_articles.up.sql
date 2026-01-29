CREATE TABLE IF NOT EXISTS articles (
    id BIGSERIAL PRIMARY KEY,
    title TEXT NOT NULL,
    slug TEXT NOT NULL UNIQUE,
    excerpt TEXT,
    author TEXT,
    published_date TIMESTAMPTZ,
    featured_image TEXT,
    is_public BOOLEAN NOT NULL DEFAULT false,
    categories JSONB DEFAULT '[]'::jsonb,
    content TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS article_images (
    id BIGSERIAL PRIMARY KEY,
    article_id BIGINT NOT NULL REFERENCES articles(id) ON DELETE CASCADE,
    file_path TEXT NOT NULL,
    alignment TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_articles_slug ON articles (slug);
CREATE INDEX IF NOT EXISTS idx_articles_is_public ON articles (is_public);
CREATE INDEX IF NOT EXISTS idx_articles_published_date ON articles (published_date);
CREATE INDEX IF NOT EXISTS idx_article_images_article_id ON article_images (article_id);
CREATE INDEX IF NOT EXISTS idx_article_images_file_path ON article_images (file_path);

