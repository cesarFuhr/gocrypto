CREATE TABLE IF NOT EXISTS keys(
  id uuid PRIMARY KEY,
  scope VARCHAR(50),
  expiration TIMESTAMP,
  creation TIMESTAMP,
  priv bytea,
  pub bytea
)