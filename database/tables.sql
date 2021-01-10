-- extensions
CREATE extension IF NOT EXISTS "uuid-ossp";

-- drop tables if they already exist
drop table if exists 
    public.compute_result,
    public.compute,
    public.event
	CASCADE;

-- event
CREATE TABLE IF NOT EXISTS public.event (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    name VARCHAR(240) UNIQUE NOT NULL,
    depth REAL NOT NULL,
    creator BIGINT NOT NULL DEFAULT 0,
    create_date TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- compute
CREATE TABLE IF NOT EXISTS public.compute (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    event_id UUID REFERENCES event (id),
    fips VARCHAR(240),
    creator BIGINT NOT NULL DEFAULT 0,
    create_date TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- compute_result
-- JSON field should probably be normalized; For MVP, use JSON field
CREATE TABLE IF NOT EXISTS public.compute_result (
    id UUID PRIMARY KEY NOT NULL DEFAULT uuid_generate_v4(),
    compute_id UUID REFERENCES compute (id),
    result JSON NOT NULL
);

-- 
-- Views
-- 
CREATE OR REPLACE VIEW v_compute AS (
    SELECT c.id AS id,
           c.fips AS fips,
           c.creator AS creator,
           c.create_date AS create_date,
           e.id AS event_id,
           e.depth AS event_depth
    FROM compute c
    INNER JOIN event e ON e.id = c.event_id
);

-- Sample Events
INSERT INTO EVENT (id, name, depth) VALUES
    ('0d107163-0467-46c7-b579-65f1cfad7359', 'Test Event 1FT Depth', 1.0),
    ('9e4e7e3e-a648-4552-aa75-ccae22530482', 'Test Event 1.5FT Depth', 1.5);
