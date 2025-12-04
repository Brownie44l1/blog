-- Add tsvector column for full-text search
ALTER TABLE blogs ADD COLUMN search_vector tsvector;

-- Create index
CREATE INDEX idx_blogs_search ON blogs USING gin(search_vector);

-- Auto-update search vector
CREATE OR REPLACE FUNCTION blogs_search_trigger() RETURNS trigger AS $$
BEGIN
  NEW.search_vector := 
    setweight(to_tsvector('english', coalesce(NEW.title, '')), 'A') ||
    setweight(to_tsvector('english', coalesce(NEW.content, '')), 'B');
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER blogs_search_update 
  BEFORE INSERT OR UPDATE ON blogs
  FOR EACH ROW EXECUTE FUNCTION blogs_search_trigger();

-- Populate existing data
UPDATE blogs SET search_vector = 
  setweight(to_tsvector('english', coalesce(title, '')), 'A') ||
  setweight(to_tsvector('english', coalesce(content, '')), 'B');