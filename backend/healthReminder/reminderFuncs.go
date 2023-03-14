package healthReminder

import (
	"time"

	"github.com/KnockOutEZ/Kodee/backend/utils"
)

var notify = utils.NotificationFunc

func LookAwayReminder() {
    for {
        time.Sleep(20 * time.Minute)

        notify("Health Reminder!!!", "Please look away for 20 seconds")
    }
}

func StandUpReminder() {
    for {
        time.Sleep(1 * time.Hour)

        notify("Health Reminder!!!", "Please stand up and walk for 5 minute. You can also do some stretching")
    }
}

func HydrateReminder() {
    for {
        time.Sleep(30 * time.Minute)

        notify("Health Reminder!!!", "Please drink a glass of water")
    }
}