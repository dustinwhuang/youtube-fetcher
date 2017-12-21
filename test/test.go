package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	fmt.Println("Running fetcher...")

	client := &http.Client{}

	req, err := http.NewRequest("GET", "https://www.googleapis.com/youtube/v3/videos", nil)
	if err != nil {
		panic(err)
	}

	q := req.URL.Query()
	q.Add("part", "snippet, contentDetails, statistics, topicDetails")
	q.Add("id", "86BrPe-y1fk,u3MJOFT2wZA,Zu2HHYJuta0,JnQrrkLeUW4,8DBx7UGdtSc,sYb7q3urPZ8,kIlajjTh8nU,bi-4Kuiya-Q,dM1jHyLZuk4,45WJpr9XZHk,gDLC50s82WU,Dc_zNmkZcWE,_H3yaggNl9s,E1is3fiTnfY,cG62F5FhpIQ,lQN4dcx6QwI,ft70a4VSRrM,t0hxtT98ZyQ,u2xGb6y3rBE,9-PkutG_Q2Q,gIcqpzaQRw0,EeEuScqMunM,vBKEDTC7yIQ,0MqP_4McZW8,NMQfje7IMfQ,ny2MkNiILj0,WVsu6AZNBHo,nIkWCCmYERo,pQAe8fLaZrw,dup-3ZSBJUs,qUt6NAlt3Dw,5sMyOmH_dUg,6MNCv-2MA_4,q7g68yH12PE,_SahcvK60FY,_Lky6lGAyA8,RX2Ttx4boRg,7fejr9YVY7Y,OSPVE5yyO8E,2reUbX39bgc,yot0ZzyIRxk,n5D9wWqjEXg,X0fflw9okwU,N3FVriLwSO0,3b3TEKHZNfA,nKzeaznhNyA,F3vuBkCQSNo,F8vDzoCD1OM,s-PkM7Jjv0w,goGeaa-fEpY,")
	q.Add("maxResults", "50")
	q.Add("key", "AIzaSyAhsPBoJk-7lzrW7E5fsE5IzEntSuzgEdc")
	q.Add("pageToken", "")

	req.URL.RawQuery = q.Encode()
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}

	respData, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(respData))

	resp.Body.Close()

}
