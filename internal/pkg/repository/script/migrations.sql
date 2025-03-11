//users
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       full_name VARCHAR(255) NOT NULL,
                       email VARCHAR(255) UNIQUE NOT NULL,
                       password VARCHAR(255) NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
                       deleted_at TIMESTAMP DEFAULT NULL
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_full_name ON users(full_name);
CREATE INDEX idx_users_deleted_at ON users(deleted_at) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at);

ALTER TABLE users OWNER TO postgres;


//movies
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