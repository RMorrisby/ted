package help

type PageVariables struct {
	Date string
	Time string
	Port string
}

type ResultStruct struct {
	Name              string
	TestRunIdentifier string
	Category          string
	Status            string
	Timestamp         string
	Message           string
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
