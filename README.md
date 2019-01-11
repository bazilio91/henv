Stupid /etc/hosts manager

To install:
- clone this repo
- ensure that your golang version supports go modules
- run go install

To start - fill your `~/.henv.yml` like this:
```yaml
services:
  app1:
    env1:
      127.0.0.1:
        - app1.example.com
    env2:
      192.168.0.1:
        - app1.example.com

  app2:
    env1:
      127.0.0.1:
        - app2.example.com
        - www.app2.example.com
```

Commands:
```bash
henv app2 env2 # apply env1 to app2
henv app2 undo # remove app2 rules

henv all env1 # apply env1 on all apps (if defined)
henv all undo # remove all rules
```
