CREATE DATABASE IF NOT EXISTS recipes;
GRANT ALL ON recipes.* TO 'recipes'@'localhost';

CREATE TABLE IF NOT EXISTS recipe_event (
    id VARCHAR(60),
    title VARCHAR(255) NOT NULL,
    schedule_date INT,
    description TEXT,
    PRIMARY KEY (id)
);
CREATE INDEX schedule_date_idx ON recipe_event (schedule_date);

CREATE TABLE IF NOT EXISTS recipe_event_to_recipe (
    recipe_event_id VARCHAR(60),
    recipe_id VARCHAR(36),
    PRIMARY KEY (recipe_id, recipe_event_id)
);
CREATE INDEX recipe_id_idx ON recipe_event_to_recipe (recipe_id);

CREATE TABLE IF NOT EXISTS recipe (
    id VARCHAR(36),
    name VARCHAR(255),
    variant VARCHAR(255),
    created_on INT,
    PRIMARY KEY (id)
);
CREATE INDEX name_idx ON recipe (name);
CREATE INDEX created_on_idx ON recipe (created_on);

CREATE TABLE IF NOT EXISTS recipe_tag (
    recipe_id VARCHAR(36),
    tag VARCHAR(255),
    PRIMARY KEY (tag, recipe_id)
);
CREATE INDEX tag_idx ON recipe_tag (tag);

CREATE TABLE IF NOT EXISTS ingredient (
    id VARCHAR(36),
    name VARCHAR(255),
    PRIMARY KEY (id)
);
CREATE INDEX name_idx ON ingredient (name);

CREATE TABLE IF NOT EXISTS recipe_ingredient (
    recipe_id VARCHAR(36),
    ingredient_id VARCHAR(36),
    quantity INT,
    unit VARCHAR(10),
    size VARCHAR(5),
    PRIMARY KEY (recipe_id, ingredient_id)
);
CREATE INDEX ingredient_id_idx ON recipe_ingredient (ingredient_id);



