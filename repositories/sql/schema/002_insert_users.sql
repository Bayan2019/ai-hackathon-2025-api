-- +goose Up
INSERT INTO users(first_name, 
                last_name, 
                email, password_hash)
    VALUES ('Bayan', 
            'Saparbayeva',
            'sapar1986@yahoo.com',
            '$2a$12$xQ6dkljWOaEDFgHoEzDK0.qGdUHcXiuZS0kEVoTiWsUaHSa3pYmBS'),
        ('Amina',
            'Yeszhanova', 
            'amina.yes03@gmail.com',
            '$2a$12$xQ6dkljWOaEDFgHoEzDK0.qGdUHcXiuZS0kEVoTiWsUaHSa3pYmBS');

-- +goose Down
DELETE FROM users;