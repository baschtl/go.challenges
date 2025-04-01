package appointment

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockDB struct {
	mock.Mock
}

func (m *MockDB) QueryEvents(startDate, endDate time.Time) []Event {
	args := m.Called(startDate, endDate)
	return args.Get(0).([]Event)
}

var mockEvents = []Event{
	// Openings - when doctor is available
	{ID: 1, Kind: "opening", StartsAt: parseTime("2025-03-30T09:00:00.000Z"), EndsAt: parseTime("2025-03-30T12:00:00.000Z")},
	{ID: 2, Kind: "opening", StartsAt: parseTime("2025-03-30T14:00:00.000Z"), EndsAt: parseTime("2025-03-30T18:00:00.000Z")},
	{ID: 3, Kind: "opening", StartsAt: parseTime("2025-03-31T10:00:00.000Z"), EndsAt: parseTime("2025-03-31T14:00:00.000Z")},
	{ID: 4, Kind: "opening", StartsAt: parseTime("2025-04-01T09:00:00.000Z"), EndsAt: parseTime("2025-04-01T12:00:00.000Z")},
	{ID: 5, Kind: "opening", StartsAt: parseTime("2025-04-01T13:00:00.000Z"), EndsAt: parseTime("2025-04-01T18:00:00.000Z")},
	{ID: 6, Kind: "opening", StartsAt: parseTime("2025-04-02T10:00:00.000Z"), EndsAt: parseTime("2025-04-02T16:00:00.000Z")},
	{ID: 7, Kind: "opening", StartsAt: parseTime("2025-04-03T09:00:00.000Z"), EndsAt: parseTime("2025-04-03T17:00:00.000Z")},
	{ID: 8, Kind: "opening", StartsAt: parseTime("2025-04-04T09:00:00.000Z"), EndsAt: parseTime("2025-04-04T16:00:00.000Z")},
	{ID: 9, Kind: "opening", StartsAt: parseTime("2025-04-05T10:00:00.000Z"), EndsAt: parseTime("2025-04-05T14:00:00.000Z")},

	// Appointments - when doctor is booked
	{ID: 101, Kind: "appointment", StartsAt: parseTime("2025-03-30T10:00:00.000Z"), EndsAt: parseTime("2025-03-30T10:30:00.000Z")},
	{ID: 102, Kind: "appointment", StartsAt: parseTime("2025-03-30T11:00:00.000Z"), EndsAt: parseTime("2025-03-30T11:30:00.000Z")},
	{ID: 103, Kind: "appointment", StartsAt: parseTime("2025-03-30T15:00:00.000Z"), EndsAt: parseTime("2025-03-30T16:00:00.000Z")},
	{ID: 104, Kind: "appointment", StartsAt: parseTime("2025-03-31T11:00:00.000Z"), EndsAt: parseTime("2025-03-31T12:00:00.000Z")},
	{ID: 105, Kind: "appointment", StartsAt: parseTime("2025-04-01T09:30:00.000Z"), EndsAt: parseTime("2025-04-01T10:30:00.000Z")},
	{ID: 106, Kind: "appointment", StartsAt: parseTime("2025-04-02T11:00:00.000Z"), EndsAt: parseTime("2025-04-02T11:30:00.000Z")},
	{ID: 107, Kind: "appointment", StartsAt: parseTime("2025-04-02T14:00:00.000Z"), EndsAt: parseTime("2025-04-02T15:00:00.000Z")},
	{ID: 108, Kind: "appointment", StartsAt: parseTime("2025-04-03T13:00:00.000Z"), EndsAt: parseTime("2025-04-03T14:00:00.000Z")},
	{ID: 109, Kind: "appointment", StartsAt: parseTime("2025-04-04T11:00:00.000Z"), EndsAt: parseTime("2025-04-04T12:00:00.000Z")},
}

func filterEvents(events []Event, startDate, endDate time.Time) []Event {
	var filtered []Event
	for _, event := range events {
		if (event.StartsAt.Equal(startDate) || event.StartsAt.After(startDate)) && event.StartsAt.Before(endDate) {
			filtered = append(filtered, event)
		}
	}
	return filtered
}

func parseTime(timeStr string) time.Time {
	t, _ := time.Parse(time.RFC3339, timeStr)

	return t
}

func formatTime(t time.Time) string {
	return t.Format(time.RFC3339)
}

// Helper for creating time slot with string representation
func makeTimeSlot(start, end string) TimeSlot {
	return TimeSlot{
		Start: parseTime(start),
		End:   parseTime(end),
	}
}

