-- +migrate Up
CREATE SCHEMA core;
CREATE SCHEMA audit;
CREATE SCHEMA support;
CREATE SCHEMA test;

CREATE TABLE core.clients (
  id uuid PRIMARY KEY,
  email text NOT NULL,
  name text NOT NULL,
-- Team avatar (S3 bucket key / url)
  avatar_url text,
-- Statistics (e.g: usage, preferences, etc.)
  metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ON core.clients (email);

CREATE TYPE core.widget_status as ENUM (
  'active',
  'cancelled'
);


CREATE TYPE core.widget_type as ENUM (
  'poll',
  'nps',
  'quiz',
  'form'
);

CREATE TABLE core.widgets (
  id uuid PRIMARY KEY,
  client_id uuid REFERENCES core.clients (id) ON DELETE CASCADE NOT NULL,
  title text NOT NULL,
  status core.widget_status NOT NULL,
  type core.widget_type NOT NULL,
  settings jsonb DEFAULT '{}'::jsonb NOT NULL,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP,
  updated_at timestamptz DEFAULT CURRENT_TIMESTAMP
);


CREATE TABLE core.questions (
  id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  widget_id uuid REFERENCES core.widgets (id) ON DELETE CASCADE NOT NULL,
  content text NOT NULL,
  sort_order integer NOT NULL DEFAULT 0,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ON core.questions (widget_id);

CREATE TABLE core.options (
  id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  question_id bigint REFERENCES core.questions (id) ON DELETE CASCADE NOT NULL,
  content text NOT NULL,
  -- Boolean used for quizzes
  is_correct boolean DEFAULT false ,
  sort_order integer NOT NULL DEFAULT 0
);

CREATE TABLE core.submissions (
  id uuid PRIMARY KEY,
  widget_id uuid REFERENCES core.widgets (id) ON DELETE CASCADE NOT NULL,
  -- Metadata for IP, user_agent, etc.
  metadata jsonb DEFAULT '{}'::jsonb NOT NULL,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE core.answers (
  id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
  widget_id uuid REFERENCES core.widgets (id) ON DELETE CASCADE NOT NULL,
  question_id bigint REFERENCES core.questions (id) ON DELETE CASCADE NOT NULL,
  option_id bigint REFERENCES core.options (id) ON DELETE CASCADE,
  submission_id uuid REFERENCES core.submissions (id) ON DELETE CASCADE NOT NULL,
  text_value text,
  numeric_value integer,
  created_at timestamptz DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX ON core.answers(submission_id);
CREATE INDEX ON core.answers(question_id);
CREATE INDEX ON core.answers (widget_id); 


-- +migrate Down
DROP TABLE IF EXISTS core.answers CASCADE;
DROP TABLE IF EXISTS core.submissions CASCADE;
DROP TABLE IF EXISTS core.options CASCADE;
DROP TABLE IF EXISTS core.questions CASCADE;
DROP TABLE IF EXISTS core.widgets CASCADE;
DROP TABLE IF EXISTS core.clients CASCADE;

DROP TYPE IF EXISTS core.widget_type;
DROP TYPE IF EXISTS core.widget_status;

DROP SCHEMA IF EXISTS core CASCADE;
DROP SCHEMA IF EXISTS audit CASCADE;
DROP SCHEMA IF EXISTS support CASCADE;
DROP SCHEMA IF EXISTS test CASCADE;
