package main

import (
  "github.com/jinzhu/gorm"
  _ "github.com/jinzhu/gorm/dialects/sqlite"
  "github.com/kataras/go-errors"
  "golang.org/x/net/context"
  "net/http"

  "github.com/wscherfel/go-microservices-test/model" // model and DAO are here

  transport "github.com/go-kit/kit/transport/http"
  "encoding/json"
  "github.com/go-kit/kit/endpoint"
  "log"
  "fmt"
)

func panicIf(err error) {
  if err != nil {
    panic(err)
  }
}

/* this is moved in other package
// @dao
type StringModel struct {
  gorm.Model

  Key string
  Value string
}
*/

type StringService interface {
  Save(string, string) error
  Load(string, bool) (string, error)
}

type DAOInterface interface {
  Create(*model.StringModel)
  Read(*model.StringModel)[]model.StringModel
}

type StringServiceImplementation struct {
  dao DAOInterface

  otherServices []StringService
}

type saveRequest struct {
  Key string `json:"key"`
  Value string `json:"value"`
}

type saveResponse struct {
  Err string `json:"err,omitempty"`
}

func decodeSaveRequest(_ context.Context, r *http.Request) (interface{}, error) {
  var request saveRequest
  if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
    return nil, err
  }
  return request, nil
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

func (s *StringServiceImplementation) Save(key, value string) error {
  if values := s.dao.Read(&model.StringModel{Key:key}); len(values) != 0 {
    return errors.New("Key already exists")
  }

  // if key is found in other databases return error
  for _, serv := range s.otherServices {
    if _, err := serv.Load(key, false); err == nil {
      return errors.New("Key already exists")
    }
  }

  s.dao.Create(&model.StringModel{Key:key, Value:value})
  fmt.Println("saved key:", key, "value:", value)
  return nil
}

type loadRequest struct {
  Key string `json:"key"`
}

type loadResponse struct {
  Value string `json:"value"`
  Err string `json:"err,omitempty"`
}

func decodeLoadRequest(_ context.Context, r *http.Request) (interface{}, error) {
  var request loadRequest
  if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
    return nil, err
  }
  return request, nil
}

func makeLoadEndpoint(svc StringService) endpoint.Endpoint {
  return func(ctx context.Context, request interface{}) (interface{}, error) {
    req := request.(loadRequest)
    val, err := svc.Load(req.Key, true)
    if err != nil {
      return loadResponse{val, err.Error()}, nil
    }
    return loadResponse{val, ""}, nil
  }
}

func (s *StringServiceImplementation) Load(key string, askOthers bool) (string, error) {
  values := s.dao.Read(&model.StringModel{Key:key})

  // look for key in this db
  if len(values) > 0 {
    return values[0].Value, nil
  }

  // look for keys in other db
  if askOthers {
    for _, serv := range s.otherServices {
      if value, err := serv.Load(key, false); err == nil {
        return value, nil
      }
    }
  }

  return "", errors.New("Key not found")
}

func encodeResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
  return json.NewEncoder(w).Encode(response)
}

func main() {
  ctx := context.Background()
  dbA, err := gorm.Open("sqlite3", "dbA.sql")
  panicIf(err)
  dbB, err := gorm.Open("sqlite3", "dbB.sql")
  panicIf(err)
  dbC, err := gorm.Open("sqlite3", "dbC.sql")
  panicIf(err)
  dbD, err := gorm.Open("sqlite3", "dbD.sql")

  dbA.AutoMigrate(&model.StringModel{})
  dbB.AutoMigrate(&model.StringModel{})
  dbC.AutoMigrate(&model.StringModel{})
  dbD.AutoMigrate(&model.StringModel{})

  daoA := model.NewStringModelDAO(dbA)
  daoB := model.NewStringModelDAO(dbB)
  daoC := model.NewStringModelDAO(dbC)
  daoD := model.NewStringModelDAO(dbD)

  servA := &StringServiceImplementation{dao:daoA, otherServices:[]StringService{}}
  servB := &StringServiceImplementation{dao:daoB, otherServices:[]StringService{}}
  servC := &StringServiceImplementation{dao:daoC, otherServices:[]StringService{}}
  servD := &StringServiceImplementation{dao:daoD, otherServices:[]StringService{}}


  servA.otherServices = append(servA.otherServices, servB)
  servA.otherServices = append(servA.otherServices, servC)
  servA.otherServices = append(servA.otherServices, servD)

  servB.otherServices = append(servB.otherServices, servA)
  servB.otherServices = append(servB.otherServices, servC)
  servB.otherServices = append(servB.otherServices, servD)

  servC.otherServices = append(servC.otherServices, servA)
  servC.otherServices = append(servC.otherServices, servC)
  servC.otherServices = append(servC.otherServices, servD)

  servD.otherServices = append(servD.otherServices, servA)
  servD.otherServices = append(servD.otherServices, servB)
  servD.otherServices = append(servD.otherServices, servC)

  servAsaveHandler := transport.NewServer(
    ctx,
    makeSaveEndpoint(servA),
    decodeSaveRequest,
    encodeResponse,
  )

  servBsaveHandler := transport.NewServer(
    ctx,
    makeSaveEndpoint(servB),
    decodeSaveRequest,
    encodeResponse,
  )

  servCsaveHandler := transport.NewServer(
    ctx,
    makeSaveEndpoint(servC),
    decodeSaveRequest,
    encodeResponse,
  )

  servDsaveHandler := transport.NewServer(
    ctx,
    makeSaveEndpoint(servD),
    decodeSaveRequest,
    encodeResponse,
  )

  servAloadHandler := transport.NewServer(
    ctx,
    makeLoadEndpoint(servA),
    decodeLoadRequest,
    encodeResponse,
  )

  servBloadHandler := transport.NewServer(
    ctx,
    makeLoadEndpoint(servB),
    decodeLoadRequest,
    encodeResponse,
  )

  servCloadHandler := transport.NewServer(
    ctx,
    makeLoadEndpoint(servC),
    decodeLoadRequest,
    encodeResponse,
  )

  servDloadHandler := transport.NewServer(
    ctx,
    makeLoadEndpoint(servD),
    decodeLoadRequest,
    encodeResponse,
  )

  http.Handle("/asave", servAsaveHandler)
  http.Handle("/bsave", servBsaveHandler)
  http.Handle("/csave", servCsaveHandler)
  http.Handle("/dsave", servDsaveHandler)

  http.Handle("/aload", servAloadHandler)
  http.Handle("/bload", servBloadHandler)
  http.Handle("/cload", servCloadHandler)
  http.Handle("/dload", servDloadHandler)

  fmt.Println("server is starting")

  log.Fatal(http.ListenAndServe(":8080", nil))
}
