package v1

type OData struct {
	Top   int  `form:"$top"`
	Skip  int  `form:"$skip"`
	Count bool `form:"$count"`
}
