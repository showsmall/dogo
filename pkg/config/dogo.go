package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type DogoConfig struct {
	Dsn   string
	Addr  []string
	Guacd *Guacd
}

type Guacd struct {
	Addr string
	RDP  *RDP
	SSH  *SSH
	VNC  *VNC
}

type RDP struct {
	DriveName string
	DrivePath string
}

type SSH struct {
	FontName    string
	FontSize    string
	ColorScheme string
}

type VNC struct {
	Autoretry string
}

func SetupConfig() *DogoConfig {

	viper.SetConfigName("dogo")        // name of config file (without extension)
	viper.SetConfigType("yaml")        // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("/etc/dogo/")  // path to look for the config file in
	viper.AddConfigPath("$HOME/.dogo") // call multiple times to add many search paths
	viper.AddConfigPath(".")           // optionally look for config in the working directory
	err := viper.ReadInConfig()        // Find and read the config file
	if err != nil {                    // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	rdp := &RDP{
		DriveName: viper.GetString("dogo.guacd.rdp.drive-name"),
		DrivePath: viper.GetString("dogo.guacd.rdp.drive-path"),
	}

	ssh := &SSH{
		FontName:    viper.GetString("dogo.guacd.ssh.font-name"),
		FontSize:    viper.GetString("dogo.guacd.ssh.font-size"),
		ColorScheme: viper.GetString("dogo.guacd.ssh.color-scheme"),
	}

	vnc := &VNC{
		Autoretry: viper.GetString("dogo.guacd.vnc.autoretry"),
	}

	var guacd = &Guacd{
		Addr: viper.GetString("dogo.guacd.addr"),
		RDP:  rdp,
		SSH:  ssh,
		VNC:  vnc,
	}

	var config = &DogoConfig{
		Dsn:   viper.GetString("dogo.dsn"),
		Addr:  viper.GetStringSlice("dogo.addr"),
		Guacd: guacd,
	}

	return config
}
