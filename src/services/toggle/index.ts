import TimeEntry from "./data/time-entry.js";
import axios, {AxiosResponse} from "axios";

class Toggl {

    static async getTimeEntries(startDate: string): Promise<TimeEntry[]> {
        const endDate = new Date().toISOString();
        const token = process.env.TOGGL_TOKEN || "";
        const response: AxiosResponse = await axios.get(
            `https://api.track.toggl.com/api/v9/me/time_entries?start_date=${encodeURIComponent(startDate)}&end_date=${encodeURIComponent(endDate)}`,
            {
                auth: {
                    username: token,
                    password: 'api_token'
                },
                headers: {
                    'Content-Type': 'application/json'
                }
            }
        );

        return response.data;
    }
}

export default Toggl;
