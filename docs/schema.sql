-- NeoBarter 数据库 Schema (重新设计)
-- PostgreSQL 15+

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- ========================================
-- 用户相关
-- ========================================

-- 用户表
CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    phone VARCHAR(20) UNIQUE NOT NULL,
    nickname VARCHAR(50),
    avatar_url VARCHAR(255),
    user_type VARCHAR(10) NOT NULL DEFAULT 'personal', -- personal / enterprise
    status VARCHAR(20) NOT NULL DEFAULT 'active',      -- active / banned / deleted
    credit_score INTEGER NOT NULL DEFAULT 100,
    real_name VARCHAR(50),
    id_card VARCHAR(30),
    real_name_verified BOOLEAN NOT NULL DEFAULT FALSE,
    enterprise_name VARCHAR(100),
    enterprise_license_url VARCHAR(255),
    enterprise_verified BOOLEAN NOT NULL DEFAULT FALSE,
    location VARCHAR(100),
    bio TEXT,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_phone ON users(phone);
CREATE INDEX idx_users_status ON users(status);

-- 用户收货地址
CREATE TABLE user_addresses (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(50) NOT NULL,
    phone VARCHAR(20) NOT NULL,
    province VARCHAR(30) NOT NULL,
    city VARCHAR(30) NOT NULL,
    district VARCHAR(30) NOT NULL,
    detail VARCHAR(200) NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_user_addresses_user ON user_addresses(user_id);

-- ========================================
-- 钱包 / 巴特币
-- ========================================

-- 钱包表（每用户一个）
CREATE TABLE wallets (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    balance DECIMAL(12,2) NOT NULL DEFAULT 0.00,  -- 巴特币余额
    frozen_balance DECIMAL(12,2) NOT NULL DEFAULT 0.00, -- 冻结余额（交易中）
    total_income DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    total_expense DECIMAL(12,2) NOT NULL DEFAULT 0.00,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_balance_non_negative CHECK (balance >= 0),
    CONSTRAINT chk_frozen_non_negative CHECK (frozen_balance >= 0)
);

-- 钱包流水
CREATE TABLE wallet_transactions (
    id BIGSERIAL PRIMARY KEY,
    wallet_id BIGINT NOT NULL REFERENCES wallets(id),
    type VARCHAR(20) NOT NULL,        -- reward / trade_in / trade_out / deposit / withdraw / freeze / unfreeze
    amount DECIMAL(12,2) NOT NULL,
    balance_after DECIMAL(12,2) NOT NULL,
    reference_type VARCHAR(30),       -- trade_request / system / usdc
    reference_id BIGINT,
    description VARCHAR(200),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_wallet_tx_wallet ON wallet_transactions(wallet_id);
CREATE INDEX idx_wallet_tx_type ON wallet_transactions(type);
CREATE INDEX idx_wallet_tx_created ON wallet_transactions(created_at);

-- ========================================
-- 物品相关
-- ========================================

-- 物品分类表
CREATE TABLE categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL,
    parent_id INTEGER REFERENCES categories(id),
    icon VARCHAR(50),
    sort_order INTEGER DEFAULT 0
);

-- 物品表
CREATE TABLE items (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(100) NOT NULL,
    description TEXT,
    category_id INTEGER REFERENCES categories(id),
    estimated_value DECIMAL(10,2),     -- 巴特币估值
    condition VARCHAR(20) NOT NULL DEFAULT 'good', -- new / like_new / good / fair
    images TEXT[],
    video_url VARCHAR(255),
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active / inactive / traded / deleted
    location VARCHAR(100),
    view_count INTEGER NOT NULL DEFAULT 0,
    want_items TEXT[],                 -- 期望交换的物品描述
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_items_user ON items(user_id);
CREATE INDEX idx_items_category ON items(category_id);
CREATE INDEX idx_items_status ON items(status);
CREATE INDEX idx_items_created ON items(created_at DESC);

-- ========================================
-- 交易相关
-- ========================================

-- 交换请求表
CREATE TABLE trade_requests (
    id BIGSERIAL PRIMARY KEY,
    initiator_id BIGINT NOT NULL REFERENCES users(id),
    target_user_id BIGINT NOT NULL REFERENCES users(id),
    target_item_id BIGINT NOT NULL REFERENCES items(id),
    offered_item_id BIGINT REFERENCES items(id),       -- 可为空（纯巴特币购买）
    barter_coin_amount DECIMAL(10,2) DEFAULT 0.00,     -- 补差价的巴特币
    status VARCHAR(20) NOT NULL DEFAULT 'pending',     -- pending / accepted / rejected / completed / cancelled / expired
    message TEXT,
    reject_reason TEXT,
    expired_at TIMESTAMP,
    completed_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_trade_initiator ON trade_requests(initiator_id);
CREATE INDEX idx_trade_target_user ON trade_requests(target_user_id);
CREATE INDEX idx_trade_status ON trade_requests(status);

-- ========================================
-- 消息相关
-- ========================================

-- 会话表
CREATE TABLE conversations (
    id BIGSERIAL PRIMARY KEY,
    type VARCHAR(20) NOT NULL DEFAULT 'private', -- private / group
    last_message_id BIGINT,
    last_message_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- 会话参与者
CREATE TABLE conversation_participants (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    unread_count INTEGER NOT NULL DEFAULT 0,
    last_read_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(conversation_id, user_id)
);

CREATE INDEX idx_conv_part_user ON conversation_participants(user_id);

-- 消息表
CREATE TABLE messages (
    id BIGSERIAL PRIMARY KEY,
    conversation_id BIGINT NOT NULL REFERENCES conversations(id) ON DELETE CASCADE,
    sender_id BIGINT NOT NULL REFERENCES users(id),
    content TEXT NOT NULL,
    message_type VARCHAR(20) NOT NULL DEFAULT 'text', -- text / image / voice / item_card
    extra_data JSONB,                                  -- 附加数据（图片URL、物品卡片等）
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_messages_conv ON messages(conversation_id, created_at DESC);
CREATE INDEX idx_messages_sender ON messages(sender_id);

-- ========================================
-- 评价相关
-- ========================================

-- 评价表
CREATE TABLE reviews (
    id BIGSERIAL PRIMARY KEY,
    trade_request_id BIGINT NOT NULL REFERENCES trade_requests(id),
    reviewer_id BIGINT NOT NULL REFERENCES users(id),
    reviewee_id BIGINT NOT NULL REFERENCES users(id),
    rating INTEGER NOT NULL CHECK (rating >= 1 AND rating <= 5),
    comment TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(trade_request_id, reviewer_id)
);

CREATE INDEX idx_reviews_reviewee ON reviews(reviewee_id);

-- ========================================
-- 通知相关
-- ========================================

-- 通知表
CREATE TABLE notifications (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    type VARCHAR(30) NOT NULL,         -- trade_request / trade_accepted / trade_rejected / message / system
    title VARCHAR(100) NOT NULL,
    content TEXT,
    reference_type VARCHAR(30),
    reference_id BIGINT,
    is_read BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_notifications_user ON notifications(user_id, is_read, created_at DESC);

-- ========================================
-- 初始数据
-- ========================================

INSERT INTO categories (name, icon, sort_order) VALUES
('数码电子', 'laptop', 1),
('家用电器', 'home', 2),
('服饰鞋包', 'shopping', 3),
('图书教材', 'book', 4),
('美妆护肤', 'heart', 5),
('运动户外', 'sports', 6),
('家居家具', 'shop', 7),
('母婴用品', 'gift', 8),
('食品饮料', 'coffee', 9),
('其他', 'ellipsis', 99);
