package dm

import (
	"net/url"
	"strconv"
	"strings"
	"fmt"
)

type Attribute struct {
	Name      string `json:"name, omitempty"`
	Region	  string `json: region, omitmepty"`
	Extension struct {
		Data    map[string]interface{} `json:"data, omitempty"`
		Version string                 `json:"version, omitempty"`
		Type    string                 `json:"type, omitempty"`
		Schema  struct {
			Href string `json:"href, omitempty"`
		} `json:"schema, omitempty"`
	} `json:"extension, omitempty"`
}

type Content struct {
	Relationships 	Relationships 	`json:"relationships, omitempty"`
	Attributes 	  	Attribute 		`json:"attributes, omitempty"`
	Type       		string    		`json:"type, omitempty"`
	Id         		string    		`json:"id, omitempty"`
	Links      		Link      		`json:"links, omitempty"`
}

type DataDetails struct {
	Data    	[]Content 	`json:"data, omitempty"`
	JsonApi 	JsonAPI   	`json:"jsonapi, omitempty"`
	Links   	Link      	`json:"links, omitempty"`
}

type ItemDetails struct {
	Data    	Content 	`json:"data, omitempty"`
	JsonApi 	JsonAPI   	`json:"jsonapi, omitempty"`
	Links   	Link      	`json:"links, omitempty"`
	Included 	[]Content 	`json:"included, omitempty"`
}

type Hub struct {
	Links 		Link 		`json:"links, omitempty"`
	Data 		[]Content 	`json:"data, omitempty"`
}

type JsonAPI struct {
	Version 	string 		`json:"version, omitempty"`
}

type Link struct {
	Self struct {
		Href 	string `json:"href, omitempty"`
	} `json:"self, omitempty"`
	First struct {
		Href string `json:"href, omitempty"`
	} `json:"first, omitempty"`
	Prev struct {
		Href string `json:"href, omitempty"`
	} `json:"prev, omitempty"`
	Next struct {
		Href string `json:"href, omitempty"`
	} `json:"next, omitempty"`
	Related struct {
		Href string `json:"href, omitempty"`
	} `json:"related, omitempty"`
}

type Project struct {
	Links Link `json:"links, omitempty"`
}

type Relationships struct {
	Projects 	Project 	`json:"projects, omitempty"`
	Hub 		struct {
		Links 		Link 		`json:"links, omitempty"`
		Data 		struct {
			ID 		string 	`json:"id, omitempty"`
			Type 	string 	`json:"type, omitempty"`
		} 	`json:"data, omitempty"`
	} 		`json:"hub, omitempty"`
	RootFolder 	RootFolder 	`json:"rootfolder, omitempty"`
	TopFolders 	TopFolders 	`json:"topfolders, omitempty"`
}

type RootFolder struct {
	Meta struct {
		Links Link `json:"links, omitempty"`
	} `json:"meta, omitempty"`
	Data struct {
		ID 		string 	`json:"id, omitempty"`
		Type 	string 	`json:"type, omitempty"`
	} 	`json:"data, omitempty"`
}

type TopFolders struct {
	Links 		Link 	`	json:"links, omitempty"`
}

//Query model
type Query struct {
	includePathInProject		bool
	filterType					[]string
	filterId					[]string
	filterExtensionType			[]string
	filterVersionNumber			[]int
	pageNumber					int
	pageLimit					int
}

//NewQuery initilazies a new query
func NewQuery() *Query {
	return &Query{
		includePathInProject:		false,
		filterType:					[]string{},
		filterId:					[]string{},
		filterExtensionType:		[]string{},
		filterVersionNumber:		[]int{},
		pageNumber:					0,
		pageLimit:					0,
	}
}

//Include query
func (q *Query) Include(include bool) *Query {
	q.includePathInProject = include
	return q
}

//FilterType query
func (q *Query) FilterType(ft string) *Query {
	q.filterType = append(q.filterType, ft)
	return q
}

//FilterId query
func (q *Query) FilterId(fi string) *Query {
	q.filterId = append(q.filterId, fi)
	return q
}

//FilterExtensionType query
func (q *Query) FilterExtensionType(fet string) *Query {
	q.filterExtensionType = append(q.filterExtensionType, fet)
	return q
}

//FilterVersionNumber equality query
func (q *Query) FilterVersionNumber(fvn []int) *Query {
	q.filterVersionNumber = fvn
	return q
}

//PageNumber equality query
func (q *Query) PageNumber(pn int) *Query {
	q.pageNumber = pn
	return q
}

//PageLimit equality query
func (q *Query) PageLimit(pl int) *Query {
	q.pageLimit = pl
	return q
}

// Values constructs url.Values
func (q *Query) Values() url.Values {
	params := url.Values{}

	if q.includePathInProject != false {
		params.Set("includePathInProject", strconv.FormatBool(q.includePathInProject))
	}

	if len(q.filterType) != 0 {
		params.Set("filter[type]", strings.Join(q.filterType, ","))
	}

	if len(q.filterId) != 0 {
		params.Set("filter[id]", strings.Join(q.filterId, ","))
	}

	if len(q.filterExtensionType) != 0 {
		params.Set("filter[extension.type]", strings.Join(q.filterExtensionType, ","))
	}

	if len(q.filterVersionNumber) != 0 {

		var stringArray []string
		for i := range q.filterVersionNumber {
			stringVal := strconv.Itoa(q.filterVersionNumber[i])
			stringArray := append(stringArray, stringVal)
			fmt.Println(stringArray)
		}

		params.Set("filter[versionNumber]", strings.Join(stringArray, "+"))
	}

	if q.pageNumber != 0 {
		params.Set("page[number]", strconv.Itoa(q.pageNumber))
	}

	if q.pageLimit != 0 {
		params.Set("page[limit]", strconv.Itoa(q.pageLimit))
	}

	return params
}

func (q *Query) String() string {
	return q.Values().Encode()
}