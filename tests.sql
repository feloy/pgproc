BEGIN TRANSACTION;

DROP SCHEMA IF EXISTS tests CASCADE;
CREATE SCHEMA tests;

CREATE TABLE tests.content (
  cnt_id serial PRIMARY KEY,
  cnt_name text NOT NULL
);

CREATE TYPE tests.enumtype AS ENUM ('val1', 'val2', 'val3');

CREATE FUNCTION tests.test_returns_integer()
RETURNS integer
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT 42;
$$;

CREATE FUNCTION tests.test_returns_setof_integer()
RETURNS SETOF integer
LANGUAGE PLPGSQL
IMMUTABLE
AS $$
BEGIN
  RETURN NEXT 42;
  RETURN NEXT 43;
  RETURN NEXT 44;
END;
$$;

CREATE FUNCTION tests.test_returns_integer_as_string()
RETURNS character varying
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT '42'::varchar;
$$;

CREATE FUNCTION tests.test_returns_string()
RETURNS character varying
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT 'hello'::varchar;
$$;

CREATE FUNCTION tests.test_returns_setof_string()
RETURNS SETOF character varying
LANGUAGE PLPGSQL
IMMUTABLE
AS $$
BEGIN
  RETURN NEXT 'hello'::varchar;
  RETURN NEXT 'world'::varchar;
  RETURN NEXT '!'::varchar;
END;
$$;

CREATE FUNCTION tests.test_returns_numeric()
RETURNS numeric
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT 3.14159::numeric;
$$;

CREATE FUNCTION tests.test_returns_setof_numeric()
RETURNS SETOF numeric
LANGUAGE PLPGSQL
IMMUTABLE
AS $$
BEGIN
  RETURN NEXT 3.14159::numeric;
  RETURN NEXT 4.49::numeric;
END;
$$;

CREATE FUNCTION tests.test_returns_real()
RETURNS real
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT 3.14::real;
$$;

CREATE FUNCTION tests.test_returns_setof_real()
RETURNS SETOF real
LANGUAGE PLPGSQL
IMMUTABLE
AS $$
BEGIN
  RETURN NEXT 3.14::real;
  RETURN NEXT 4.49::real;
END;
$$;

CREATE FUNCTION tests.test_returns_bool_true()
RETURNS boolean
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT true;
$$;

CREATE FUNCTION tests.test_returns_bool_false()
RETURNS boolean
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT false;
$$;

CREATE FUNCTION tests.test_returns_setof_bool()
RETURNS SETOF boolean
LANGUAGE PLPGSQL
IMMUTABLE
AS $$
BEGIN
  RETURN NEXT false;
  RETURN NEXT true;
  RETURN NEXT true;
  RETURN NEXT false;
END;
$$;

CREATE FUNCTION tests.test_returns_date()
RETURNS date
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT CURRENT_TIMESTAMP::date;
$$;

CREATE FUNCTION tests.test_returns_infinity_date()
RETURNS date
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT 'infinity'::date;
$$;

CREATE FUNCTION tests.test_returns_minus_infinity_date()
RETURNS date
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT '-infinity'::date;
$$;

CREATE FUNCTION tests.test_returns_64bits_date()
RETURNS date
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT '2040-01-01'::date;
$$;

CREATE FUNCTION tests.test_returns_setof_date()
RETURNS SETOF date
LANGUAGE PLPGSQL
IMMUTABLE
AS $$
BEGIN
  RETURN NEXT '2015-01-01';
  RETURN NEXT '2016-02-02';
END;
$$;

CREATE FUNCTION tests.test_returns_timestamp()
RETURNS timestamp
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT CURRENT_TIMESTAMP::timestamp;
$$;

CREATE FUNCTION tests.test_returns_time()
RETURNS time
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT CURRENT_TIMESTAMP::time;
$$;

CREATE TYPE tests.composite1 AS (
  a integer,
  b varchar
);

CREATE FUNCTION tests.test_returns_composite()
RETURNS tests.composite1
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT (1, 'hello')::tests.composite1;
$$;

CREATE FUNCTION tests.test_returns_setof_composite()
RETURNS SETOF tests.composite1
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT (1, 'hello')::tests.composite1
  UNION SELECT (2, 'bye')::tests.composite1;
$$;

CREATE FUNCTION tests.test_returns_enum() 
RETURNS tests.enumtype 
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT 'val1'::tests.enumtype;
$$;

CREATE FUNCTION tests.test_returns_enum_array() 
RETURNS tests.enumtype[] 
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT ARRAY['val1', 'val2']::tests.enumtype[];
$$;

