SET character_set_client = utf8;

CREATE TABLE `advert` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `locality` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `link` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `hash_id` int(10) unsigned DEFAULT NULL,
  `price` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `name` varchar(255) COLLATE utf8_unicode_ci DEFAULT NULL,
  `status` enum('new','sent','call','reject','approve') COLLATE utf8_unicode_ci DEFAULT NULL,
  `created` int(10) unsigned DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `hash_id` (`hash_id`)
) ENGINE=InnoDB AUTO_INCREMENT=480 DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci;

CREATE TABLE `location` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `advert_id` int(10) unsigned DEFAULT NULL,
  `lat` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  `lon` varchar(15) COLLATE utf8_unicode_ci DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `advert_id` (`advert_id`),
  CONSTRAINT `fk_advert_id` FOREIGN KEY (`advert_id`) REFERENCES `advert` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci

