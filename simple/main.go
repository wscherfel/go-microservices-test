package main

import (
  "errors"
  "github.com/go-kit/kit/endpoint"
  "golang.org/x/net/context"
  "net/http"
  "encoding/json"
  transport "github.com/go-kit/kit/transport/http"
  "log"
)

type StringService interface {
  Save(key string, value string) error // save value of key in DB
  Load(key string) (string, error) // load value of key from DB
}


type StringServiceImplementation struct {
  values map[string]string // map of values indexed by keys
}

func (s *StringServiceImplementation) Save(key, value string) error {
  // if key exists return error
  if _, ok := s.values[key]; ok == true {
    return errors.New("Key already exists")
  }

  // key does not yet exist, save it
  s.values[key] = value
  return nil
}


func (s *StringServiceImplementation) Load(key string) (string, error) {
  // look for key
  value, ok := s.values[key]
  if ok == true {
    return value, nil
  }

  // key does not exist, return error
  return "", errors.New("Key not found")
}

type saveRequest struct {
  Key string `json:"key"`
  Value string `json:"value"`
}

type saveResponse struct {
  Err string `json:"err,omitempty"`
}

type loadRequest struct {
  Key string `json:"key"`
}

type loadResponse struct {
  Value string `json:"value"`
  Err string `json:"err,omitempty"`
}

func makeSaveEndpoint(svc StringService) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(saveRequest)
    err := svc.Save(req.Key, req.Value)
    if err != nil {
      return saveResponse{err.Error()}, nil
    }
    return saveResponse{""}, nil
  }
}

func makeLoadEndpoint(svc StringService) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(loadRequest)
    val, err := svc.Load(req.Key)
    if err != nil {
      return loadResponse{val, err.Error()}, nil
    }
    return loadResponse{val, ""}, nil
  }
}

func decodeSaveRequest(_ context.Context, r *http.Request) (interface{}, error) {
  var request saveRequest
  if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
    return nil, err
  }
  return request, nil
}

func decodeLoadRequest(_ context.Context, r *http.Request) (interface{}, error) {
  var request loadRequest
  if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
    return nil, err
  }
  return request, nil
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
  return json.NewEncoder(w).Encode(response)
}

func main() {
  ctx := context.Background()

  serv := &StringServiceImplementation{values:make(map[string]string)}

  servSaveHandler := transport.NewServer(
    ctx,
    makeSaveEndpoint(serv),
    decodeSaveRequest,
    encodeResponse,
  )

  servLoadHandler := transport.NewServer(
    ctx,
    makeLoadEndpoint(serv),
    decodeLoadRequest,
    encodeResponse,
  )

  http.Handle("/save", servSaveHandler)
  http.Handle("/load", servLoadHandler)

  log.Fatal(http.ListenAndServe(":8080", nil))
}
