import Toggl from "./index";
import MockAdapter from "axios-mock-adapter";
import TimeEntry from "./data/time-entry";
import axios from "axios";

let mockAxios: MockAdapter;

describe("Test Toggl Service", () => {

  beforeAll(() => {
    mockAxios = new MockAdapter(axios);
  });

  afterEach(() => {
    mockAxios.resetHistory();
  });

  it("Should return time entries", async () => {
    const mockResponse: TimeEntry[] = [
      {
        "id": 1951596187,
        "wid": 1391549,
        "billable": false,
        "start": "2021-04-01T12:26:22+00:00",
        "stop": "2021-04-01T13:00:22+00:00",
        "duration": 2040,
        "description": "ISSUE-52 doing work",
        "at": "2021-04-01T13:04:51+00:00"
      },
      {
        "id": 1951664141,
        "wid": 1391549,
        "billable": false,
        "start": "2021-04-01T13:00:17+00:00",
        "stop": "2021-04-01T13:29:59+00:00",
        "duration": 1782,
        "description": "ISSUE-112 daily",
        "at": "2021-04-01T13:30:00+00:00"
      }
    ]
    const date: Date = new Date();
    date.setDate(date.getDate() - 1);

    mockAxios.onGet()
        .reply(200, mockResponse);
    const result = await Toggl.getTimeEntries(date.toISOString());
    expect(mockAxios.history.get.length).toBe(1)
    expect(result).toEqual(mockResponse);
  });

  it("Should error on non-200 response", () => {
    const date: Date = new Date();
    date.setDate(date.getDate() - 1);

    mockAxios.onGet().reply(403);
    const promise = Toggl.getTimeEntries(date.toISOString());

    promise.then(() => {
      throw new Error("Should enter .catch not .then");
    }).catch(err => {
      expect(err.isAxiosError).toBeTruthy();
    })
  });
});
