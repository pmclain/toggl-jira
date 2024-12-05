import Toggl from "./services/toggle/index.js";
import TimeEntry from "./services/toggle/data/time-entry.js";
import Jira from "./services/jira/index.js";

export default async function main(): Promise<void> {
    const date: Date = new Date();

    date.setDate(date.getDate() - 2);
    const timeEntries: TimeEntry[] = await Toggl.getTimeEntries(date.toISOString());

    if (timeEntries.length === 0) {
        console.log("No time entries found in the last 48 hours");
        return;
    }

    for (const timeEntry of timeEntries) {
        await Jira.addTimeToWorkLog(timeEntry);
    }
}
