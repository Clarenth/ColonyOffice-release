CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

CREATE TABLE IF NOT EXISTS accounts_employee
(
  id_code uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  email varchar NOT NULL UNIQUE,
  password varchar NOT NULL,
  phone_number varchar NOT NULL UNIQUE,
  job_title varchar NOT NULL DEFAULT 'unassigned',
  office_address varchar NOT NULL DEFAULT '',
  security_access_level varchar NOT NULL DEFAULT 'official',
  language varchar NOT NULL,
  created_at timestamp with time zone,
  updated_at timestamp with time zone
);

CREATE TABLE IF NOT EXISTS accounts_identity
(
  id_code uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  first_name varchar,
  middle_name varchar,
  last_name varchar,
  sex varchar,
  gender varchar,
  age bigint,
  height varchar,
  home_address varchar NOT NULL DEFAULT '',
  birthdate varchar,
  birthplace varchar,
  updated_at timestamp with time zone
  --FOREIGN KEY (id_code) REFERENCES accounts_employee (id_code)
);

CREATE TABLE IF NOT EXISTS security_access_level
(
  official varchar,
  secret varchar,
  top_secret varchar
);

CREATE TABLE IF NOT EXISTS documents
(
	document_id uuid DEFAULT uuid_generate_v4() PRIMARY KEY,
	document_title varchar,
	author_name varchar,
	author_id uuid DEFAULT uuid_generate_v4(),
	description varchar,
	cdn_url varchar,
	security_access_level varchar,
	created_at timestamp,
	updated_at timestamp,
	language  varchar
);

CREATE TABLE IF NOT EXISTS files
(
	file_id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  document_id uuid NOT NULL DEFAULT uuid_generate_v4(),
  file_title varchar(255),
  title_hash varchar (255),
  author_id uuid DEFAULT uuid_generate_v4(),
  author_name varchar,
  security_access_level varchar,
  created_at timestamp,
	updated_at timestamp
);

CREATE TABLE IF NOT EXISTS telephone_number
(
  id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
  account_employee uuid NOT NULL DEFAULT uuid_generate_v4(),
  phone_number varchar(20)
);

CREATE TABLE IF NOT EXISTS signin_history
(
  id int,
  email varchar PRIMARY KEY,
  password varchar,
  status boolean,
  signin_count int,
  ip_address inet,
  user_agent varchar NOT NULL,
  login_timestamp timestamp,
  locked boolean
);

CREATE TABLE IF NOT EXISTS employee_identity (
  employee_id_code uuid DEFAULT uuid_generate_v4(),
  identity_id_code uuid DEFAULT uuid_generate_v4()
);
