-- +goose Up
-- +goose StatementBegin
CREATE TYPE service_type AS ENUM (
    'Construction',
    'Delivery',
    'Manufacture'
    );

CREATE TYPE status AS ENUM (
    'Created',
    'Published',
    'Closed'
    );

 CREATE TABLE IF NOT EXISTS tender (
                        id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                        name VARCHAR(100) NOT NULL,
                        description TEXT,
                        service_type service_type NOT NULL ,
                        status status DEFAULT 'Created',
                        organization_id UUID NOT NULL REFERENCES organization(id),
                        creator_username VARCHAR(50) NOT NULL REFERENCES employee(username),
                        version INT DEFAULT 1,
                        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TYPE author_type AS ENUM (
    'Organization',
    'User'
    );


CREATE TYPE bid_status AS ENUM (
    'Created',
    'Published',
    'Canceled'
    );

CREATE TYPE decision AS ENUM (
    'Accepted',
    'Rejected'
    );

 CREATE TABLE IF NOT EXISTS bid (
                                    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                                    name VARCHAR(100) NOT NULL,
                                    description TEXT,
                                    status bid_status NOT NULL DEFAULT 'Created',
                                    tender_id UUID NOT NULL REFERENCES tender(id),
                                    decision decision,
                                    author_type author_type NOT NULL ,
                                    author_id UUID NOT NULL REFERENCES employee(id),
                                    version INT DEFAULT 1,
                                    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
 );
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bid;

DROP TYPE IF EXISTS author_type;

DROP TYPE IF EXISTS bid_status;

DROP TYPE IF EXISTS decision;

DROP TABLE IF EXISTS tender;

DROP TYPE IF EXISTS status;

DROP TYPE IF EXISTS service_type;
-- +goose StatementEnd
