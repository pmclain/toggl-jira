# Toggl > JIRA Sync

Sync Toggl time entries into JIRA Work Logs.

### Setup

* `yarn install`
* `cp .env.sample .env`

Add values for `.env` variables

| Variable | Description |
| --- | --- |
| `JIRA_TOKEN` | JIRA API Key [Atlassian Docs](https://support.atlassian.com/atlassian-account/docs/manage-api-tokens-for-your-atlassian-account/) |
| `JIRA_USER` | The email or username of your JIRA account. |
| `JIRA_PROJECTS` | Comma separated list of project keys to sync ie `PROJONE,PROJTWO`. Only issues from these projects will by synced. |
| `JIRA_BASE_URI` | The base URI of your JIRA instance. |
| `TOGGL_TOKEN` | Toggle API Key [Toggl Docs](https://github.com/toggl/toggl_api_docs#api-token) |

### Usage

`yarn sync-time`

* Pulls Toggl entries for the last 24hrs
* Converts Toggl entries to JIRA work logs for `JIRA_PROJECTS` based on the
  Toggle description ie `ISSUE-1 Taking care of business`
* Rounds time entries up to next 15 minute increment.
* Updates work log duration when JIRA ticket and Toggl ID match.

### Known Limitation

* Time entries cannot be disassociated with tickets after sync. For example:
  A time entry with a description `ISSUE-1 I did work` is synced to JIRA.
  Changes to the duration will update the associated work log, but changes to
  the ticket number will result in a new work log and the old work log will NOT
  be deleted.
* Time entries cannot apply to multiple tickets
