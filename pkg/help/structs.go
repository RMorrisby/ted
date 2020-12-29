package help

type PageVariables struct {
	Date         string
	Time         string
	Port         string
	SuccessCount int
	FailCount    int
}

type ResultStruct struct {
	Name              string
	TestRunIdentifier string
	Category          string
	Status            string
	Timestamp         string
	Message           string
}

func NewResultStruct(csvLine []string) *ResultStruct {
	r := new(ResultStruct)
	r.Name = csvLine[0]
	r.TestRunIdentifier = csvLine[1]
	r.Category = csvLine[2]
	r.Status = csvLine[3]
	r.Timestamp = csvLine[4]
	r.Message = csvLine[5]
	return r
}

func (r ResultStruct) ToA() []string {
	resultArray := []string{
		r.Name,
		r.TestRunIdentifier,
		r.Category,
		r.Status,
		r.Timestamp,
		r.Message,
	}
	return resultArray
}

func ResultHeader() []string {
	header := []string{
		"Name",
		"TestRunIdentifier",
		"Category",
		"Status",
		"Timestamp",
		"Message",
	}
	return header
}
