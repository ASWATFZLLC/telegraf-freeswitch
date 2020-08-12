package main

import (
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/rif/telegraf-freeswitch/utils"
	"log"
	"net/http"
	"os"
	"strconv"
)

var (
	host          = flag.String("host", "localhost", "freeswitch host address")
	port          = flag.Int("port", 8021, "freeswitch port")
	pass          = flag.String("pass", "ClueCon", "freeswitch password")
	serve         = flag.Bool("serve", false, "run as a server")
	listenAddress = flag.String("listen_address", "127.0.0.1", "listen on address")
	listenPort    = flag.Int("listen_port", 9191, "listen on port")
)

func handler(w http.ResponseWriter, route string) {
}

func main() {
	flag.Parse()
	fetcher, err := utils.NewFetcher(*host, *port, *pass)
	if err != nil {
		fmt.Print("error connecting to fs: ", err)
	}
	defer fetcher.Close()
	if !*serve {
		if err := fetcher.GetData(); err != nil {
			fmt.Print(err.Error())
		}
		fmt.Print(fetcher.FormatOutput(utils.InfluxFormat))
		os.Exit(0)
	}
	r := gin.Default()
	r.GET("/status", func(c *gin.Context) {
		fetcher, err = utils.NewFetcher(*host, *port, *pass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fetcher.Close()
		if err := fetcher.GetData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		status, _, _ := fetcher.FormatOutput(utils.JSONFormat)
		c.String(http.StatusOK, status)
	})

	r.GET("/profiles", func(c *gin.Context) {
		fetcher, err = utils.NewFetcher(*host, *port, *pass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fetcher.Close()
		if err := fetcher.GetData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, fetcher.SofiaProfiles)
	})

	r.GET("/gateways", func(c *gin.Context) {
		fetcher, err = utils.NewFetcher(*host, *port, *pass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fetcher.Close()
		if err := fetcher.GetData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, fetcher.SofiaGateways)
	})

	r.GET("/gateways/check", func(c *gin.Context) {
		fetcher, err = utils.NewFetcher(*host, *port, *pass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fetcher.Close()
		if err := fetcher.GetData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		globalStatus := true
		for _, sofiaGateway := range fetcher.SofiaGateways {
			if sofiaGateway.Status == "0" {
				globalStatus = false
			}
		}
		if globalStatus == true {
			c.JSON(http.StatusOK, fetcher.SofiaGateways)
		} else {
			c.JSON(http.StatusInternalServerError, fetcher.SofiaGateways)
		}
	})
	r.GET("/gateways/check-latency", func(c *gin.Context) {
		fetcher, err = utils.NewFetcher(*host, *port, *pass)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		defer fetcher.Close()
		if err := fetcher.GetData(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		globalStatus := true
		for _, sofiaGateway := range fetcher.SofiaGateways {
			if sofiaGateway.Status == "0" {
				globalStatus = false
			}
			ping, _ := strconv.ParseFloat(sofiaGateway.Ping, 8)
			if ping > 100 {
				globalStatus = false
			}
		}
		if globalStatus == true {
			c.JSON(http.StatusOK, fetcher.SofiaGateways)
		} else {
			c.JSON(http.StatusInternalServerError, fetcher.SofiaGateways)
		}
	})
	listen := fmt.Sprintf("%s:%d", *listenAddress, *listenPort)
	fmt.Printf("Listening on %s...", listen)
	log.Fatal(r.Run(listen))
}
