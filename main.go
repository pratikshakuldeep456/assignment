package main

import (
	"fmt"
)

type User struct {
	ID           string
	Name         string
	Gender       string
	Age          int
	CurrentState string
	CurrentDist  string
}

type VaccinationCenter struct {
	State      string
	Dist       string
	CenterID   string
	Capacities map[int]int             //capacity of each day
	Appoinment map[int]map[string]bool //for booking appointment for each day
}

type Appoinment struct {
	UserID    string
	CenterID  string
	Day       int
	BookingID string
}

type VaccinationSystem struct {
	Users              map[string]*User
	VaccinationCenters map[string]*VaccinationCenter
	Bookings           map[string]*Appoinment
}

// method to create a new user
func (vs *VaccinationSystem) AddUser(id, name, gender string, age int, state, dist string) {
	if age < 18 || (gender != "Male " && gender != "Female") {
		fmt.Println("Invalid user details ")
		return
	}

	user := &User{
		ID:           id,
		Name:         name,
		Age:          age,
		CurrentState: state,
		CurrentDist:  dist,
	}
	vs.Users[id] = user

}

// to onboard a new centre
func (vs *VaccinationSystem) AddVaccinationCentre(state, dist, centerID string) {

	if _, exists := vs.VaccinationCenters[centerID]; exists {
		fmt.Println("center already exists")
		return
	}

	center := &VaccinationCenter{
		State:      state,
		Dist:       dist,
		CenterID:   centerID,
		Capacities: make(map[int]int),
		Appoinment: make(map[int]map[string]bool),
	}

	vs.VaccinationCenters[centerID] = center
}

//capacity addition to specicfic day

func (vs *VaccinationSystem) AddCapacity(centerID string, day, capacity int) {

	center := vs.VaccinationCenters[centerID]
	if center == nil || day < 1 || day > 7 {

		fmt.Println("invalid id or day")
		return
	}
	if _, exists := center.Capacities[day]; exists {
		fmt.Println("capacity already exisrts")
		return
	}

	center.Capacities[day] = capacity
	center.Appoinment[day] = make(map[string]bool)

}

//method to list of vaccination centre

func (vs *VaccinationSystem) ListVaccinationCenters(dist string) {

	for _, center := range vs.VaccinationCenters {
		if center.Dist == dist {
			fmt.Printf(center.CenterID)
			fmt.Printf(center.State)
			fmt.Printf(center.Dist)

			fmt.Println("details of capacity")

			for day, capacity := range center.Capacities {
				fmt.Println(day, capacity) //recheck
			}
		}
	}
}

//method for booking a appointment

func (vs *VaccinationSystem) BookVaccination(centerID string, day int, userID string) string {
	center := vs.VaccinationCenters[centerID]
	user := vs.Users[userID]
	if center == nil || user == nil || user.Age < 18 || center.Dist != user.CurrentDist {
		return "Invalid booking req "
	}

	if day < 1 || day > 7 {
		return "invalid day"
	}

	if center.Capacities[day] <= len(center.Appoinment[day]) {
		return "no slot avialble "
	}

	if _, alreadyBooked := center.Appoinment[day][userID]; alreadyBooked {
		return "already booked"
	}
	bookingID := generateBookingID(vs)

	appointment := &Appoinment{
		UserID:    userID,
		CenterID:  centerID,
		Day:       day,
		BookingID: bookingID,
	}
	center.Appoinment[day][userID] = true

	vs.Bookings[bookingID] = appointment
	return bookingID
}

// list all booking from perticular center
func (vs *VaccinationSystem) ListAllBooking(centerID string, day int) {
	center := vs.VaccinationCenters[centerID]
	if center == nil {
		fmt.Println("invlaid")
		return
	}

	fmt.Printf("booknig for %s on Day %d:\n", centerID, day)

	for bookingID, appointment := range vs.Bookings {
		if appointment.CenterID == centerID && appointment.Day == day {
			user := vs.Users[appointment.UserID]
			fmt.Printf("Booknig ID: %s/n", bookingID)
			fmt.Printf(user.ID)
			fmt.Printf(user.Name)
			fmt.Printf("user Age : %d/n ", user.Age)
		}

	}
}

// cancel booking

func (vs *VaccinationSystem) CancelBooking(centerID, bookingID, userId string) string {

	center := vs.VaccinationCenters[centerID]
	if center == nil {
		return "invlaid center id"
	}

	appointment, exists := vs.Bookings[bookingID]
	if !exists || appointment.CenterID != centerID || appointment.UserID != userId {
		return "INvalid booknig id"
	}
	day := appointment.Day
	delete(center.Appoinment[day], userId)
	delete(vs.Bookings, bookingID)

	return "succesfully removed"

}

func (vs *VaccinationSystem) SearchVaccinationcenter(day int, userid string) {
	user := vs.Users[userid]
	if user == nil || user.Age < 18 {
		fmt.Println("invalid")

		return
	}

	for _, center := range vs.VaccinationCenters {
		if center.Dist == user.CurrentDist && center.Capacities[day] > len(center.Appoinment[day]) {
			fmt.Printf(center.CenterID)
			fmt.Printf(center.State)
			fmt.Printf(center.Dist)

			//	fmt.Printf(center.Capacities[day] - len(center.Appoinment[day]))

		}
	}
}
func generateBookingID(vs *VaccinationSystem) string {
	return fmt.Sprintf("BK%d", len(vs.Bookings)+1)
}
func main() {

	vs := &VaccinationSystem{
		Users:              make(map[string]*User),
		VaccinationCenters: make(map[string]*VaccinationCenter),
		Bookings:           make(map[string]*Appoinment),
	}
	vs.AddUser("U1", "Harry", "Male", 30, "karnataka", "bangalore")
	vs.AddUser("U2", "Harsht", "Male", 30, "karnataka", "bangalore")
	vs.AddUser("U3", "Hera", "Female", 30, "karnataka", "bangalore")
	vs.AddUser("U4", "Harryanka", "Female", 30, "karnataka", "bangalore")

	vs.AddVaccinationCentre("karnataka", "bangalore", "vc1")

	vs.AddCapacity("vc1", 1, 5)

	fmt.Println("details fetch")

	fmt.Println(vs.BookVaccination("vc1", 1, "U1"))
	fmt.Println(vs.BookVaccination("vc1", 1, "U2"))

	vs.ListVaccinationCenters("bangalore")
	vs.ListAllBooking("vc1", 1)
	fmt.Println(vs.CancelBooking("vc1", "bk3", "U1"))
	vs.SearchVaccinationcenter(1, "U1")
}
