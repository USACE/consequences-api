-- Users and Roles for Consequences Database

-- User consequences_user
-- Note: Substitute real password for 'password'
CREATE USER consequences_user WITH ENCRYPTED PASSWORD 'password';
CREATE ROLE consequences_reader;
CREATE ROLE consequences_writer;
CREATE ROLE postgis_reader;

--------------------------------------------------------------------------
-- NOTE: IF USERS ALREADY EXIST ON DATABASE, JUST RUN FROM THIS POINT DOWN
--------------------------------------------------------------------------

-- Role consequences_reader
-- Tables specific to instrumentation app
GRANT SELECT ON
    event,
    compute,
    compute_result,
    v_compute
TO consequences_reader;

-- Role consequences_writer
-- Tables specific to instrumentation app
GRANT INSERT,UPDATE,DELETE ON
    event,
    compute,
    compute_result
TO consequences_writer;

-- Role postgis_reader
GRANT SELECT ON geometry_columns TO postgis_reader;
GRANT SELECT ON geography_columns TO postgis_reader;
GRANT SELECT ON spatial_ref_sys TO postgis_reader;
-- Grant Permissions to instrument_user
GRANT postgis_reader TO consequences_user;
GRANT consequences_reader TO consequences_user;
GRANT consequences_writer TO consequences_user;
