CREATE TABLE IF NOT EXISTS products (
    id bigint(50) not null PRIMARY KEY,
    name text not null,
    image_url text not null,
    created_at datetime not null    
)
