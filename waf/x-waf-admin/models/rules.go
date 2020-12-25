package models

import "log"

type Rules struct {
	Id       int64
	RuleType string `xorm:"varchar(25) not null"`
	RuleItem string `xorm:"varchar(255) not null"`
}

var RuleInfo = map[string]string{
	"args":      "args规则",
	"blackip":   "黑名单",
	"whiteip":   "白名单",
	"cookie":    "cookies规则",
	"url":       "url规则",
	"useragent": "UserAgent规则",
	"headers":   "headers规则",
	"whiteUrl":  "URL白名单",
	"post":      "post规则",
}

const DefaultRules = `INSERT INTO rules VALUES (4,'args','\\.\\./'),(2,'blackip','8.8.8.8'),(3,'blackip','1.1.1.1'),(5,'args','\\:\\$'),(6,'args','\\$\\{'),(7,'args','select.+(from|limit)'),(8,'args','(?:(union(.*?)select))'),(9,'args','having|rongjitest'),(10,'args','sleep\\((\\s*)(\\d*)(\\s*)\\)'),(11,'args','benchmark\\((.*)\\,(.*)\\)'),(12,'args','base64_decode\\('),(13,'args','(?:from\\W+information_schema\\W)'),(14,'args','(?:(?:current_)user|database|schema|connection_id)\\s*\\('),(15,'args','(?:etc\\/\\W*passwd)'),(16,'args','into(\\s+)+(?:dump|out)file\\s*'),(17,'args','group\\s+by.+\\('),(18,'args','xwork.MethodAccessor'),(19,'args','(?:define|eval|file_get_contents|include|require|require_once|shell_exec|phpinfo|system|passthru|preg_\\w+|execute|echo|print|print_r|var_dump|(fp)open|alert|showmodaldialog)\\('),(20,'args','xwork\\.MethodAccessor'),(21,'args','(gopher|doc|php|glob|file|phar|zlib|ftp|ldap|dict|ogg|data)\\:\\/'),(22,'args','java\\.lang'),(23,'args','\\$_(GET|post|cookie|files|session|env|phplib|GLOBALS|SERVER)\\['),(24,'args','\\<(iframe|script|body|img|layer|div|meta|style|base|object|input)'),(25,'args','(onmouseover|onerror|onload)\\='),(26,'cookie','\\.\\./'),(27,'cookie','\\:\\$'),(28,'cookie','\\$\\{'),(29,'cookie','select.+(from|limit)'),(30,'cookie','(?:(union(.*?)select))'),(31,'cookie','having|rongjitest'),(32,'cookie','sleep\\((\\s*)(\\d*)(\\s*)\\)'),(33,'cookie','benchmark\\((.*)\\,(.*)\\)'),(34,'cookie','base64_decode\\('),(35,'cookie','(?:from\\W+information_schema\\W)'),(36,'cookie','(?:(?:current_)user|database|schema|connection_id)\\s*\\('),(37,'cookie','(?:etc\\/\\W*passwd)'),(38,'cookie','into(\\s+)+(?:dump|out)file\\s*'),(39,'cookie','group\\s+by.+\\('),(40,'cookie','xwork.MethodAccessor'),(41,'cookie','(?:define|eval|file_get_contents|include|require|require_once|shell_exec|phpinfo|system|passthru|preg_\\w+|execute|echo|print|print_r|var_dump|(fp)open|alert|showmodaldialog)\\('),(42,'cookie','xwork\\.MethodAccessor'),(43,'cookie','(gopher|doc|php|glob|file|phar|zlib|ftp|ldap|dict|ogg|data)\\:\\/'),(44,'cookie','java\\.lang'),(45,'cookie','\\$_(GET|post|cookie|files|session|env|phplib|GLOBALS|SERVER)\\['),(46,'post','\\.\\./'),(47,'post','select.+(from|limit)'),(48,'post','(?:(union(.*?)select))'),(49,'post','having|rongjitest'),(50,'post','sleep\\((\\s*)(\\d*)(\\s*)\\)'),(51,'post','benchmark\\((.*)\\,(.*)\\)'),(52,'post','base64_decode\\('),(53,'post','(?:from\\W+information_schema\\W)'),(54,'post','(?:(?:current_)user|database|schema|connection_id)\\s*\\('),(55,'post','(?:etc\\/\\W*passwd)'),(56,'post','into(\\s+)+(?:dump|out)file\\s*'),(57,'post','group\\s+by.+\\('),(58,'post','xwork.MethodAccessor'),(59,'post','(?:define|eval|file_get_contents|include|require|require_once|shell_exec|phpinfo|system|passthru|preg_\\w+|execute|echo|print|print_r|var_dump|(fp)open|alert|showmodaldialog)\\('),(60,'post','xwork\\.MethodAccessor'),(61,'post','(gopher|doc|php|glob|file|phar|zlib|ftp|ldap|dict|ogg|data)\\:\\/'),(62,'post','java\\.lang'),(63,'post','\\$_(GET|post|cookie|files|session|env|phplib|GLOBALS|SERVER)\\['),(64,'post','\\<(iframe|script|body|img|layer|div|meta|style|base|object|input)'),(65,'post','(onmouseover|onerror|onload)\\='),(66,'url','\\.(htaccess|bash_history)'),(67,'url','\\.(bak|inc|old|mdb|sql|backup|java|class|tgz|gz|tar|zip)$'),(68,'url','(phpmyadmin|jmx-console|admin-console|jmxinvokerservlet)'),(69,'url','java\\.lang'),(70,'url','\\.(svn|git|sql|bak)\\/'),(71,'url','/(attachments|upimg|images|css|uploadfiles|html|uploads|templets|static|template|data|inc|forumdata|upload|includes|cache|avatar)/(\\\\w+).(php|jsp)'),(72,'useragent','(HTTrack|harvest|audit|dirbuster|pangolin|nmap|sqln|-scan|hydra|Parser|libwww|BBBike|sqlmap|w3af|owasp|Nikto|fimap|havij|PycURL|zmeu|BabyKrokodil|netsparker|httperf|bench)'),(73,'whiteUrl','/news/'),(74,'whiteip','8.8.8.8'),(75,'args','and\\s+(1=1|1=2)');`

func ListRules() (rules []Rules, err error) {
	rules = make([]Rules, 0)
	err = Engine.Find(&rules)
	log.Println(err, rules)
	return rules, err
}

// list rules by rule_type
func ListRulesByType(ruleType string) (rules []Rules, err error) {
	rules = make([]Rules, 0)
	err = Engine.Where("rule_type = ?", ruleType).Find(&rules)
	return rules, err
}

func ListAllRules() (rulesMap map[string][]Rules, err error) {
	rulesMap = make(map[string][]Rules)
	rules := make([]Rules, 0)
	err = Engine.Find(&rules)
	for _, item := range rules {
		rulesMap[item.RuleType] = append(rulesMap[item.RuleType], item)
	}
	return rulesMap, err
}

func NewRule(ruleType string, rule string) (err error) {
	_, err = Engine.Insert(&Rules{RuleType: ruleType, RuleItem: rule})
	return err
}

func EditRule(ruleId int64, ruleItem string) (err error) {
	Rule := new(Rules)
	Engine.Id(ruleId).Get(Rule)
	Rule.RuleItem = ruleItem
	_, err = Engine.Id(ruleId).Update(Rule)
	return err
}

func DelRule(ruleId int64) (err error) {
	_, err = Engine.Delete(&Rules{Id: ruleId})
	return err
}
