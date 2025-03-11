CREATE TABLE movies (
                        id SERIAL PRIMARY KEY,
                        title VARCHAR(255) NOT NULL,
                        director VARCHAR(255) NOT NULL,
                        year INTEGER NOT NULL,
                        plot TEXT,
                        rating FLOAT,
                        created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
                        deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

CREATE INDEX idx_movies_title ON movies(title);
CREATE INDEX idx_movies_director ON movies(director);
CREATE INDEX idx_movies_year ON movies(year);
CREATE INDEX idx_movies_rating ON movies(rating);
CREATE INDEX idx_movies_deleted_at ON movies(deleted_at) WHERE deleted_at IS NULL;

ALTER TABLE movies OWNER TO postgres;