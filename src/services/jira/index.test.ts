import Jira from "./index";
import axios from "axios";
import MockAdapter from "axios-mock-adapter";
import TimeEntry from "../toggle/data/time-entry";
import WorkLog from "./data/work-log";

const timeEntry: TimeEntry = {
  "id": 1951596187,
  "wid": 1391549,
  "billable": false,
  "start": "2021-04-01T12:26:22+00:00",
  "stop": "2021-04-01T13:00:22+00:00",
  "duration": 2040,
  "description": "ISSUE-52 doing work",
  "at": "2021-04-01T13:04:51+00:00"
};

const timeEntryMultipleIssues: TimeEntry = {
  "id": 1951596187,
  "wid": 1391549,
  "billable": false,
  "start": "2021-04-01T12:26:22+00:00",
  "stop": "2021-04-01T13:00:22+00:00",
  "duration": 2040,
  "description": "ISSUE-52,ISSUE-55 doing work",
  "at": "2021-04-01T13:04:51+00:00"
};

const timeEntryMultipleIssuesWhiteSpace: TimeEntry = {
  "id": 1951596187,
  "wid": 1391549,
  "billable": false,
  "start": "2021-04-01T12:26:22+00:00",
  "stop": "2021-04-01T13:00:22+00:00",
  "duration": 2040,
  "description": "ISSUE-52, ISSUE-55 doing work",
  "at": "2021-04-01T13:04:51+00:00"
};

const timeEntryMultipleIssuesNoComma: TimeEntry = {
  "id": 1951596187,
  "wid": 1391549,
  "billable": false,
  "start": "2021-04-01T12:26:22+00:00",
  "stop": "2021-04-01T13:00:22+00:00",
  "duration": 2040,
  "description": "ISSUE-52 ISSUE-55 doing work",
  "at": "2021-04-01T13:04:51+00:00"
};

const timeEntryMultipleIssuesNoSeparator: TimeEntry = {
  "id": 1951596187,
  "wid": 1391549,
  "billable": false,
  "start": "2021-04-01T12:26:22+00:00",
  "stop": "2021-04-01T13:00:22+00:00",
  "duration": 2040,
  "description": "ISSUE-52ISSUE-55 doing work",
  "at": "2021-04-01T13:04:51+00:00"
};

let mockAxios: MockAdapter;

describe("Test Jira Service", () => {

  beforeAll(() => {
    mockAxios = new MockAdapter(axios);
  });

  beforeEach(() => {
    process.env.JIRA_PROJECTS = "ISSUE";
  });

  afterEach(() => {
    mockAxios.resetHistory();
  });

  it("Should skip entry without issue", async () => {
    const timeEntryWithoutIssue: TimeEntry = {
      "id": 1951596187,
      "wid": 1391549,
      "billable": false,
      "start": "2021-04-01T12:26:22+00:00",
      "stop": "2021-04-01T13:00:22+00:00",
      "duration": 2040,
      "description": "doing work",
      "at": "2021-04-01T13:04:51+00:00"
    };

    await Jira.addTimeToWorkLog(timeEntryWithoutIssue);
    expect(mockAxios.history.get.length).toBe(0);
  });

  it("Should create new worklog", async () => {
    mockAxios.onGet()
        .replyOnce(200, { worklogs: [] });
    mockAxios.onPost().reply(200);
    await Jira.addTimeToWorkLog(timeEntry);
    const postData = JSON.parse(mockAxios.history.post[0].data);
    expect(postData.comment).toContain(`TogglID: ${timeEntry.id}`);
    expect(postData.timeSpentSeconds % 900).toBe(0)
  });

  it("Should ignore running time entries", async () => {
    const runningTimeEntry: TimeEntry = {
      "id": 1952399902,
      "wid": 1391549,
      "billable": false,
      "start": "2021-04-01T21:17:21+00:00",
      "duration": -1617311842,
      "description": "ISSUE-52 hello",
      "at": "2021-04-01T21:17:31+00:00"
    };

    await Jira.addTimeToWorkLog(runningTimeEntry);
    expect(mockAxios.history.get.length).toBe(0);
  });

  it("Should update existing worklog", async () => {
    const mockWorkLog: WorkLog = {
      id: "123",
      issueId: "456",
      comment: `TogglID: ${timeEntry.id}`,
      timeSpentSeconds: 5
    };

    mockAxios.onGet()
        .replyOnce(200, { worklogs: [mockWorkLog] });
    mockAxios.onPut().reply(200);
    await Jira.addTimeToWorkLog(timeEntry);
    expect(mockAxios.history.post.length).toBe(0);
    const putData = JSON.parse(mockAxios.history.put[0].data);
    expect(putData.comment).toContain(`TogglID: ${timeEntry.id}`);
    expect(putData.timeSpentSeconds % 900).toBe(0)
  });

  it("Should ignore unsupported project keys", async () => {
    process.env.JIRA_PROJECTS = "OTHER";
    await Jira.addTimeToWorkLog(timeEntry);
    expect(mockAxios.history.get.length).toBe(0);
  });

  it.skip("Should split time between issues", async () => {
    expect(false).toBeTruthy();
  });
});
