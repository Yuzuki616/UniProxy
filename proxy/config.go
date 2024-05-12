package proxy

import (
	"UniProxy/common/file"
	"UniProxy/geo"
	"UniProxy/v2b"
	"encoding/base64"
	"errors"
	"fmt"
	C "github.com/sagernet/sing-box/constant"
	"github.com/sagernet/sing-box/option"
	"net/netip"
	"net/url"
	"os"
	"path"
	"strconv"
)

func GetSingBoxConfig(uuid string, server *v2b.ServerInfo) (option.Options, error) {
	in := option.Inbound{}
	if TunMode {
		in.Type = "tun"
		in.TunOptions = option.TunInboundOptions{
			Inet4Address: option.Listable[netip.Prefix]{
				netip.MustParsePrefix("172.19.0.1/30"),
			},
			Inet6Address: option.Listable[netip.Prefix]{
				netip.MustParsePrefix("fdfe:dcba:9876::1/126"),
			},
			MTU:         9000,
			AutoRoute:   true,
			StrictRoute: true,
			Inet4RouteAddress: option.Listable[netip.Prefix]{
				netip.MustParsePrefix("0.0.0.0/1"),
				netip.MustParsePrefix("128.0.0.0/1"),
			},
			Inet6RouteAddress: option.Listable[netip.Prefix]{
				netip.MustParsePrefix("::/1"),
				netip.MustParsePrefix("8000::/1"),
			},
			Stack: "gvisor",
		}
	} else {
		in.Type = "mixed"
		addr, _ := netip.ParseAddr("127.0.0.1")
		in.MixedOptions = option.HTTPMixedInboundOptions{
			ListenOptions: option.ListenOptions{
				Listen:     (*option.ListenAddress)(&addr),
				ListenPort: uint16(InPort),
			},
			SetSystemProxy: SystemProxy,
		}
	}
	so := option.ServerOptions{
		Server:     server.Host,
		ServerPort: uint16(server.Port),
	}
	var out option.Outbound
	switch server.Type {
	case "vmess", "vless":
		transport := &option.V2RayTransportOptions{
			Type: server.Network,
		}
		switch transport.Type {
		case "tcp":
			transport.Type = ""
		case "http":
		case "ws":
			var u *url.URL
			u, err := url.Parse(server.NetworkSettings.Path)
			if err != nil {
				return option.Options{}, err
			}
			ed, _ := strconv.Atoi(u.Query().Get("ed"))
			transport.WebsocketOptions.EarlyDataHeaderName = "Sec-WebSocket-Protocol"
			transport.WebsocketOptions.MaxEarlyData = uint32(ed)
			transport.WebsocketOptions.Path = u.Path
		case "grpc":
			transport.GRPCOptions.ServiceName = server.ServerName
		}
		out = option.Outbound{
			Tag:  "p",
			Type: server.Type,
		}
		if server.Type == "vmess" {
			out.VMessOptions = option.VMessOutboundOptions{
				UUID:                uuid,
				Security:            "auto",
				AuthenticatedLength: true,
				Network:             "tcp",
				ServerOptions:       so,
				Transport:           transport,
			}
			if server.Tls == 1 {
				out.VMessOptions.TLS = &option.OutboundTLSOptions{
					Enabled:    true,
					ServerName: server.ServerName,
					Insecure:   server.TlsSettings.AllowInsecure != "0",
				}
			}
		} else {
			out.VLESSOptions = option.VLESSOutboundOptions{
				UUID:          uuid,
				ServerOptions: so,
				Flow:          server.Flow,
				Transport:     transport,
			}
			switch server.Tls {
			case 1:
				out.VLESSOptions.TLS = &option.OutboundTLSOptions{
					Enabled:    true,
					ServerName: server.ServerName,
					Insecure:   server.TlsSettings.AllowInsecure != "0",
				}
			case 2:
				out.VLESSOptions.TLS = &option.OutboundTLSOptions{
					Enabled:    true,
					ServerName: server.TlsSettings.RealityDest,
					Insecure:   true,
					UTLS: &option.OutboundUTLSOptions{
						Enabled:     true,
						Fingerprint: "chrome",
					},
					Reality: &option.OutboundRealityOptions{
						Enabled:   true,
						ShortID:   server.TlsSettings.ShortId,
						PublicKey: server.TlsSettings.PublicKey,
					},
				}
			}
		}
	case "shadowsocks":
		var keyLength int
		switch server.Cipher {
		case "2022-blake3-aes-128-gcm":
			keyLength = 16
		case "2022-blake3-aes-256-gcm":
			keyLength = 32
		}
		var pw string
		if keyLength != 0 {
			pw = base64.StdEncoding.EncodeToString([]byte(uuid[:keyLength]))
		} else {
			pw = uuid
		}
		out = option.Outbound{
			Type: "shadowsocks",
			Tag:  "p",
			ShadowsocksOptions: option.ShadowsocksOutboundOptions{
				ServerOptions: so,
				Password:      pw,
				Method:        server.Cipher,
			},
		}
	case "trojan":
		out = option.Outbound{
			Type: "trojan",
			Tag:  "p",
			TrojanOptions: option.TrojanOutboundOptions{
				ServerOptions: so,
				Password:      uuid,
			},
		}
		if server.Tls != 0 {
			out.TrojanOptions.TLS = &option.OutboundTLSOptions{
				Enabled:    true,
				ServerName: server.ServerName,
				Insecure:   server.TlsSettings.AllowInsecure != "0",
			}
		}
	case "hysteria":
		out = option.Outbound{
			Tag:  "p",
			Type: "hysteria",
			HysteriaOptions: option.HysteriaOutboundOptions{
				ServerOptions: option.ServerOptions{
					Server:     server.Host,
					ServerPort: uint16(server.Port),
				},
				UpMbps:     server.UpMbps,
				DownMbps:   server.DownMbps,
				Obfs:       server.ServerKey,
				AuthString: uuid,
			},
		}
		out.HysteriaOptions.TLS = &option.OutboundTLSOptions{
			Enabled:    true,
			Insecure:   server.AllowInsecure != 0,
			ServerName: server.ServerName,
		}
	default:
		return option.Options{}, errors.New("server type is unknown")
	}
	r, err := getRules(GlobalMode)
	if err != nil {
		return option.Options{}, fmt.Errorf("get rules error: %s", err)
	}
	return option.Options{
		Log: &option.LogOptions{
			//Output: path.Join(DataPath, "proxy.log"),
		},
		Inbounds: []option.Inbound{
			in,
		},
		Outbounds: []option.Outbound{
			out,
			{
				Tag:  "d",
				Type: "direct",
			},
		},
		Route: r,
	}, nil
}

