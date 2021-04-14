import TimeEntry from "./data/time-entry";
import axios, {AxiosResponse} from "axios";

class Toggl {

    static async getTimeEntries(startDate: string): Promise<TimeEntry[]> {
        const response: AxiosResponse = await axios.get(
            `https://api.track.toggl.com/api/v8/time_entries?start_date=${encodeURIComponent(startDate)}`,
            {
                auth: {
                    username: process.env.TOGGL_TOKEN || "",
                    password: "api_token"
                }
            }
        );

        return response.data;
    }
}

export default Toggl;
