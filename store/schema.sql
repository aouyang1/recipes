CREATE DATABASE IF NOT EXISTS recipes;

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
    recipe_id BIGINT(20) UNSIGNED,
    PRIMARY KEY (recipe_id, recipe_event_id)
);
CREATE INDEX recipe_id_idx ON recipe_event_to_recipe (recipe_id);

CREATE TABLE IF NOT EXISTS recipe (
    id BIGINT(20) UNSIGNED AUTO_INCREMENT,
    name VARCHAR(255),
    variant VARCHAR(255),
    created_on BIGINT(20),
    PRIMARY KEY (id)
);
CREATE INDEX name_idx ON recipe (name);
CREATE INDEX created_on_idx ON recipe (created_on);

CREATE TABLE IF NOT EXISTS recipe_to_tag (
    recipe_id BIGINT(20) UNSIGNED,
    tag_id BIGINT(20) UNSIGNED,
    PRIMARY KEY (recipe_id, tag_id)
);
CREATE INDEX tag_idx ON recipe_to_tag (tag);

CREATE TABLE IF NOT EXISTS tag (
    id BIGINT(20) UNSIGNED AUTO_INCREMENT,
    name VARCHAR(255),
    PRIMARY KEY (id)
);
CREATE INDEX name_idx ON tag (name);

CREATE TABLE IF NOT EXISTS ingredient (
    id BIGINT(20) UNSIGNED AUTO_INCREMENT,
    name VARCHAR(255),
    PRIMARY KEY (id)
);
CREATE INDEX name_idx ON ingredient (name);

CREATE TABLE IF NOT EXISTS recipe_to_ingredient (
    recipe_id BIGINT(20) UNSIGNED,
    ingredient_id BIGINT(20) UNSIGNED,
    quantity INT,
    unit VARCHAR(10),
    size VARCHAR(5),
    PRIMARY KEY (recipe_id, ingredient_id)
);
CREATE INDEX ingredient_id_idx ON recipe_to_ingredient (ingredient_id);

GRANT ALL ON recipes.* TO 'recipes'@'localhost';
