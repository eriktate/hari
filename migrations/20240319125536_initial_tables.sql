-- migrate:up
CREATE OR REPLACE FUNCTION manage_table_updated_at()
RETURNS TRIGGER AS $$
BEGIN
	NEW.updated_at = now();
	RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TABLE IF NOT EXISTS accounts(
	id UUID PRIMARY KEY,
	name VARCHAR(512),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON accounts FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();


CREATE TABLE IF NOT EXISTS users(
	id UUID PRIMARY KEY,
	email VARCHAR(320) NOT NULL,
	name VARCHAR(512) NOT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON users FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();


CREATE TABLE IF NOT EXISTS account_users(
	account_id UUID NOT NULL REFERENCES accounts(id),
	user_id UUID NOT NULL REFERENCES users(id),
	PRIMARY KEY (account_id, user_id)
);


CREATE TABLE IF NOT EXISTS webhooks(
	id UUID PRIMARY KEY,
	account_id UUID NOT NULL REFERENCES accounts(id),
	name VARCHAR(512) NOT NULL,
	"key" VARCHAR(512) NOT NULL,
	static_data BYTEA DEFAULT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	UNIQUE(account_id, key)
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON webhooks FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();


CREATE TABLE IF NOT EXISTS target_status(
	status VARCHAR(64) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS targets(
	id UUID PRIMARY KEY,
	webhook_id UUID NOT NULL REFERENCES webhooks(id),
	url TEXT NOT NULL,
	status VARCHAR(64) NOT NULL REFERENCES target_status(status),
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON targets FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();


CREATE TABLE IF NOT EXISTS hook_status(
	status VARCHAR(64) PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS hooks(
	id UUID PRIMARY KEY,
	webhook_id UUID NOT NULL REFERENCES webhooks(id),
	status VARCHAR(64) NOT NULL REFERENCES hook_status(status),
	payload BYTEA DEFAULT NULL,
	run_at TIMESTAMP DEFAULT NULL,
	created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER manage_updated_at BEFORE UPDATE
ON hooks FOR EACH ROW EXECUTE PROCEDURE manage_table_updated_at();

INSERT INTO target_status
	(status)
VALUES
	('ok'),
	('not_responding'),
	('error')
ON CONFLICT DO NOTHING;

INSERT INTO hook_status
	(status)
VALUES
	('success'),
	('pending'),
	('scheduled'),
	('failure')
ON CONFLICT DO NOTHING;

INSERT INTO accounts
	(id, name)
VALUES
	('688eeae7-eae7-4712-b68d-8fbadd6bd4d5', 'Hari Test Account')
ON CONFLICT DO NOTHING;

-- create app user
DO
$do$
BEGIN
	IF NOT EXISTS (
		SELECT FROM pg_catalog.pg_user
		WHERE usename = 'hari') THEN

		CREATE USER hari WITH ENCRYPTED PASSWORD 'hari';
		REVOKE CONNECT ON DATABASE hari FROM PUBLIC;

		GRANT CONNECT
		ON DATABASE hari
		TO hari;

		GRANT SELECT, INSERT, UPDATE, DELETE
		ON ALL TABLES IN SCHEMA public
		TO hari;
	END IF;
END
$do$;

-- migrate:down
DROP TRIGGER IF EXISTS manage_updated_at ON hooks;
DROP TRIGGER IF EXISTS manage_updated_at ON targets;
DROP TRIGGER IF EXISTS manage_updated_at ON users;
DROP TRIGGER IF EXISTS manage_updated_at ON accounts;
DROP TRIGGER IF EXISTS manage_updated_at ON webhooks;

DROP TABLE IF EXISTS hooks;
DROP TABLE IF EXISTS hook_status;
DROP TABLE IF EXISTS targets;
DROP TABLE IF EXISTS target_status;
DROP TABLE IF EXISTS webhooks;
DROP TABLE IF EXISTS account_users;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS accounts;
