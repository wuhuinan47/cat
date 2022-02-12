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

CREATE TABLE `config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `conf_name` varchar(255) NOT NULL DEFAULT '',
  `conf_key` varchar(255) NOT NULL DEFAULT '',
  `conf_value` varchar(800) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `conf_key` (`conf_key`)
) ENGINE=InnoDB AUTO_INCREMENT=34 DEFAULT CHARSET=utf8;


INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`)
VALUES
	(1, '拉动物uid列表', 'animalUids', '302691822,309392050,301807377,309433834,374289806,375912362,380576240,381034522,381909995,382292124,385498006,403573789,406961861,408385382,408871361,410572648,412305015,412670685,425190502,439943689,440204933,441013912,444729748,445291795,446231345,446399085,447729427,690364007,690423894,690708340,690840661,690980615,692326562,692732133,693419844,693437767,694068717,694316841,694571893,694981971,695021679,695923850,696100351,696280163,696360023,696453763,696528833,696636309,697058069,697068758,697132831,697625165'),
	(2, '牛牛Boss组1', 'cowBoss1', '301807377,309433834,690708340,694068717,693419844,697068758'),
	(3, '蜜蜜Boss组1', 'mmBoss1', '375912362,374289806,439943689,382292124,385498006,381909995'),
	(4, '牛大哥', 'cowBoy', '302691822'),
	(5, '蜜蜜', 'mm', '309392050'),
	(6, '拉动物定时器状态', 'pullAnimalGoStatus', '1'),
	(7, '微信登录二维码', 'wechatLoginQrcode', ''),
	(8, 'BOSS组3', 'boss3', '445291795,381034522,446231345,447729427,694981971,696636309'),
	(9, '摇一摇ids', 'drawIds', '302691822,309392050,301807377,695923850,692326562,690980615,381909995,690423894,694981971,382292124,439943689,309433834,695923850,694571893,696453763,697132831,697058069,696360023,696280163,694316841,374289806,406378614,690364007,375912362,445291795,696636309,696528833,696100351,694068717,693419844,690708340,380576240,381034522,38190,693437767,695021679,697625165,692732133,690840661,697068758,385498006'),
	(10, '新BOSS组1', 'newBoss1', '301807377,309433834,445291795,381034522,446231345,447729427,694981971,696636309'),
	(11, '新BOSS组2', 'newBoss2', '697068758,694068717,693419844,690708340,439943689,385498006,382292124,381909995'),
	(12, '新BOSS组3', 'newBoss3', '374289806,375912362,380576240'),
	(13, '摇一摇状态', 'drawStatus', '1'),
	(14, '挨打小号', 'attackIslandUid', '446231345'),
	(15, '挨打小号2', 'attackIslandUid2', '447729427'),
	(17, '牛OPENID', 'openid', 'o9od753rqL522SlZkuIEc3NVHBKA'),
	(18, '摇一摇备份', 'drawidsbakup', '302691822,309392050,695923850,692326562,690980615,381909995,690423894,694981971,382292124,439943689,309433834,695923850,694571893,696453763,697132831,697058069,696360023,696280163,694316841,374289806,406378614,690364007,375912362,445291795,696636309,696528833,696100351,694068717,693419844,690708340,380576240,381034522,38190,693437767,695021679,697625165,692732133,690840661'),
	(19, '汤圆领取定时器', 'steamBoxStatus', '0'),
	(20, '海滩助力铲子uid列表', 'beachUidList', '301807377,302691822,309392050,309433834,374289806,375912362,380576240,381034522,381909995,382292124,385498006,439943689,445291795,446231345,690708340,693419844,694068717,694981971,696636309,697068758,690364007'),
	(21, '摇一摇切宠物', 'drawChangePet', '1'),
	(22, '今日动物统计', 'todayAnimalsData', ''),
	(23, '今日已计算动物列表', 'todayAlreadyCalculateAnimals', '[]'),
	(24, '昨日动物统计', 'yesterdayAnimalsData', ''),
	(25, '今日敌方动物统计', 'enemyAnimalsData', ''),
	(26, '昨日敌方动物统计', 'enemyYesterdayAnimalsData', ''),
	(27, '动物待计算列表', 'animalList', ''),
	(28, 'animalUid', 'animalUid', '690364007'),
	(29, 'todayInitAnimal', 'todayInitAnimal', '{\"76\":7,\"77\":1,\"78\":1,\"79\":0,\"80\":0,\"81\":0}'),
	(30, 'crawlerStatus', 'crawlerStatus', '1'),
	(31, '检测账号时是否互送拼图', 'checkPiece', '1'),
	(32, '海滩助力uids', 'beachHelpUids', '408385382,408871361,410572648,412305015,412670685,425190502,440204933,441013912,444729748,446399085,690423894,690840661,690980615,692326562,693437767,694571893,695021679,696453763,697058069,697132831,697625165'),
	(33, 'beachStatus', 'beachStatus', '0');
