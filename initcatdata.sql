/*
 Navicat Premium Data Transfer

 Source Server         : cat
 Source Server Type    : MySQL
 Source Server Version : 50740 (5.7.40)
 Source Host           : 127.0.0.1:43306
 Source Schema         : data_cat

 Target Server Type    : MySQL
 Target Server Version : 50740 (5.7.40)
 File Encoding         : 65001

 Date: 24/07/2023 15:17:45
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for boss_list
-- ----------------------------
DROP TABLE IF EXISTS `boss_list`;
CREATE TABLE `boss_list` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `boss_id` varchar(40) DEFAULT NULL,
  `hp` int(11) DEFAULT '5000' COMMENT '血量',
  `state` tinyint(2) DEFAULT '1' COMMENT '1=正常 2=删除 3=执行中',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11027 DEFAULT CHARSET=utf8mb4;

-- ----------------------------
-- Records of boss_list
-- ----------------------------

-- ----------------------------
-- Table structure for config
-- ----------------------------
DROP TABLE IF EXISTS `config`;
CREATE TABLE `config` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `conf_name` varchar(255) NOT NULL DEFAULT '',
  `conf_key` varchar(255) NOT NULL DEFAULT '',
  `conf_value` varchar(800) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`),
  UNIQUE KEY `conf_key` (`conf_key`)
) ENGINE=InnoDB AUTO_INCREMENT=979 DEFAULT CHARSET=utf8;

-- ----------------------------
-- Records of config
-- ----------------------------
BEGIN;
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (1, '拉动物uid列表', 'animalUids', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (2, '牛牛Boss组1', 'cowBoss1', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (3, '蜜蜜Boss组1', 'mmBoss1', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (4, '牛大哥', 'cowBoy', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (5, '蜜蜜', 'mm', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (6, '拉动物定时器状态', 'pullAnimalGoStatus', '1');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (7, '微信登录二维码', 'wechatLoginQrcode', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (8, 'BOSS组3', 'boss3', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (9, '摇一摇ids', 'drawIds', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (10, '新BOSS组1', 'newBoss1', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (11, '新BOSS组2', 'newBoss2', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (12, '新BOSS组3', 'newBoss3', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (13, '摇一摇状态', 'drawStatus', '1');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (14, '挨打小号', 'attackIslandUid', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (15, '挨打小号2', 'attackIslandUid2', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (17, '牛OPENID', 'openid', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (18, '摇一摇备份', 'drawidsbakup', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (19, '汤圆领取定时器', 'steamBoxStatus', '1');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (20, '海滩助力铲子uid列表', 'beachUidList', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (21, '摇一摇切宠物', 'drawChangePet', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (22, '今日动物统计', 'todayAnimalsData', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (23, '今日已计算动物列表', 'todayAlreadyCalculateAnimals', '[]');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (24, '昨日动物统计', 'yesterdayAnimalsData', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (25, '今日敌方动物统计', 'enemyAnimalsData', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (26, '昨日敌方动物统计', 'enemyYesterdayAnimalsData', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (27, '动物待计算列表', 'animalList', '');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (28, 'animalUid', 'animalUid', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (29, 'todayInitAnimal', 'todayInitAnimal', '{}');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (30, 'crawlerStatus', 'crawlerStatus', '1');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (31, '检测账号时是否互送拼图', 'checkPiece', '1');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (32, '海滩助力uids', 'beachHelpUids', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (33, 'beachStatus', 'beachStatus', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (34, '是否可以运行下一个Runner', 'isRunDone', '1');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (35, '检测token定时器开关', 'checkTokenStatus', '1');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (36, 'MaxDraw', 'maxDraw', '1995');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (304, '不允许删除的账号', 'cannotdeleteusers', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (306, '外部账号海滩需要助力配置', 'outbeachneedhelps', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (910, '让公会拉动物', 'animal_family_ids', '43000');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (934, '每日6点30是否领取汤圆等活动奖励', 'exchangeRiceCakeStatus', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (945, '幸运多宝状态', 'PlayLuckyWheelGo', '0');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (946, '周六是否攻击龙？', 'SaturdayAttackMyBossStatus', '1');
INSERT INTO `config` (`id`, `conf_name`, `conf_key`, `conf_value`) VALUES (952, '打完龙是否直接领取奖励', 'AttackBossGoGetbossPrizeLogic', '0');
COMMIT;

-- ----------------------------
-- Table structure for tokens
-- ----------------------------
DROP TABLE IF EXISTS `tokens`;
CREATE TABLE `tokens` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `token` varchar(80) COLLATE utf8mb4_unicode_ci NOT NULL,
  `serverURL` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `zoneToken` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `pull_rows` varchar(80) COLLATE utf8mb4_unicode_ci DEFAULT '1,2,3,4,5,6',
  `init_animals` varchar(500) COLLATE utf8mb4_unicode_ci DEFAULT NULL,
  `all_animals` text COLLATE utf8mb4_unicode_ci,
  `password` varchar(20) COLLATE utf8mb4_unicode_ci DEFAULT 'Aa112211',
  `beach_runner` int(11) DEFAULT '0',
  `add_firewood_types` varchar(10) COLLATE utf8mb4_unicode_ci DEFAULT '2,3',
  `hit_boss_nums` int(11) DEFAULT '5' COMMENT '打龙的炮弹',
  `update_time` timestamp NULL DEFAULT NULL ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
  `create_time` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `state` tinyint(4) DEFAULT '1',
  `follow_uids` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '跟随拉动物',
  `familyId` varchar(255) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '公会id',
  `draw_status` int(11) DEFAULT NULL COMMENT '摇能量状态',
  `attack_sub_id` int(10) DEFAULT NULL COMMENT '攻击小号的岛',
  `is_show` tinyint(4) NOT NULL DEFAULT '1' COMMENT '1=显示 0=不显示',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=697936103 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ----------------------------
-- Records of tokens
-- ----------------------------

SET FOREIGN_KEY_CHECKS = 1;
