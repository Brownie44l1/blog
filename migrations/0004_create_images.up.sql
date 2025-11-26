CREATE TABLE images ( 
    id UUID  PRIMARY KEY DEFAULT uuid_generate_v4(),
    blog_id  UUID NOT NULL REFERENCES blogs(id) ON DELETE CASCADE,
    url      TEXT NOT NULL
);