CREATE TABLE `users` (
 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
 `email` varchar(255) NOT NULL,
 `name` varchar(255) NOT NULL,
 `dropbox_token` varchar(255) DEFAULT NULL,
 `drive_token` varchar(255) DEFAULT NULL,
 PRIMARY KEY (`id`),
 UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8

CREATE TABLE `links` (
 `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
 `user_id` int(10) unsigned NOT NULL,
 `title` varchar(255) CHARACTER SET utf8mb4 NOT NULL,
 `password` varchar(255) NOT NULL,
 `slug` varchar(255) NOT NULL,
 `description` text CHARACTER SET utf8mb4 NOT NULL,
 `deadline` datetime NOT NULL,
 PRIMARY KEY (`id`),
 UNIQUE KEY `slug` (`slug`),
 KEY `links_user_id_users_id_foreign` (`user_id`),
 CONSTRAINT `links_user_id_users_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8