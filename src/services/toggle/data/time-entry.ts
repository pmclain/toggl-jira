export default interface TimeEntry {
    id: number;
    wid: number;
    billable: boolean;
    start: string;
    stop?: string;
    duration: number;
    description: string;
    at: string;
}
