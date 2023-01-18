package json_fetchers

import (
	"log"
	"testing"
)

func TestReadJsonFromText(t *testing.T) {
	data := `GET$$$/_msearch?rest_total_hits_as_int=true$$${"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu adidas 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu adidas nemeziz"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"futsal adidas 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"futsal nemeziz 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"futsal adidas nemeziz"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu nemeziz 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu futsal 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu futsal nemeziz"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu futsal adidas"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"adidas nemeziz 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}`
	ReadJsonFromText(data)
	t.Fail()
}

func TestReadSearchQueryFromMSearch(t *testing.T) {
	k := `POST$$$_msearch$$${"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu adidas 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"shop_product"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu adidas nemeziz"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"futsal adidas 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"futsal nemeziz 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"futsal adidas nemeziz"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu nemeziz 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu futsal 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu futsal nemeziz"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"sepatu futsal adidas"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}{"index":"product_ace"}{"query":{"bool":{"filter":[{"bool":{"must_not":{"exists":{"field":"parent_id"}}}},{"query_string":{"default_operator":"and","escape":true,"fields":["extended_name"],"fuzziness":"0","query":"adidas nemeziz 19 3"}}]}},"size":0,"terminate_after":5,"timeout":"100ms"}`
	res := ReadSearchQueryFromMSearch(k)

	log.Print(res)
}

func Test_getIndexFromURL(t *testing.T) {
	data := "/product_ace/_search"
	res := GetIndexFromURL(data)

	if res != "product_ace" {
		t.FailNow()
	}

	data = "/shop/_search?routing=adkjalsda"
	res = GetIndexFromURL(data)

	if res != "shop" {
		t.FailNow()
	}

}
