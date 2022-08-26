CREATE TABLE IF NOT EXISTS users(
    id SERIAL PRIMARY KEY ,
    name TEXT NOT NULL ,
    email TEXT UNIQUE CHECK (email <>'') NOT NULL,
    phone_no INTEGER NOT NULL ,
    age      INTEGER NOT NULL ,
    gender   TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS restaurants(
    id SERIAL PRIMARY KEY ,
    name TEXT NOT NULL ,
    created_by INTEGER REFERENCES users(id) ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS dishes(
    id SERIAL PRIMARY KEY ,
    name TEXT NOT NULL ,
    price INTEGER NOT NULL ,
    user_id INTEGER REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS dish_per_restaurant(
    id SERIAL PRIMARY KEY ,
    restaurant_id INTEGER REFERENCES restaurants(id),
    dish_id       INTEGER REFERENCES dishes(id)
);

CREATE TABLE IF NOT EXISTS user_address(
    id SERIAL PRIMARY KEY ,
    user_id INTEGER REFERENCES users(id),
    address TEXT NOT NULL ,
    location point,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    archived_at TIMESTAMP WITH TIME ZONE
);

CREATE TABLE IF NOT EXISTS restaurant_address(
    id SERIAL PRIMARY KEY ,
    restaurant_id INTEGER REFERENCES restaurants(id),
    address TEXT NOT NULL ,
    location point,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    archived_at TIMESTAMP WITH TIME ZONE
);

create type user_type as enum('admin', 'sub_admin', 'user');

CREATE TABLE IF NOT EXISTS roles(
    id SERIAL PRIMARY KEY ,
    roles user_type ,
    user_id INTEGER REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS sessions(
    id uuid primary key default gen_random_uuid() not null ,
    user_id INTEGER REFERENCES users(id) NOT NULL ,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT now() NOT NULL ,
    expires_at TIMESTAMP WITH TIME ZONE
)

