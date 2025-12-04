-- Users table
DROP TABLE IF EXISTS blogs CASCADE;
DROP TABLE IF EXISTS users CASCADE;

CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(20) UNIQUE NOT NULL,
    password VARCHAR(60) NOT NULL
);

-- Blogs table
CREATE TABLE blogs (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id),
    title VARCHAR(200) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    view_count INTEGER DEFAULT 0,
    search_vector tsvector
);

-- Indexes
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_blogs_user_id ON blogs(user_id);
CREATE INDEX idx_blogs_title ON blogs(title);
CREATE INDEX idx_blogs_created_at_id ON blogs(created_at DESC, id DESC);
CREATE INDEX idx_blogs_view_count ON blogs(view_count DESC);
CREATE INDEX idx_blogs_search ON blogs USING gin(search_vector);

-- updated_at trigger
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER update_blogs_updated_at
    BEFORE UPDATE ON blogs
    FOR EACH ROW
    EXECUTE FUNCTION update_updated_at_column();

-- full-text search trigger
CREATE OR REPLACE FUNCTION blogs_search_trigger()
RETURNS TRIGGER AS $$
BEGIN
  NEW.search_vector :=
    setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(NEW.content, '')), 'B');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER blogs_search_update
  BEFORE INSERT OR UPDATE ON blogs
  FOR EACH ROW
  EXECUTE FUNCTION blogs_search_trigger();

-- initial population
UPDATE blogs SET search_vector =
  setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
  setweight(to_tsvector('english', coalesce(content, '')), 'B');

-- utility: increment view count
CREATE OR REPLACE FUNCTION increment_blog_view(blog_id BIGINT)
RETURNS void AS $$
BEGIN
    UPDATE blogs SET view_count = view_count + 1 WHERE id = blog_id;
END;
$$ LANGUAGE plpgsql;
