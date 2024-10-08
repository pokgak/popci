# popci

- [X] one endpoint /webhook
- [X] receives payload from Github containing the event that triggered the webhook
    - payload contains
        * git commit
        * workflow file?
- [X] runner executes the script in the workflow file
- [ ] store result for user to read later
- [ ] update github commit status

## questions

- how to multiplex the same endpoint to accept different payloads from Github, Gitlab, etc?
- how to grant CI access to (private) repository?