import axios, {AxiosResponse} from "axios";
import TimeEntry from "../toggle/data/time-entry";
import WorkLog from "./data/work-log";
import moment from "moment";

class Jira {
    private static auth() {
        return {
            username: process.env.JIRA_USER || "",
            password: process.env.JIRA_TOKEN || ""
        }
    }

    private static jiraBaseUri() {
        return process.env.JIRA_BASE_URI || "";
    }

    static async addTimeToWorkLog(time: TimeEntry): Promise<void> {
        if (time.duration < 1 || !time.stop) {
            console.log(`${time.description}: Skipped running timer`);
            return;
        }
        const timeEntryIssues: string[] = Jira.extractIssuesFromTimeEntry(time);
        if (timeEntryIssues.length === 0) {
            console.log(`${time.description}: No supported issue found`);
            return;
        }

        //TODO: split time between issues
        for (const issueKey of timeEntryIssues) {
            try {
                const existingWorkLog = await Jira.getWorkLogsForIssue(issueKey, time);
                if (existingWorkLog) {
                    await Jira.updateWorkLog(issueKey, existingWorkLog, time);
                } else {
                    await Jira.createWorkLog(issueKey, time);
                }
            } catch (e) {
                // TODO add an actual logger
                console.error(`${issueKey}: Unable to load issue`);
            }
        }
    }

    private static async updateWorkLog(issueKey: string, workLog: WorkLog, time: TimeEntry): Promise<void> {
        if (workLog.timeSpentSeconds === Jira.roundUp(time.duration)) {
            console.log(`${time.description}: up to date`);
            return;
        }

        console.log(`${time.description}: updating existing worklog`);

        await axios.put(
            `${Jira.jiraBaseUri()}/rest/api/latest/issue/${issueKey}/worklog/${workLog.id}?notifyUsers=false`,
            {
                comment: `TogglID: ${String(time.id)} ${time.description}`,
                timeSpentSeconds: Jira.roundUp(time.duration),
                started: Jira.formatDateTime(time.start)
            },
            {
                auth: Jira.auth()
            }
        );
    }

    private static async createWorkLog(issueKey: string, time: TimeEntry): Promise<void> {
        console.log(`${time.description}: creating worklog`);

        await axios.post(
            `${Jira.jiraBaseUri()}/rest/api/latest/issue/${issueKey}/worklog?notifyUsers=false`,
            {
                comment: `TogglID: ${String(time.id)} ${time.description}`,
                timeSpentSeconds: Jira.roundUp(time.duration),
                started: Jira.formatDateTime(time.start)
            },
            {
                auth: Jira.auth()
            }
        );
    }

    private static roundUp(duration: number): number {
        const roundUpToNext = 900;
        return duration + (roundUpToNext - (duration % roundUpToNext));
    }

    private static formatDateTime(time: string): string {
        return moment(time).format("YYYY-MM-DDTHH:mm:ss.SSSZZ");
    }

    private static getSupportedJiraKeys(): string[] {
        const keys: string[] = (process.env.JIRA_PROJECTS || "").split(",");
        return keys.map(string => string.trim());
    }

    private static extractIssuesFromTimeEntry(time: TimeEntry): string[] {
        const supportedKeys: string[] = Jira.getSupportedJiraKeys();
        if (supportedKeys.length === 0) {
            return [];
        }

        const keyRegex = new RegExp(`((?:${supportedKeys.join("|")})-\\d+)`, "g");
        const matches = time.description?.match(keyRegex);

        return matches || [];
    }

    private static async getWorkLogsForIssue(issueKey: string, time: TimeEntry): Promise<WorkLog | null> {
        const response: AxiosResponse = await axios.get(
            `${Jira.jiraBaseUri()}/rest/api/latest/issue/${issueKey}/worklog?comment`,
            {
                auth: Jira.auth()
            }
        );

        for (const workLog of response.data.worklogs || []) {
            if (workLog.comment?.includes(String(time.id))) {
                return workLog;
            }
        }

        return null;
    }
}

export default Jira;
