import Toggl from "./services/toggle";
import TimeEntry from "./services/toggle/data/time-entry";
import Jira from "./services/jira";

export default async function main(): Promise<void> {
    const date: Date = new Date();

    date.setDate(date.getDate() - 1);
    const timeEntries: TimeEntry[] = await Toggl.getTimeEntries(date.toISOString());

    if (timeEntries.length === 0) {
        console.log("No time entries found in the last 24 hours");
        return;
    }

    for (const timeEntry of timeEntries) {
        await Jira.addTimeToWorkLog(timeEntry);
    }
}
