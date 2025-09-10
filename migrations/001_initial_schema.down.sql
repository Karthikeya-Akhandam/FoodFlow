-- Drop triggers
DROP TRIGGER IF EXISTS update_donation_offers_updated_at ON donation_offers;
DROP TRIGGER IF EXISTS update_collaborators_updated_at ON collaborators;
DROP TRIGGER IF EXISTS update_organizations_updated_at ON organizations;
DROP TRIGGER IF EXISTS update_profiles_updated_at ON profiles;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables (in reverse dependency order)
DROP TABLE IF EXISTS events;
DROP TABLE IF EXISTS remote_org_assignments;
DROP TABLE IF EXISTS redemptions;
DROP TABLE IF EXISTS tokens;
DROP TABLE IF EXISTS credits;
DROP TABLE IF EXISTS donation_claims;
DROP TABLE IF EXISTS donation_offers;
DROP TABLE IF EXISTS collaborators;
DROP TABLE IF EXISTS organizations;
DROP TABLE IF EXISTS profiles;
DROP TABLE IF EXISTS users;

-- Drop extensions
DROP EXTENSION IF EXISTS "uuid-ossp";