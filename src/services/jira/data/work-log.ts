interface Author {
    self: string;
    accountId: string;
    emailAddress: string;
    displayName: string;
    active: boolean;
    timeZone: string;
    accountType: string;
}

export default interface WorkLog {
    self?: string;
    author?: Author;
    updateAuthor?: Author;
    comment: string;
    created?: string;
    updated?: string;
    started?: string;
    timeSpent?: string;
    timeSpentSeconds: number;
    id: string;
    issueId: string;
}
