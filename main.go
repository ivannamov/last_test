package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	trainsLength    = 3
	trainTimeLayout = "15:04:05"
)

type Trains []Train

type Train struct {
	TrainID            int
	DepartureStationID int
	ArrivalStationID   int
	Price              float32
	ArrivalTime        time.Time
	DepartureTime      time.Time
}

type TrainJSON struct {
	TrainID            int     `json:"trainId"`
	DepartureStationID int     `json:"departureStationId"`
	ArrivalStationID   int     `json:"arrivalStationId"`
	Price              float32 `json:"price"`
	ArrivalTime        string  `json:"arrivalTime"`
	DepartureTime      string  `json:"departureTime"`
}

func main() {
	// input
	fmt.Println("Enter departure station:")
	departureStation, inputErr := ReadInput()

	fmt.Println("Enter arrival station:")
	arrivalStation, inputErr := ReadInput()

	fmt.Println("Enter criteria:")
	criteria, inputErr := ReadInput()
	if inputErr != nil {
		fmt.Println("reading input failed", inputErr)
	}

	result, err := FindTrains(departureStation, arrivalStation, criteria)
	if err != nil {
		fmt.Println(err)
	}

	//	print result
	fmt.Println("Result:", result)
}

func FindTrains(departureStation, arrivalStation, criteria string) (Trains, error) {
	// validate departure station if empty
	if departureStation == "" {
		return nil, errors.New("empty departure station")
	}
	// validate arrival station if empty
	if arrivalStation == "" {
		return nil, errors.New("empty arrival station")
	}
	// validate criteria
	if criteria != "price" && criteria != "arrival-time" && criteria != "departure-time" {
		return nil, errors.New("unsupported criteria")
	}

	departureStationId, err := strconv.Atoi(departureStation)
	if err != nil {
		return nil, errors.New("bad departure station input")
	}

	arrivalStationId, err := strconv.Atoi(arrivalStation)
	if err != nil {
		return nil, errors.New("bad arrival station input")
	}

	trains, jsonError := ReadTrainsJson()
	if jsonError != nil {
		fmt.Println("reading json failed", jsonError)
	}

	filteredTrains := FilterTrains(trains, departureStationId, arrivalStationId)

	sortedTrains := SortTrains(filteredTrains, criteria)

	lengthTrains := trainsLength
	if len(sortedTrains) < lengthTrains {
		lengthTrains = len(sortedTrains)
	}

	return sortedTrains[:lengthTrains], nil
}

func SortTrains(trains Trains, criteria string) Trains {
	sort.Slice(trains, func(i, j int) bool {
		switch criteria {
		case "price":
			return trains[i].Price < trains[j].Price
		case "departure-time":
			return trains[i].DepartureTime.Before(trains[j].DepartureTime)
		case "arrival-time":
			return trains[i].ArrivalTime.Before(trains[j].ArrivalTime)
		}
		return false
	})

	return trains
}

func FilterTrains(trains Trains, departureStationId int, arrivalStationId int) Trains {
	var res Trains
	for _, train := range trains {
		if train.DepartureStationID == departureStationId && train.ArrivalStationID == arrivalStationId {
			res = append(res, train)
		}
	}
	return res
}

func (train *Train) UnmarshalJSON(data []byte) error {
	var v TrainJSON

	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	departureTime, err := time.Parse(trainTimeLayout, v.DepartureTime)
	arrivalTime, err := time.Parse(trainTimeLayout, v.ArrivalTime)

	if err != nil {
		return err
	}

	train.TrainID = v.TrainID
	train.DepartureStationID = v.DepartureStationID
	train.ArrivalStationID = v.ArrivalStationID
	train.Price = v.Price
	train.DepartureTime = departureTime
	train.ArrivalTime = arrivalTime

	return nil
}

func ReadTrainsJson() (Trains, error) {
	// read json file
	content, err := ioutil.ReadFile("./data.json")
	if err != nil {
		return nil, err
	}

	// unmarshall data
	var data Trains
	err = json.Unmarshal(content, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func ReadInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)
	// ReadString will block until the delimiter is entered
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// remove the delimiter from the string
	return strings.TrimSuffix(strings.TrimSuffix(input, "\n"), "\r"), nil
}
