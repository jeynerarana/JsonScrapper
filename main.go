package main

import (
	"encoding/json"
    "fmt"
    "io/ioutil"
    "os"
    "sort"
    "math"
    "bufio"
    "strings"
//    "strconv"
    
)
//Let's defne our Log struct
type Logs struct{
	Log Log `json:"log"`
}
type Log struct{
	Creator   Creator `json:"creator"`
	Browser   Browser `json:"browser"`
	Pages     []Pages   `json:"pages"`
	Entries   []Entries `json:"entries"`
	Version   string   `json:"version"`

}
type Creator struct{
	Name string `json:"name"`
	Version string `json:"version"`

}
type Browser struct{
	Name string `json:"name"`
	Version string `json:"version"`

}
type Pages struct{
	FirstPage       FirstPage   `json:"0"`
	StartedDateTime string      `json:"startedDateTime"`
	Id              string      `json:"id"`
	Title           string      `json:"title"`
	PageTimings     PageTimings `json:"pageTimings"`


}
type Entries struct{
	Pageref         string   `json:"pageref"`
	StartedDateTime string   `json:"startedDateTime"`
	Request         Request  `json:"request"`
	Response        Response `json:"response"`
	//Cache           Cache    'json:"cache"'
	//Timing          Timing   'json:"timings"'
	Time            int      `json:"time"`
	SecurityState   string   `json:"_securityState"`

}
type Request struct{
	Url         string `json:"url"`
	HttpVersion string `json:"httpVersion"`
	HeadersSize  int   `json:"headersSize"`

}
type Response struct{
	HeadersSize     int      `json:"headersSize"`
}
type FirstPage struct{
	StartedDateTime string      `json:"startedDateTime"`
	Id              string      `json:"id"`
	Title           string      `json:"title"`
	PageTimings     PageTimings `json:"pageTimings"`
}
type PageTimings struct{
	OnContentLoad int `json:"onContentLoad"`
	OnLoad        int `json:"onLoad"`
}
//lets calculate the median number of requests
func mNumRequest(b [] int)float32{ //DONE
	//save result
	
	//sort.Float64s(l.Log.Entries) // sort the numbers
	//numOfRequests := l.Log.Entries
	//algo
	mNumber := len(b) / 2
	//odd
	if (len(b)%2 !=0) {
		return float32(b[int32(((len(b)-1)/2)+1)])
	} else {
		return float32((b[mNumber -1] + b[mNumber +1])/2)

	}
	
}
//lets find percentage of URLs using HTTPS
func urlUsingHttps(l Logs, total *int){//DONE
	for value:=0 ;value < len(l.Log.Entries) ; value++{
		if(l.Log.Entries[value].SecurityState =="secure"){
			*total++
		}
	}

}
//percentage of requests over HTTP/1.1, HTTP/2 and HTTP/3
func httpOnetoThree(l Logs,one *int,two *int,three *int,count *int){
	//DONE
	for value:=0 ;value < len(l.Log.Entries) ; value++{
		if (l.Log.Entries[value].Request.HttpVersion == "HTTP/1.1") {
			*one++
			*count++
		} else if (l.Log.Entries[value].Request.HttpVersion == "HTTP/2") {
			*two++
			*count++
		} else if(l.Log.Entries[value].Request.HttpVersion == "HTTP/3"){
			*three++
			*count++
		}
	}
}
//median, min and max page load time
func mmmLoadTimes(pageTimings []float64, min *int, max *int , median *float32){
	//DONE
	mNumber := len(pageTimings) / 2
	//Min, MAx
	for value:= 0; value < len(pageTimings); value++{
			if int(math.Abs((pageTimings[value]))) < *min {
				*min = int(math.Abs(pageTimings[value]))
			}
			if int(math.Abs(pageTimings[value])) > *max {
				*max = int(math.Abs(pageTimings[value]))
			}

	}
	//Median
	//odd
	if (len(pageTimings)%2 !=0) {
		*median = float32(math.Abs(pageTimings[int32(((len(pageTimings)-1)/2)+1)]))
	}
	//even 
		*median = float32(math.Abs(pageTimings[mNumber-1] + pageTimings[mNumber+1])/2)




}
//median, min and max number of bytes spent on Response HTTP headers
func resonseHeadersByte(l Logs, rsmin *int, rsmax *int, rqmin *int, rqmax *int){
	//TODO

	var a [] Entries = l.Log.Entries

	//lets check Request Headers
	for value:=0 ;value < len(l.Log.Entries) ; value++{
		if a[value].Request.HeadersSize < *rqmin {
			*rqmin += a[value].Request.HeadersSize
			//fmt.Println(a[value].Request.HeadersSize)
		}
		if a[value].Request.HeadersSize  > *rqmax {
			*rqmax = a[value].Request.HeadersSize
		}
	}
	
	//Now Response HTTP headers
	for value:=0 ;value < len(l.Log.Entries) ; value++{
		if a[value].Response.HeadersSize < *rsmin {
			*rsmin = a[value].Response.HeadersSize
			//fmt.Println(a[value].Response.HeadersSize)
		}
		if a[value].Response.HeadersSize  > *rsmax {
			*rsmax = a[value].Response.HeadersSize
		}
	}

}
//median total kilobytes transferred per har page 
func mTotalBytes(b [] int)float32{ //DONE
	mNumber := len(b) / 2
	//odd
	if (len(b)%2 !=0) {
		return float32(b[int32(((len(b)-1)/2)+1)])
	}
	//even 
		return float32((b[mNumber-1] + b[mNumber+1])/2)
		


}
//For onload
func countOnload(l Logs)int{ //DONE
	for value:=0 ;value < len(l.Log.Entries) ; value++{
		return l.Log.Pages[value].PageTimings.OnLoad
	}
	return 0


}



