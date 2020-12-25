-- MySQL dump 10.13  Distrib 5.6.30, for debian-linux-gnu (x86_64)
--
-- Host: 10.99.184.139    Database: lottery
-- ------------------------------------------------------
-- Server version	5.1.66

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `rules`
--

DROP TABLE IF EXISTS `rules`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `rules` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `rule_type` varchar(25) NOT NULL,
  `rule_item` varchar(255) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=MyISAM AUTO_INCREMENT=75 DEFAULT CHARSET=utf8;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `rules`
--

LOCK TABLES `rules` WRITE;
/*!40000 ALTER TABLE `rules` DISABLE KEYS */;
INSERT INTO `rules` VALUES (4,'args','\\.\\./'),(2,'blackip','8.8.8.8'),(3,'blackip','1.1.1.1'),(5,'args','\\:\\$'),(6,'args','\\$\\{'),(7,'args','select.+(from|limit)'),(8,'args','(?:(union(.*?)select))'),(9,'args','having|rongjitest'),(10,'args','sleep\\((\\s*)(\\d*)(\\s*)\\)'),(11,'args','benchmark\\((.*)\\,(.*)\\)'),(12,'args','base64_decode\\('),(13,'args','(?:from\\W+information_schema\\W)'),(14,'args','(?:(?:current_)user|database|schema|connection_id)\\s*\\('),(15,'args','(?:etc\\/\\W*passwd)'),(16,'args','into(\\s+)+(?:dump|out)file\\s*'),(17,'args','group\\s+by.+\\('),(18,'args','xwork.MethodAccessor'),(19,'args','(?:define|eval|file_get_contents|include|require|require_once|shell_exec|phpinfo|system|passthru|preg_\\w+|execute|echo|print|print_r|var_dump|(fp)open|alert|showmodaldialog)\\('),(20,'args','xwork\\.MethodAccessor'),(21,'args','(gopher|doc|php|glob|file|phar|zlib|ftp|ldap|dict|ogg|data)\\:\\/'),(22,'args','java\\.lang'),(23,'args','\\$_(GET|post|cookie|files|session|env|phplib|GLOBALS|SERVER)\\['),(24,'args','\\<(iframe|script|body|img|layer|div|meta|style|base|object|input)'),(25,'args','(onmouseover|onerror|onload)\\='),(26,'cookie','\\.\\./'),(27,'cookie','\\:\\$'),(28,'cookie','\\$\\{'),(29,'cookie','select.+(from|limit)'),(30,'cookie','(?:(union(.*?)select))'),(31,'cookie','having|rongjitest'),(32,'cookie','sleep\\((\\s*)(\\d*)(\\s*)\\)'),(33,'cookie','benchmark\\((.*)\\,(.*)\\)'),(34,'cookie','base64_decode\\('),(35,'cookie','(?:from\\W+information_schema\\W)'),(36,'cookie','(?:(?:current_)user|database|schema|connection_id)\\s*\\('),(37,'cookie','(?:etc\\/\\W*passwd)'),(38,'cookie','into(\\s+)+(?:dump|out)file\\s*'),(39,'cookie','group\\s+by.+\\('),(40,'cookie','xwork.MethodAccessor'),(41,'cookie','(?:define|eval|file_get_contents|include|require|require_once|shell_exec|phpinfo|system|passthru|preg_\\w+|execute|echo|print|print_r|var_dump|(fp)open|alert|showmodaldialog)\\('),(42,'cookie','xwork\\.MethodAccessor'),(43,'cookie','(gopher|doc|php|glob|file|phar|zlib|ftp|ldap|dict|ogg|data)\\:\\/'),(44,'cookie','java\\.lang'),(45,'cookie','\\$_(GET|post|cookie|files|session|env|phplib|GLOBALS|SERVER)\\['),(46,'post','\\.\\./'),(47,'post','select.+(from|limit)'),(48,'post','(?:(union(.*?)select))'),(49,'post','having|rongjitest'),(50,'post','sleep\\((\\s*)(\\d*)(\\s*)\\)'),(51,'post','benchmark\\((.*)\\,(.*)\\)'),(52,'post','base64_decode\\('),(53,'post','(?:from\\W+information_schema\\W)'),(54,'post','(?:(?:current_)user|database|schema|connection_id)\\s*\\('),(55,'post','(?:etc\\/\\W*passwd)'),(56,'post','into(\\s+)+(?:dump|out)file\\s*'),(57,'post','group\\s+by.+\\('),(58,'post','xwork.MethodAccessor'),(59,'post','(?:define|eval|file_get_contents|include|require|require_once|shell_exec|phpinfo|system|passthru|preg_\\w+|execute|echo|print|print_r|var_dump|(fp)open|alert|showmodaldialog)\\('),(60,'post','xwork\\.MethodAccessor'),(61,'post','(gopher|doc|php|glob|file|phar|zlib|ftp|ldap|dict|ogg|data)\\:\\/'),(62,'post','java\\.lang'),(63,'post','\\$_(GET|post|cookie|files|session|env|phplib|GLOBALS|SERVER)\\['),(64,'post','\\<(iframe|script|body|img|layer|div|meta|style|base|object|input)'),(65,'post','(onmouseover|onerror|onload)\\='),(66,'url','\\.(htaccess|bash_history)'),(67,'url','\\.(bak|inc|old|mdb|sql|backup|java|class|tgz|gz|tar|zip)$'),(68,'url','(phpmyadmin|jmx-console|admin-console|jmxinvokerservlet)'),(69,'url','java\\.lang'),(70,'url','\\.(svn|git|sql|bak)\\/'),(71,'url','/(attachments|upimg|images|css|uploadfiles|html|uploads|templets|static|template|data|inc|forumdata|upload|includes|cache|avatar)/(\\\\w+).(php|jsp)'),(72,'useragent','(HTTrack|harvest|audit|dirbuster|pangolin|nmap|sqln|-scan|hydra|Parser|libwww|BBBike|sqlmap|w3af|owasp|Nikto|fimap|havij|PycURL|zmeu|BabyKrokodil|netsparker|httperf|bench)'),(73,'whiteUrl','/news/'),(74,'whiteip','8.8.8.8');
/*!40000 ALTER TABLE `rules` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2016-08-06 13:23:11
