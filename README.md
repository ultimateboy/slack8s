# Kubernetes Events Slack Bot (slack8s)

Kubernetes Slack Integration - Infinite loop go program which queries the Kubernetes Event Stream API and
posts messages to slack for important events.

![Slack8s demo showing creation of pod and then failed backoff loop alerts via Slack.](images/slack8s-demo.png)

## Building

Given this is an active work in progress, you'll probably want to modify
and build the code yourself. If you just want to get started quickly,
use the image in [docker hub](https://hub.docker.com/r/ultimateboy/slack8s/).

1. `docker build .`
2. Tag and push to your favorite registry

## Running

1. [Create a new Slack Bot](https://my.slack.com/services/new/bot).
2. Copy the example configmap file:  
`cp examples/example.slack8s-configmap.yaml examples/slack8s-configmap.yaml`
3. Modify `slack-token` and `slack-channel` variables in your new file.
4. Create the config map using kubectl:  
`kubectl create -f examples/slack8s-configmap.yaml`
5. Create the slack8s deployment:  
`kubectl create -f examples/slack8s-deployment.yaml`

## Limitations

1. This does not keep track of which notifications it has sent. You should not
trust that slack has every event.
2. When the container dies, if the container takes more than 1 minute to spawn,
some events will not be posted to slack.
3. When the container dies, events within 1 minute of the restart will be posted
to Slack twice.
4. The watch api seems to hit EOF causing this container to die more often than
desired. On a busy cluster, the above limitations are more obvious because of this.
5. Do not attempt to use this as a replacement for a pager.

## Todo

1. Refactor the way in which the types of alerts to send to slack are configured
2. Add more types of alerts which get posted to slack
3. Better formatting of slack attachment output
4. Alerting thresholds