func getRules(global bool) (*option.RouteOptions, error) {
	var r option.RouteOptions
	if !global {
		err := checkRes(DataPath)
		if err != nil {
			return nil, fmt.Errorf("check res err: %s", err)
		}
		r = option.RouteOptions{
			GeoIP: &option.GeoIPOptions{
				DownloadURL: ResUrl + "/geoip.db",
				Path:        path.Join(DataPath, "geoip.dat"),
			},
			Geosite: &option.GeositeOptions{
				DownloadURL: ResUrl + "/geosite.db",
				Path:        path.Join(DataPath, "geosite.dat"),
			},
		}
		r.Rules = []option.Rule{
			{
				Type: C.RuleTypeDefault,
				DefaultOptions: option.DefaultRule{
					GeoIP: option.Listable[string]{
						"cn", "private",
					},
					Geosite: option.Listable[string]{
						"cn",
					},
					Outbound: "d",
				},
			},
		}
		return &r, nil
	} else {
		r = option.RouteOptions{
			AutoDetectInterface: true,
		}
	}
	return &r, nil
}

func checkRes(p string) error {
	if !file.IsExist(path.Join(p, "geoip.dat")) {
		f, err := os.OpenFile(path.Join(p, "geoip.dat"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(geo.Ip)
		if err != nil {
			return err
		}
	}
	if !file.IsExist(path.Join(p, "geosite.dat")) {
		f, err := os.OpenFile(path.Join(p, "geosite.dat"), os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
		if err != nil {
			return err
		}
		defer f.Close()
		_, err = f.Write(geo.Site)
		if err != nil {
			return err
		}
	}
	return nil
}
