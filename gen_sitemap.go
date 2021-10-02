package main

import (
	"bytes"
	"fmt"
	"time"

	"github.com/snabb/sitemap"
)

func genSiteMap(articles []*article, fqdn string) string {
	sm := sitemap.New()
	now := time.Now()
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s", fqdn), LastMod: &now})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/ja/", fqdn), LastMod: &now})
	sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/about/", fqdn), LastMod: &now})

	for _, a := range articles {
		sm.Add(&sitemap.URL{Loc: fmt.Sprintf("https://%s/%s", fqdn, link(a)), LastMod: &a.timestamp})
	}

	buff := &bytes.Buffer{}
	sm.WriteTo(buff)
	return buff.String()
}
