-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) NOT NULL CHECK (role IN ('ADMIN', 'ORG', 'COLLAB')),
    status VARCHAR(20) NOT NULL DEFAULT 'ACTIVE' CHECK (status IN ('ACTIVE', 'SUSPENDED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Profiles table
CREATE TABLE profiles (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    phone VARCHAR(15),
    pincode VARCHAR(6),
    city VARCHAR(50),
    district VARCHAR(50),
    state VARCHAR(50),
    address VARCHAR(200),
    is_remote_org BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Organizations table
CREATE TABLE organizations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    org_type VARCHAR(50) NOT NULL,
    purpose_focus TEXT[] NOT NULL DEFAULT '{}',
    last_month_people_fed INTEGER DEFAULT 0 CHECK (last_month_people_fed >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Collaborators table
CREATE TABLE collaborators (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    collab_type VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Donation offers table
CREATE TABLE donation_offers (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    collab_id UUID NOT NULL REFERENCES collaborators(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description VARCHAR(500),
    ready_from TIMESTAMP WITH TIME ZONE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    estimated_servings INTEGER NOT NULL CHECK (estimated_servings > 0),
    purpose VARCHAR(20) CHECK (purpose IN ('CHILDREN', 'ELDERLY', 'WOMEN', 'GENERAL', 'EMERGENCY')),
    pincode VARCHAR(6),
    city VARCHAR(50),
    state VARCHAR(50),
    status VARCHAR(20) NOT NULL DEFAULT 'OPEN' CHECK (status IN ('OPEN', 'PENDING_CONFIRM', 'CLAIMED', 'EXPIRED', 'CANCELLED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_timing CHECK (expires_at > ready_from),
    CONSTRAINT valid_servings CHECK (estimated_servings <= 10000)
);

-- Donation claims table
CREATE TABLE donation_claims (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL REFERENCES donation_offers(id) ON DELETE CASCADE,
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    requested_servings INTEGER NOT NULL CHECK (requested_servings > 0),
    priority_score DECIMAL(5,2) DEFAULT 0.0 CHECK (priority_score >= 0 AND priority_score <= 100),
    status VARCHAR(20) NOT NULL DEFAULT 'REQUESTED' CHECK (status IN ('REQUESTED', 'WON', 'LOST', 'CANCELLED')),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    CONSTRAINT valid_servings CHECK (requested_servings <= 10000)
);

-- Credits table (monthly credits for organizations)
CREATE TABLE credits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    month VARCHAR(6) NOT NULL, -- Format: YYYYMM
    credits_issued INTEGER NOT NULL DEFAULT 0 CHECK (credits_issued >= 0),
    credits_remaining INTEGER NOT NULL DEFAULT 0 CHECK (credits_remaining >= 0),
    is_remote_bonus BOOLEAN DEFAULT FALSE,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(org_id, month),
    CONSTRAINT valid_credits CHECK (credits_remaining <= credits_issued)
);

-- Tokens table (monthly tokens for collaborators)
CREATE TABLE tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    collab_id UUID NOT NULL REFERENCES collaborators(id) ON DELETE CASCADE,
    month VARCHAR(6) NOT NULL, -- Format: YYYYMM
    tokens_earned INTEGER NOT NULL DEFAULT 0 CHECK (tokens_earned >= 0),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(collab_id, month)
);

-- Redemptions table (completed transactions)
CREATE TABLE redemptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    offer_id UUID NOT NULL REFERENCES donation_offers(id),
    org_id UUID NOT NULL REFERENCES organizations(id),
    collab_id UUID NOT NULL REFERENCES collaborators(id),
    servings_accepted INTEGER NOT NULL CHECK (servings_accepted > 0),
    credits_spent INTEGER NOT NULL CHECK (credits_spent > 0),
    tokens_awarded INTEGER NOT NULL CHECK (tokens_awarded >= 0),
    confirmed_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Remote organization assignments table
CREATE TABLE remote_org_assignments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    remote_org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    proxy_org_id UUID NOT NULL REFERENCES organizations(id) ON DELETE CASCADE,
    active BOOLEAN DEFAULT TRUE,
    notes VARCHAR(500),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(remote_org_id, proxy_org_id)
);

-- Audit events table
CREATE TABLE events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID,
    event_type VARCHAR(50) NOT NULL,
    payload JSONB,
    user_id UUID REFERENCES users(id),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);

-- Create indexes for performance
CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_role ON users(role);
CREATE INDEX idx_profiles_pincode ON profiles(pincode);
CREATE INDEX idx_profiles_city ON profiles(city);
CREATE INDEX idx_profiles_state ON profiles(state);
CREATE INDEX idx_organizations_user_id ON organizations(user_id);
CREATE INDEX idx_collaborators_user_id ON collaborators(user_id);
CREATE INDEX idx_donation_offers_collab_id ON donation_offers(collab_id);
CREATE INDEX idx_donation_offers_status ON donation_offers(status);
CREATE INDEX idx_donation_offers_expires_at ON donation_offers(expires_at);
CREATE INDEX idx_donation_offers_location ON donation_offers(pincode, city, state);
CREATE INDEX idx_donation_offers_purpose ON donation_offers(purpose);
CREATE INDEX idx_donation_claims_offer_id ON donation_claims(offer_id);
CREATE INDEX idx_donation_claims_org_id ON donation_claims(org_id);
CREATE INDEX idx_donation_claims_status ON donation_claims(status);
CREATE INDEX idx_credits_org_month ON credits(org_id, month);
CREATE INDEX idx_credits_expires_at ON credits(expires_at);
CREATE INDEX idx_tokens_collab_month ON tokens(collab_id, month);
CREATE INDEX idx_redemptions_offer_id ON redemptions(offer_id);
CREATE INDEX idx_redemptions_org_id ON redemptions(org_id);
CREATE INDEX idx_redemptions_collab_id ON redemptions(collab_id);
CREATE INDEX idx_redemptions_confirmed_at ON redemptions(confirmed_at);
CREATE INDEX idx_remote_assignments_active ON remote_org_assignments(active);
CREATE INDEX idx_events_entity ON events(entity_type, entity_id);
CREATE INDEX idx_events_created_at ON events(created_at);

-- Create updated_at trigger function
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply updated_at triggers to relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_profiles_updated_at BEFORE UPDATE ON profiles
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_organizations_updated_at BEFORE UPDATE ON organizations
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_collaborators_updated_at BEFORE UPDATE ON collaborators
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_donation_offers_updated_at BEFORE UPDATE ON donation_offers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();