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

-- Sample Events
INSERT INTO EVENT (id, name, depth) VALUES
    ('0d107163-0467-46c7-b579-65f1cfad7359', 'Test Event 1FT Depth', 1.0),
    ('9e4e7e3e-a648-4552-aa75-ccae22530482', 'Test Event 1.5FT Depth', 1.5);