CREATE FUNCTION tests.test_returns_null_enum_array() 
RETURNS tests.enumtype[] 
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT NULL::tests.enumtype[];
$$;

CREATE FUNCTION tests._hidden_function()
RETURNS boolean
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT true;
$$;

CREATE FUNCTION tests.function_in_tests_schema()
RETURNS boolean
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT true;
$$;

CREATE FUNCTION tests.function_raising_exception()
RETURNS boolean
LANGUAGE PLPGSQL
IMMUTABLE
AS $$
BEGIN
  RAISE EXCEPTION '"a particular exception message"';
  SELECT true;
END;
$$;

-- test arguments
CREATE FUNCTION tests.test_returns_incremented_integer(n integer)
RETURNS integer
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT $1 + 1;
$$;

CREATE FUNCTION tests.test_returns_incremented_numeric(n numeric)
RETURNS numeric
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT $1 + 1.5;
$$;

CREATE FUNCTION tests.test_returns_incremented_real(n real)
RETURNS real
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT ($1 + 1.42)::real;
$$;

CREATE FUNCTION tests.test_returns_cat_string(s varchar)
RETURNS varchar
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT $1 || '.';
$$;

CREATE FUNCTION tests.test_returns_same_bool(b boolean)
RETURNS boolean
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT $1;
$$;

CREATE FUNCTION tests.test_returns_same_date(d date)
RETURNS date
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT $1;
$$;

CREATE FUNCTION tests.test_returns_same_timestamp(t timestamp)
RETURNS timestamp
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT $1;
$$;

CREATE FUNCTION tests.test_returns_same_time(t time)
RETURNS time
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT $1;
$$;

CREATE FUNCTION tests.test_integer_array_arg(list integer[]) 
RETURNS SETOF integer
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE 
  i integer;
BEGIN
  FOREACH i IN ARRAY list LOOP
    RETURN NEXT i;
  END LOOP;
END;
$$;

CREATE FUNCTION tests.test_varchar_array_arg(list varchar[]) 
RETURNS SETOF varchar
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE 
  i varchar;
BEGIN
  FOREACH i IN ARRAY list LOOP
    RETURN NEXT i;
  END LOOP;
END;
$$;

CREATE FUNCTION tests.test_enum_arg(enumval tests.enumtype) 
RETURNS tests.enumtype
LANGUAGE plpgsql
IMMUTABLE
AS $$
BEGIN
  RETURN enumval;
END;
$$;

CREATE FUNCTION tests.test_enum_array_arg(list tests.enumtype[]) 
RETURNS SETOF tests.enumtype
LANGUAGE plpgsql
IMMUTABLE
AS $$
DECLARE 
  i varchar;
BEGIN
  FOREACH i IN ARRAY list LOOP
    RETURN NEXT i;
  END LOOP;
END;
$$;

CREATE FUNCTION tests.test_returns_accented_string()
RETURNS character varying
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT 'héllo'::varchar;
$$;

CREATE FUNCTION tests.test_returns_empty_array()
RETURNS integer[]
LANGUAGE SQL
IMMUTABLE
AS $$
  SELECT '{}'::integer[];
$$;

CREATE FUNCTION tests.content_add(prm_name text) 
RETURNS integer
LANGUAGE plpgsql
VOLATILE
AS $$
DECLARE
  ret integer;
BEGIN
  INSERT INTO tests.content (cnt_name) VALUES (prm_name)
    RETURNING cnt_id INTO ret;
  RETURN ret;
END;
$$;

CREATE FUNCTION tests.content_get(prm_id integer)
RETURNS tests.content
LANGUAGE plpgsql
STABLE
AS $$
DECLARE
  ret tests.content;
BEGIN
  SELECT * INTO ret FROM tests.content WHERE cnt_id = prm_id;
  RETURN ret;
END;
$$;

DROP FUNCTION IF EXISTS public.tests_get_one();
CREATE FUNCTION public.tests_get_one()
RETURNS integer
LANGUAGE plpgsql
IMMUTABLE
as $$
BEGIN
  RETURN 1;
END;
$$;

CREATE FUNCTION tests.returns_json() 
RETURNS json
LANGUAGE plpgsql
STABLE 
AS $$
DECLARE
  ret json;
BEGIN
  SELECT json_agg(row_to_json(d)) INTO ret FROM (
    SELECT 1 AS id, 'One' AS name
    UNION
    SELECT 2, 'Two'
  ) d;
  RETURN ret;
END;
$$;

COMMIT;
