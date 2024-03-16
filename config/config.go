package config

import (
	"flag"

	"github.com/vharitonsky/iniflags"
)

var (
	AssetsPath 	  = flag.String("assets", "./assets", "Assets path for images")
	Env           = flag.String("environment", "local", "Environment")
	DB            = flag.String("db", "", "Database connection url")
	LoginSize     = flag.Int("login_size", 6, "Login size")
	Port          = flag.Int("port", 8082, "Port for running application")
	RefreshSecret = flag.String("refresh_secret", "§1puR?*1f", "Secret used for refresh jwt")
	TokenSecret   = flag.String("token_secret", "}Sn+#u||j?9n", "Secret used for access jwt")
	TokenLife     = flag.Int64("token_life", 3600, "Token life jwt")
	TokenSize     = flag.Int("token_size", 30, "Token size jwt")
	PasswordSize  = flag.Int("password_size", 6, "Password size")
	ProxyBluetooth = flag.String("proxy", "192.168.1.17:2000", "Proxy bluetooth")
	SMTPFrom      = flag.String("email_from", "admin@neli.com", "Sender address for neli email")
	SMTPLogin     = flag.String("email_login", "", "SMTP login for email authentication")
	SMTPPassword  = flag.String("email_password", "", "SMTP password for email authentication")
	SMTPHost      = flag.String("email_host", "localhost", "SMTP host for email authentication")
	SMTPPort      = flag.Int("email_port", 1025, "SMTP port for email authentication")
	Stub          = flag.Int("stub", 0, "Stub in case of hub unavailable. Automatically assign a string to a content.")
	URL           = flag.String("url", "http://neli", "Base url for forms")
	URLContent    = flag.String("url_content", "http://neli", "Base url for content (video and image)")
	ZombieID      = flag.Int64("zombie", 99, "Zombie id")
	StubPath      = "123456"
)

func init() {
	iniflags.Parse()
}
