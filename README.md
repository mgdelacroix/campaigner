# Campaigner

Command line tool to create and manage community campaigns.

## Application flow

1. Run `campaigner set-token` to set the tokens for `jira/github`.
2. Run `campaigner create` to create a new community campaign.
3. Run `campaigner add` to add new tickets based either on a `grep/ag`
   command or a `govet` check.
4. Run `campaigner state` to see the status of each one of the
   tickets. The tickets can be in an `unpublished`, `jira`, `github`
   and `completed` state.
5. Run `campaigner template` to edit the ticket template.
6. Run `campaigner publish` to create the tickets in `jira` based on
   the template.
7. Possible next step to publish the tickets from `jira` to `github`.