func TestDoctorAppointmentAvailability(t *testing.T) {
	// Setup mock database
	mockDB := new(MockDB)

	t.Run("should return available slots for the next 7 days", func(t *testing.T) {
		// Reset mock
		mockDB.ExpectedCalls = nil

		startDate := parseTime("2025-03-30T00:00:00.000Z")
		endDate := parseTime("2025-04-06T00:00:00.000Z")

		// Setup mock to return filtered events
		mockDB.On("QueryEvents", startDate, endDate).Return(filterEvents(mockEvents, startDate, endDate))

		// Call the function under test
		result := CalculateAvailableSlots(mockDB, startDate)

		// Verify that the result contains data for 7 days
		assert.Len(t, result, 7)

		// Check that the first date is the start date
		_, hasFirstDay := result["2025-03-30"]
		assert.True(t, hasFirstDay)

		// Check that the last date is 6 days after the start date
		_, hasLastDay := result["2025-04-05"]
		assert.True(t, hasLastDay)

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("should correctly calculate available slots for a specific day", func(t *testing.T) {
		// Reset mock
		mockDB.ExpectedCalls = nil

		startDate := parseTime("2025-03-30T00:00:00.000Z")
		endDate := parseTime("2025-04-06T00:00:00.000Z")

		// Setup mock
		mockDB.On("QueryEvents", startDate, endDate).Return(filterEvents(mockEvents, startDate, endDate))

		// Call the function under test
		result := CalculateAvailableSlots(mockDB, startDate)

		// On March 30, doctor has openings from 9-12 and 14-18
		// With appointments at 10-10:30, 11-11:30, and 15-16
		// So available slots should be: 9-10, 10:30-11, 11:30-12, 14-15, 16-18
		march30Slots := result["2025-03-30"]

		expectedSlots := []TimeSlot{
			makeTimeSlot("2025-03-30T09:00:00.000Z", "2025-03-30T10:00:00.000Z"),
			makeTimeSlot("2025-03-30T10:30:00.000Z", "2025-03-30T11:00:00.000Z"),
			makeTimeSlot("2025-03-30T11:30:00.000Z", "2025-03-30T12:00:00.000Z"),
			makeTimeSlot("2025-03-30T14:00:00.000Z", "2025-03-30T15:00:00.000Z"),
			makeTimeSlot("2025-03-30T16:00:00.000Z", "2025-03-30T18:00:00.000Z"),
		}

		assert.Len(t, march30Slots, len(expectedSlots))

		// Compare each slot
		for i, expected := range expectedSlots {
			assert.Equal(t, formatTime(expected.Start), formatTime(march30Slots[i].Start))
			assert.Equal(t, formatTime(expected.End), formatTime(march30Slots[i].End))
		}

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("should return empty array for days with no availability", func(t *testing.T) {
		// Reset mock
		mockDB.ExpectedCalls = nil

		startDate := parseTime("2025-04-06T00:00:00.000Z")
		endDate := parseTime("2025-04-13T00:00:00.000Z")

		// Setup mock to return empty events for this date range
		mockDB.On("QueryEvents", startDate, endDate).Return([]Event{})

		// Call the function under test
		result := CalculateAvailableSlots(mockDB, startDate)

		// Check the first day has empty slots
		assert.Empty(t, result["2025-04-06"])

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("should handle days with openings but all slots booked", func(t *testing.T) {
		// Reset mock
		mockDB.ExpectedCalls = nil

		startDate := parseTime("2025-04-06T00:00:00.000Z")
		endDate := parseTime("2025-04-13T00:00:00.000Z")

		// Add a day where there's an opening but it's fully booked
		fullyBookedEvents := []Event{
			{ID: 10, Kind: "opening", StartsAt: parseTime("2025-04-06T10:00:00.000Z"), EndsAt: parseTime("2025-04-06T11:00:00.000Z")},
			{ID: 110, Kind: "appointment", StartsAt: parseTime("2025-04-06T10:00:00.000Z"), EndsAt: parseTime("2025-04-06T11:00:00.000Z")},
		}

		// Setup mock
		mockDB.On("QueryEvents", startDate, endDate).Return(fullyBookedEvents)

		// Call the function under test
		result := CalculateAvailableSlots(mockDB, startDate)

		// Check that April 6 has no available slots
		assert.Empty(t, result["2025-04-06"])

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("should handle appointments that partially overlap with openings", func(t *testing.T) {
		// Reset mock
		mockDB.ExpectedCalls = nil

		startDate := parseTime("2025-04-06T00:00:00.000Z")
		endDate := parseTime("2025-04-13T00:00:00.000Z")

		// Add a day with a partial overlap
		partialOverlapEvents := []Event{
			{ID: 11, Kind: "opening", StartsAt: parseTime("2025-04-06T13:00:00.000Z"), EndsAt: parseTime("2025-04-06T15:00:00.000Z")},
			{ID: 111, Kind: "appointment", StartsAt: parseTime("2025-04-06T12:30:00.000Z"), EndsAt: parseTime("2025-04-06T13:30:00.000Z")},
		}

		// Setup mock
		mockDB.On("QueryEvents", startDate, endDate).Return(partialOverlapEvents)

		// Call the function under test
		result := CalculateAvailableSlots(mockDB, startDate)

		// The opening is 13:00-15:00, but appointment overlaps 12:30-13:30
		// So available slot should be 13:30-15:00
		april06Slots := result["2025-04-06"]

		expectedSlot := makeTimeSlot("2025-04-06T13:30:00.000Z", "2025-04-06T15:00:00.000Z")

		assert.Len(t, april06Slots, 1)
		assert.Equal(t, formatTime(expectedSlot.Start), formatTime(april06Slots[0].Start))
		assert.Equal(t, formatTime(expectedSlot.End), formatTime(april06Slots[0].End))

		// Verify mock expectations
		mockDB.AssertExpectations(t)
	})

	t.Run("should query the database for the correct date range", func(t *testing.T) {
		// Reset mock
		mockDB.ExpectedCalls = nil

		startDate := parseTime("2025-03-30T00:00:00.000Z")
		endDate := parseTime("2025-04-06T00:00:00.000Z")

		// Setup mock and capture the arguments
		mockDB.On("QueryEvents", startDate, endDate).Return([]Event{})

		// Call the function under test
		CalculateAvailableSlots(mockDB, startDate)

		// Verify that the database was called with the correct date range (7 days)
		mockDB.AssertCalled(t, "QueryEvents", startDate, endDate)
	})
}
