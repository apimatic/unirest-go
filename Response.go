package unirest

import(
    "io/ioutil"
    "net/http"
)

type Response struct {
    Code int
    RawBody []byte
    Body string
    Headers map[string][]string
}

func NewStringResponse(resp *http.Response) (*Response, error) {
    res, err :=	NewBinaryResponse(resp)
    if err != nil{
        return nil, err
    }
    
    res.Body = string(res.RawBody)
    return res, err
}

func NewBinaryResponse(resp *http.Response) (*Response, error) {
    //read response body
    res, err := ioutil.ReadAll(resp.Body)
    resp.Body.Close()
     
    //error reading response
    if err != nil {
        return nil, err
    }
    
    return makeResponse(resp.StatusCode, res, resp.Header), nil
}

func makeResponse(statusCode int, buff []byte, headers map[string][]string) *Response {
    //prepare a new request object
    response := new(Response)
    
    response.Code = statusCode
    response.RawBody = buff
    response.Headers = headers
    
    return response
}