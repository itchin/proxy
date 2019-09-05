package parser

import (
    "github.com/itchin/proxy/client/config"
)

// 客户端所注册的域名
func getDomains() []string {
    domains := make([]string, 0)
    for key, _ := range config.DOMAINS {
        domains = append(domains, key)
    }
    return domains
}