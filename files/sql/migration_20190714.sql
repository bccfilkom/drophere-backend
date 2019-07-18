ALTER TABLE `users`
ADD `recover_password_token` varchar(255) NULL,
ADD `recover_password_token_expiry` datetime NULL;