func main() {

	//array used to store size of file
	var fSize [] int
	// array used to store the number of requestsper page(har file)
	var requests [] int
	//array for ratio of http version
	var httpv1 int =0
	var httpv2 int =0
	var httpv3 int =0
	var overallHttp int =0
	//variable for secure headers
	var secureHeaders int = 0
	//variables fro Min, Media, Max Load times
	var minLoad    int =0
	var maxLoad    int =0
	var medianLoad float32 =0
	var pageTimings [] float64 
	//varibles for response Headers
	var responseMinBytes int =0
	var requestMinBytes  int =0
	var responseMaxBytes int =0
	var requestMaxBytes  int =0
	// ask the user for the file
	fmt.Println("Please type the director where the Har files are: ")
	reader := bufio.NewReader(os.Stdin)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)
	//fmt.Println(text)
	//Lets Iterate through the whole Directory
	items, _ := ioutil.ReadDir(text)
	fmt.Println("Loading...Please wait")
    for _, item := range items {


		jsonFile, err := os.Open(text+item.Name())
		//prints out an error if there is nothing 
		if err != nil {
	    		fmt.Println(err)
		}
		//closes the file
		defer jsonFile.Close()
		//Lets find the files size of current file
		fi, err := jsonFile.Stat()
		if err != nil {
  			// Could not obtain stat, handle error
			}
		//lets pass the size to our 
		fSize =append(fSize,int(fi.Size()))
		//Let's initialize our user array
		var logs Logs
		//read all bytes
		byteValue, _ := ioutil.ReadAll(jsonFile)

	    // we unmarshal our byteArray which contains our
	    // jsonFile's content into 'logs' which we defined above
	    json.Unmarshal(byteValue, &logs)

	    //fmt.Println("Browser: " + mNumber(logs))
	    //mNumRequest(logs)
	    //bytesSpendonHeaders(logs)

	    //Let's count number of requests per page
	    var rCount int = len(logs.Log.Entries)
	    //save the request to our array of requests
	    requests =append(requests,rCount)

	    //Lets count the HTTP versions
	    httpOnetoThree(logs,&httpv1,&httpv2,&httpv3,&overallHttp)
	    //lets count number of secure websites
	    urlUsingHttps(logs, &secureHeaders)
	    //lets save the Onload page timings to our designated array
	    pageTimings = append(pageTimings,float64(countOnload(logs)))

	    //Now lets save the headers
	    resonseHeadersByte(logs, &responseMinBytes, &responseMaxBytes, &requestMinBytes, &requestMaxBytes)


	}
	//sort the array
	sort.Ints(fSize)
	fmt.Println("\n---------Size of Har Files----------\n")
	fmt.Println("Median total kilobytes transferred per har page:", int(mTotalBytes(fSize)), "bytes")
	//sort requests array
	sort.Ints(requests)
	fmt.Println("Median number of requests per page", int(mNumRequest(requests))," Requests")
	//Lets see how many times the http versions showed up in Percent
	//turn into percent
	//array for ratio of http version
	pHttpv1:=(float32(httpv1)/float32(overallHttp)) * 100
	pHttpv2:=(float32(httpv2)/float32(overallHttp)) * 100
	pHttpv3:=(float32(httpv3)/float32(overallHttp)) * 100
	fmt.Println("\n---------HTTP/1.1 vs HTTP/2 vs HTTP/3----------\n")
	fmt.Println("HTTP/1:",pHttpv1,"Percent\n",
	"HTTP/2:",pHttpv2,"Percent\n","HTTP/3:",pHttpv3,
	"Percent\nTotal Requests Made:",overallHttp)
	//print the ratio of secure headers vs non secure
	fmt.Println("\n-----------HTTPS on URLs------------\n")
	psecureHeaders:=(float32(secureHeaders)/float32(overallHttp)) * 100
	fmt.Println("Secure Headers",psecureHeaders,"Percent")
	//lets fing min,max, median page timings
	fmt.Println("-----------Page Load Times------------")
	mmmLoadTimes(pageTimings,&minLoad,&maxLoad,&medianLoad)
	fmt.Println("Min Load Time:",minLoad,"ms\nMax Load Time:",maxLoad,"ms\nMedian Load Time:",medianLoad,"ms\n")
	//let's Move on to Headers
	fmt.Println("-----------Page Response Headers Bytes------------\n")
	fmt.Println("Minimum Number of Response bytes:",responseMinBytes)
	fmt.Println("Maximum Number of Response bytes:",responseMaxBytes)

	fmt.Println("\n-----------Page Request Headers Bytes------------\n")
	fmt.Println("Minimum Number of Request bytes:",requestMinBytes)
	fmt.Println("Maximum Number of Request bytes:",requestMaxBytes)
  
}
