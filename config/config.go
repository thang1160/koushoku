package config

import (
	_ "embed"

	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	"github.com/jessevdk/go-flags"
	"gopkg.in/ini.v1"
)

//go:embed config.ini
var buf []byte

var opts struct {
	Path string `short:"c" long:"config" description:"Path to config file"`
	Mode string `short:"m" long:"mode" description:"App mode"`
}

var Config struct {
	file *ini.File
	mu   sync.RWMutex

	Mode string

	Meta struct {
		BaseURL     string `json:"baseURL"`
		Title       string `json:"title"`
		Description string `json:"description"`
		Language    string `json:"language"`
	}

	Database struct {
		Host    string
		Port    int
		Name    string
		User    string
		Passwd  string
		SSLMode string
	}

	Redis struct {
		Host   string
		Port   int
		DB     int
		Passwd string
	}

	Server struct {
		Port int
	}

	Cache struct {
		DefaultTTL   time.Duration
		TemplatesTTL time.Duration
	}

	Directories struct {
		Root string
		Data string

		Symlinks   string
		Thumbnails string
		Covers     string
	}

	Paths struct {
		Batches   string
		Singles   string
		Blacklist string
	}
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	exec, err := os.Executable()
	if err != nil {
		log.Fatalln(err)
	} else if exec, err = filepath.EvalSymlinks(exec); err != nil {
		log.Fatalln(err)
	}

	flags.NewParser(&opts, flags.PassDoubleDash).Parse()

	Config.Directories.Root = filepath.Dir(exec)
	Config.Directories.Symlinks = filepath.Join(Config.Directories.Root, "symlinks")
	Config.Directories.Thumbnails = filepath.Join(Config.Directories.Root, "thumbnails")
	Config.Directories.Covers = filepath.Join(Config.Directories.Root, "covers")

	Config.Paths.Batches = filepath.Join(Config.Directories.Root, "batches.txt")
	Config.Paths.Singles = filepath.Join(Config.Directories.Root, "singles.txt")
	Config.Paths.Blacklist = filepath.Join(Config.Directories.Root, "blacklist.txt")

	if len(opts.Path) == 0 {
		opts.Path = filepath.Join(Config.Directories.Root, "config.ini")
	}
	_, err = os.Stat(opts.Path)
	if os.IsNotExist(err) {
		dir := filepath.Dir(opts.Path)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				log.Fatalln(err)
			}
		}
		if err := os.WriteFile(opts.Path, buf, 0755); err != nil {
			log.Fatalln(err)
		}
	}

	var file *ini.File
	if file, err = ini.Load(opts.Path); err != nil {
		log.Fatalln(err)
	}

	Config.file = file
	Config.Mode = file.Section("").Key("mode").MustString("production")

	Config.Meta.BaseURL = file.Section("meta").Key("base_url").MustString("http://localhost:42073")
	Config.Meta.Title = file.Section("meta").Key("title").MustString("Koushoku")
	Config.Meta.Description = file.Section("meta").Key("description").String()
	Config.Meta.Language = file.Section("meta").Key("language").MustString("en-US")

	Config.Database.Host = file.Section("database").Key("host").MustString("localhost")
	Config.Database.Port = file.Section("database").Key("port").MustInt(5432)
	Config.Database.Name = file.Section("database").Key("name").MustString("koushoku")
	Config.Database.User = file.Section("database").Key("user").MustString("koushoku")
	Config.Database.Passwd = file.Section("database").Key("passwd").MustString("koushoku")
	Config.Database.SSLMode = file.Section("database").Key("ssl_mode").MustString("disable")

	Config.Redis.Host = file.Section("redis").Key("host").MustString("localhost")
	Config.Redis.Port = file.Section("redis").Key("port").MustInt(6379)
	Config.Redis.DB = file.Section("redis").Key("db").MustInt(0)
	Config.Redis.Passwd = file.Section("redis").Key("passwd").String()

	Config.Server.Port = file.Section("server").Key("port").MustInt(42073)

	Config.Cache.DefaultTTL = time.Duration(file.Section("cache").Key("default_ttl").MustInt(86400000000000))
	Config.Cache.TemplatesTTL = time.Duration(file.Section("cache").Key("templates_ttl").MustInt(300000000000))

	Config.Directories.Data = file.Section("directories").Key("data").MustString(filepath.Join(Config.Directories.Root, "data"))

	if len(opts.Mode) > 0 {
		Config.Mode = opts.Mode
	}

	Save()
}

func Save() error {
	Config.mu.Lock()
	defer Config.mu.Unlock()

	Config.file.Section("").Key("mode").SetValue(Config.Mode)

	Config.file.Section("meta").Key("base_url").SetValue(Config.Meta.BaseURL)
	Config.file.Section("meta").Key("description").SetValue(Config.Meta.Description)
	Config.file.Section("meta").Key("title").SetValue(Config.Meta.Title)
	Config.file.Section("meta").Key("language").SetValue(Config.Meta.Language)

	Config.file.Section("database").Key("host").SetValue(Config.Database.Host)
	Config.file.Section("database").Key("port").SetValue(strconv.Itoa(Config.Database.Port))
	Config.file.Section("database").Key("name").SetValue(Config.Database.Name)
	Config.file.Section("database").Key("user").SetValue(Config.Database.User)
	Config.file.Section("database").Key("passwd").SetValue(Config.Database.Passwd)
	Config.file.Section("database").Key("ssl_mode").SetValue(Config.Database.SSLMode)

	Config.file.Section("redis").Key("host").SetValue(Config.Redis.Host)
	Config.file.Section("redis").Key("port").SetValue(strconv.Itoa(Config.Redis.Port))
	Config.file.Section("redis").Key("db").SetValue(strconv.Itoa(Config.Redis.DB))
	Config.file.Section("redis").Key("passwd").SetValue(Config.Redis.Passwd)

	Config.file.Section("server").Key("port").SetValue(strconv.Itoa(Config.Server.Port))

	Config.file.Section("cache").Key("default_ttl").SetValue(strconv.Itoa(int(Config.Cache.DefaultTTL)))
	Config.file.Section("cache").Key("templates_ttl").SetValue(strconv.Itoa(int(Config.Cache.TemplatesTTL)))

	Config.file.Section("directories").Key("data").SetValue(Config.Directories.Data)

	return Config.file.SaveTo(opts.Path)
}
