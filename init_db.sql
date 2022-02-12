CREATE TABLE `tokens` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `uid` bigint(20) NOT NULL DEFAULT '0',
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `token` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL,
  `serverURL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zoneToken` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pull_rows` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT '1,2,3,4,5,6',
  `init_animals` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `all_animals` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `beach_runner` tinyint(4) DEFAULT '0',
  `password` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'Aa112211',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=697625166 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
