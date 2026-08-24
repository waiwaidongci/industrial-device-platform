CREATE TABLE devices (id text PRIMARY KEY, tenant_id text NOT NULL, name text NOT NULL, kind text NOT NULL, device_group text, tags jsonb NOT NULL DEFAULT '{}', status text NOT NULL, firmware_version text, parameters jsonb NOT NULL DEFAULT '{}', last_seen timestamptz, created_at timestamptz NOT NULL, updated_at timestamptz NOT NULL);
CREATE INDEX devices_tenant_group_idx ON devices(tenant_id, device_group);
CREATE TABLE device_states (device_id text REFERENCES devices(id), observed_at timestamptz NOT NULL, status text NOT NULL, payload jsonb NOT NULL DEFAULT '{}', PRIMARY KEY(device_id, observed_at));
CREATE TABLE commands (id text PRIMARY KEY, device_id text REFERENCES devices(id), type text NOT NULL, payload jsonb NOT NULL, idempotency_key text UNIQUE, status text NOT NULL, result text, created_at timestamptz NOT NULL, updated_at timestamptz NOT NULL);
CREATE TABLE firmware_releases (id text PRIMARY KEY, version text NOT NULL, canary_percent int NOT NULL, status text NOT NULL, created_at timestamptz NOT NULL);
CREATE TABLE alert_rules (id text PRIMARY KEY, name text NOT NULL, metric text NOT NULL, threshold numeric NOT NULL, severity text NOT NULL, enabled boolean NOT NULL DEFAULT true, created_at timestamptz NOT NULL);
CREATE TABLE events (id bigserial PRIMARY KEY, event_type text NOT NULL, device_id text, payload jsonb NOT NULL, created_at timestamptz NOT NULL);
