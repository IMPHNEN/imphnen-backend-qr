CREATE TABLE qr_campaigns (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    url TEXT NOT NULL,
    qr_code_data BYTEA NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_qr_campaigns_is_active ON qr_campaigns(is_active);
CREATE UNIQUE INDEX idx_qr_campaigns_single_active ON qr_campaigns(is_active) WHERE is_active = true;
