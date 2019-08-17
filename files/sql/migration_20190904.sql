CREATE TABLE `user_storage_credentials` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(10) unsigned NOT NULL,
  `provider_id` int(10) unsigned NOT NULL,
  `provider_credential` varchar(255) NOT NULL,
  `email` varchar(255) NOT NULL DEFAULT '',
  `photo` varchar(255) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `user_provider_unique` (`user_id`, `provider_id`),
  KEY `usc_user_id_users_id_foreign` (`user_id`),
  CONSTRAINT `usc_user_id_users_id_foreign` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE ON UPDATE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8

ALTER TABLE `links`
ADD `user_storage_credential_id` int(10) unsigned NULL,
ADD KEY `links_usc_id_foreign` (`user_storage_credential_id`),
ADD CONSTRAINT `links_usc_id_foreign` FOREIGN KEY (`user_storage_credential_id`) REFERENCES `user_storage_credentials` (`id`) ON DELETE SET NULL ON UPDATE CASCADE;