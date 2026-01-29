CREATE TABLE IF NOT EXISTS rentals (
    id BIGSERIAL PRIMARY KEY,
    "order" INTEGER NOT NULL,
    unparsed_address TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_rentals_order ON rentals ("order");
CREATE INDEX IF NOT EXISTS idx_rentals_created_at ON rentals (created_at);

