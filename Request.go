package unirest

import (
    "errors"
    "bytes"
    "strconv"
    "time"
    "encoding/json"
    "encoding/base64"
    "net/http"
    "net/url"
    "reflect"
)

type Request struct {
    httpClient *http.Client
    connectTimeout int
    httpMethod HttpMethod  			//HTTP method for the outgoing request
    url string 						//Url for the outgoing request
    headers map[string]interface{}  //Headers for the outgoing request
    body interface{} 				//Parameters for raw body type request
    username string					//Basic auth password
    password string					//Basic auth password
}

type jsonable interface {
    ToJson() string
}

func NewRequest(method HttpMethod, url string,
     headers map[string]interface{}, parameters interface{},
     username string, password string) *Request {
     
     request := makeRequest(method, url, headers, username, password)
     request.body = parameters
     return request;
}
     
func makeRequest(method HttpMethod, url string,
     headers map[string]interface{},
     username string, password string) *Request {
     
     //prepare a new request object
     request := new(Request)
     
     //prepare the transport layer
     tr := &http.Transport{DisableKeepAlives: false, MaxIdleConnsPerHost: 2}
     request.httpClient = &http.Client{Transport: tr}
     request.connectTimeout = -1
     
     //perpare the request parameters
     request.httpMethod = method
     request.url = url
     request.headers = headers
     request.username = username
     request.password = password
     
     return request;
}
     
func (me *Request) PrepareRequest() (*http.Request, error) {
    var req *http.Request
    var err error    
    var method = me.httpMethod.ToString()
    
    //encode body and parameters to the request
    if(me.body != nil) {
        req, err = me.encodeBody(method)
    } else {
        req, err = http.NewRequest(method, me.url, nil)
    }
    
    //encode headers and basic auth fields
    req = me.encodeHeaders(req)
    
    //set timeout values
    me.httpClient.Transport.(*http.Transport).TLSHandshakeTimeout += 2 * time.Second
    me.httpClient.Transport.(*http.Transport).ResponseHeaderTimeout = 10 * time.Second
    
    return req, err
}
 
func (me *Request) encodeHeaders(req *http.Request) (*http.Request) {
    if(me.headers != nil) {
        for key, value := range me.headers {
            req.Header.Add(key, toString(value))
        }
    }
    if(len(me.username) > 1 || len(me.password) > 1) {
        authToken := base64.StdEncoding.EncodeToString([]byte(me.username + ":" + me.password))
        req.Header.Add("Authorization", "Basic " + authToken) 
    }
    return req
}

func (me *Request) encodeBody(method string) (*http.Request, error) {
    var req *http.Request
    var err error
        
    //given body is a param collection
    if params, ok := me.body.(map[string]interface{}); ok {
        paramValues := url.Values{} 
        for key, val := range params {
            paramValues.Add(key, toString(val))
        }
        req, err = http.NewRequest(method, me.url, nil)
        req.Form = paramValues
    } else { //given a raw body object
        bodyBytes, err := json.Marshal(me.body)
        if err != nil {
            return nil, errors.New("Invalid JSON in the query")
        }
        reader := bytes.NewReader(bodyBytes)
        req, err = http.NewRequest(method, me.url, reader)
        req.Header.Add("Content-Length", strconv.Itoa(len(string(bodyBytes))))
        req.Header.Add("Content-Type", "application/json; charset=utf-8")	
    } 	
    
    return req, err
}

func toString(val interface{}) string {
    value := reflect.ValueOf(val)
    valueType := value.Type().String()
    switch valueType {
        case "bool":
            return strconv.FormatBool(value.Bool())
        case "int", "int8", "int32", "int64",
             "uint", "uint8", "uint32", "uint64":
            return strconv.FormatInt(value.Int(), 10)
        case "float32":
            return strconv.FormatFloat(value.Float(), 'f', -1, 32)
        case "float64":
            return strconv.FormatFloat(value.Float(), 'f', -1, 64)
        case "string":
            return value.String()
        case "*time.Time":
            return value.Elem().Interface().(time.Time).String()
        case "time.Time":
            return value.Interface().(time.Time).String()
        default:
            jsonValue, _ := json.Marshal(val)
            return string(jsonValue[:])
    }    
}