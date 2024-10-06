# popci

- one endpoint /webhook
- receives payload from Github containing the event that triggered the webhook
- payload contains
    * git commit
    * workflow file?
- runner executes the script in the workflow file and returns result (where?)

## questions

- how to multiplex the same endpoint to accept different payloads from Github, Gitlab, etc?
- how to grant CI access to (private) repository?