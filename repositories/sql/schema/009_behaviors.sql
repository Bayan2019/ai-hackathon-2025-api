-- +goose Up
CREATE TABLE behaviors(
    cst_dim_id INTEGER NOT NULL REFERENCES clients(cst_dim_id) ON DELETE CASCADE,
    transdate TEXT NOT NULL,
    monthly_os_changes INTEGER NOT NULL,
    monthly_phone_model_changes INTEGER NOT NULL,
    last_phone_model_categorical TEXT NOT NULL,
    last_os_categorical TEXT NOT NULL,
    logins_last_7_days INTEGER NOT NULL,
    logins_last_30_days INTEGER NOT NULL,
    login_frequency_7d NUMERIC NOT NULL,
    freq_change_7d_vs_mean NUMERIC NOT NULL,
    logins_7d_over_30d_ratio NUMERIC NOT NULL,
    avg_login_interval_30d NUMERIC NOT NULL,
    std_login_interval_30d NUMERIC NOT NULL,
    var_login_interval_30d NUMERIC NOT NULL,
    ewm_login_interval_7d NUMERIC NOT NULL,
    burstiness_login_interval NUMERIC NOT NULL,
    fano_factor_login_interval NUMERIC NOT NULL,
    zscore_avg_login_interval_7d NUMERIC NOT NULL
);

-- +goose Down
DROP TABLE behaviors;