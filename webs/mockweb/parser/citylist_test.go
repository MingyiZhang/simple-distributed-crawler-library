package parser

import (
  "io/ioutil"
  "testing"
)

func TestParseCityList(t *testing.T) {
  contents, err := ioutil.ReadFile("citylist_test_data.html")
  if err != nil {
   panic(err)
  }
  result := ParseCityList(contents, "")

  const resultSize = 2
  expectedUrls := []string{
   "http://localhost:8080/mock/www.zhenai.com/zhenghun/aba",
   "http://localhost:8080/mock/www.zhenai.com/zhenghun/akesu",
  }

  if len(result.Requests) != resultSize {
   t.Errorf("result should have %d requests; but had %d", resultSize, len(result.Items))
  }
  for i, url := range expectedUrls {
   if result.Requests[i].Url != url {
     t.Errorf("expected url #%d: %s; but was %s", i, url, result.Requests[i].Url)
   }
  }
}
