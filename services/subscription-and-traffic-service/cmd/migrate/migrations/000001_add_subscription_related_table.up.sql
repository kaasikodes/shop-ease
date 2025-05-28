CREATE TABLE vendor_plans (
    id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    content VARCHAR(350) NOT NULL,
    price INT NOT NULL,
    user_interactions_allowed INT NOT NULL,
    duration_in_secs BIGINT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
);

CREATE TABLE vendor_subscriptions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    plan_id INT NOT NULL,
    vendor_id INT NOT NULL,
    has_paid BOOLEAN NOT NULL DEFAULT FALSE,
    limit_exceeded_at DATETIME NULL,
    paid_at DATETIME NOT NULL,
    began_at DATETIME NOT NULL,
    expires_at DATETIME NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (plan_id) REFERENCES vendor_plans(id)
    -- Add vendor foreign key if vendor table exists
);

CREATE TABLE vendor_user_interactions (
    id INT AUTO_INCREMENT PRIMARY KEY,
    vendor_id INT NOT NULL,
    user_id INT NOT NULL,
    type ENUM('user.ordered_item', 'user.interested_in_item') NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
    -- Add foreign keys if needed
);
