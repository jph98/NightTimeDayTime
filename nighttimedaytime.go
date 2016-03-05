//usr/bin/env go run $0 $@; exit

package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"strconv"
)

type Location struct {
	Ip string `json:"ip"`
    CountryCode string `json:"country_code"`
    CountryName string `json:"country_name"`
    RegionCode string `json:"region_code"`
    RegionName string `json:"region_name"`
    City string `json:"city"`
    ZipCode string `json:"zip_code"`
    TimeZone string `json:"time_zone"`
    Latitude float32 `json:"latitude"`
    Longitude float32 `json:"longitude"`
    MetroCode int `json:"metro_code"`
}

type Results struct {

	Sunrise string `json:"sunrise"`
	Sunset string `json:"sunset"`
	SolarNoon string `json:"solar_noon"`
	DayLength string `json:"day_length"`
	CivilTwilightBegin string `json:"civil_twilight_begin"`
	CivilTwilightEnd string `json:"civil_twilight_end"`
	NauticalTwilightBegin string `json:"nautical_twilight_begin"`
	NauticalTwilightEnd string `json:"nautical_twilight_end"`
	AstronomicalTwilightBegin string `json:"astronomical_twilight_begin"`
	AstronomicalTwilightEnd string `json:"astronomical_twilight_end"`
}

type SunResults struct {

	Result Results `json:"results"`
	Status string `json:"status"`
}

/* 
 * Get the current location using the geoip web service
 */
func getLocation() (Location) {

	geoip_url := "http://freegeoip.net/json/"

	res, err := http.Get(geoip_url)
	defer res.Body.Close()		

	if err != nil {
		panic(err)
	}		

	fmt.Printf("\nGetting your location...\n\n")
	var location Location
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&location)	

	if err != nil {
		panic(err)
	}
	
	return location
}

/*
 * Display the location details
 */
func displayLocationDetails(location Location) {

	fmt.Printf("\tCity: \t\t%s, %s\n", location.City, location.CountryName)

	latstr := strconv.FormatFloat(float64(location.Latitude), 'f', -1, 32) 
	longstr := strconv.FormatFloat(float64(location.Longitude), 'f', -1, 32) 
	fmt.Printf("\tLat/Long: \t%s,%s\n", latstr, longstr)	
}

/*
 * Get sun details using the sunrise-sunset web service
 */
func getSunDetails(date string, location Location) (SunResults) {

	// Example: http://api.sunrise-sunset.org/json?lat=36.7201600&lng=-4.4203400
	fmt.Printf("\nGetting times for...%s\n\n", date)	

	latstr := strconv.FormatFloat(float64(location.Latitude), 'f', -1, 32) 
	longstr := strconv.FormatFloat(float64(location.Longitude), 'f', -1, 32) 	

	sunrise_url := "http://api.sunrise-sunset.org/json"
	
	resource := sunrise_url + "?lat=" + latstr + "&lng=" + longstr + "&date=" + date 

	fmt.Printf("\tURL: \t\t%s\n", resource)
	res, err := http.Get(resource)
	defer res.Body.Close()

	var sunresults SunResults
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&sunresults)

	if err != nil {
		panic(err)
	}

	return sunresults	
}

/*
 * Display the sun details
 */
func displaySunDetails(sunresults SunResults) {

	fmt.Printf("\tStatus: \t%s\n", sunresults.Status)
	fmt.Printf("\tDay length: \t%s\n", sunresults.Result.DayLength)
	fmt.Printf("\tSunrise: \t%s\n", sunresults.Result.Sunrise)
	fmt.Printf("\tSunset: \t%s\n\n", sunresults.Result.Sunset)
}

func main() {

	fmt.Printf("\n- Nighttime/Daytime -\n")

	date := "today"

	// Can we define these as part of the function?	
	var location Location
	location = getLocation()
	displayLocationDetails(location)

	// Move to a function (sunrise, sunset)
	var sunresults SunResults
	sunresults = getSunDetails(date, location)
	displaySunDetails(sunresults)	
}