package appointment

import (
	"sort"
	"time"
)

type TimeSlot struct {
	Start time.Time
	End   time.Time
}

type Database interface {
	QueryEvents(startDate, endDate time.Time) []Event
}

func CalculateAvailableSlots(db Database, startDate time.Time) map[string][]TimeSlot {
	endDate := startDate.AddDate(0, 0, 7)
	events := db.QueryEvents(startDate, endDate)
	results := make(map[string][]TimeSlot)

	currentDate := startDate
	for range 7 {
		key := currentDate.Format("2006-01-02")
		results[key] = []TimeSlot{}
		currentDate = currentDate.AddDate(0, 0, 1)
	}

	openings, appointments := filteredEvents(events)

	for _, opening := range openings {
		openingStart := opening.StartsAt
		openingEnd := opening.EndsAt
		openingLabel := openingStart.Format("2006-01-02")

		var overlappingAppointments []Event
		for _, appointment := range appointments {
			appointmentStart := appointment.StartsAt
			appointmentEnd := appointment.EndsAt

			if appointmentEnd.After(openingStart) && (appointmentEnd.Before(openingEnd) || appointmentEnd.Equal(openingEnd)) ||
				(appointmentStart.After(openingStart) || appointmentStart.Equal(openingStart)) && appointmentStart.Before(openingEnd) ||
				(appointmentStart.Before(openingStart) || appointmentStart.Equal(openingStart)) && (appointmentEnd.Equal(openingEnd) || appointmentEnd.After(openingEnd)) {
				overlappingAppointments = append(overlappingAppointments, appointment)
			}
		}

		sort.Slice(overlappingAppointments, func(i, j int) bool {
			return overlappingAppointments[i].StartsAt.Unix() < overlappingAppointments[j].StartsAt.Unix()
		})

		slotStart := openingStart
		for _, oa := range overlappingAppointments {
			oaStart := oa.StartsAt
			oaEnd := oa.EndsAt

			if openingStart.Before(oaStart) {
				timeslot := TimeSlot{Start: slotStart, End: oa.StartsAt}
				results[openingLabel] = append(results[openingLabel], timeslot)
			}

			slotStart = oaEnd
		}

		if slotStart.Before(openingEnd) {
			timeslot := TimeSlot{Start: slotStart, End: openingEnd}
			results[openingLabel] = append(results[openingLabel], timeslot)
		}
	}

	return results
}

func filteredEvents(events []Event) (openings []Event, appointments []Event) {
	for _, e := range events {
		if e.Kind == "opening" {
			openings = append(openings, e)
		} else if e.Kind == "appointment" {
			appointments = append(appointments, e)
		}
	}

	return openings, appointments
}
