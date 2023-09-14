package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/go-co-op/gocron"
	"github.com/urfave/cli/v2"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var lastCheckIPv6 string
var api *cloudflare.API
var (
	APIToken string
	ZoneId   string
) // 在这里声明为包级别变量

func init() {
	ZoneId = os.Getenv("ZONE_ID")
	APIToken = os.Getenv("API_TOKEN")

	if ZoneId == "" || APIToken == "" {
		log.Fatal("ZONE_ID or API_TOKEN not set in environment variables")
	}

	var err error
	api, err = cloudflare.NewWithAPIToken(APIToken)
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	app := &cli.App{
		Action: func(c *cli.Context) error {
			s := gocron.NewScheduler(time.UTC)

			from := os.Getenv("FROM_DOMAIN")
			to := os.Getenv("TO_DOMAIN")
			cycleStr := os.Getenv("CYCLE_MINUTES")
			cycle, err := strconv.Atoi(cycleStr)
			if err != nil || cycle <= 0 {
				cycle = 1
			}

			log.Printf("Set detection cycle: %d minutes", cycle)

			_, err = s.Every(cycle).Minutes().Do(func() {
				UpdateRecord(ZoneId, from, to)
			})
			if err != nil {
				log.Println(err)
			}
			s.StartBlocking()
			return nil
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func UpdateRecord(zone string, name string, to string) {
	ctx := context.Background()
	ip := GetIP(name)
	if ip != "" {
		log.Printf("Current IP: %s", ip)
		if ip == lastCheckIPv6 {
			log.Println("The IPv6 has not changed, waiting for the next check.")
			return
		}
		lastCheckIPv6 = ip

	} else {
		return
	}
	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zone), cloudflare.ListDNSRecordsParams{})
	if err != nil {
		log.Println(err)
		return
	}
	for _, r := range records {
		if r.Name == to && r.Type == "AAAA" {
			UpdateByID(ctx, zone, r.ID, to, ip)
		}
	}
}

func UpdateByID(ctx context.Context, zone string, id string, name string, content string) {
	err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zone), cloudflare.UpdateDNSRecordParams{
		ID:      id,
		Name:    name,
		Content: content,
		TTL:     60,
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Printf("Update record [%s] -> [%s] successfully!", name, content)
}

func GetIP(to string) string {
	ips, err := net.LookupIP(to)
	if err != nil {
		panic(err)
	}
	if len(ips) == 0 {
		fmt.Printf("no record")
		return ""
	}
	for _, ip := range ips {
		if ip.To4() == nil {
			log.Println(ip.String())
			return ip.String()
		}
	}
	return ""
}
