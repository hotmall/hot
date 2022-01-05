package commands

const goStart = `#!/bin/bash
EXEC=%s
PORT=53720
APP_HOME=$(cd "$(dirname $0)/..";pwd)
CMD="$APP_HOME/bin/$EXEC"
nohuplog="$APP_HOME/var/log/nohup.out"
nohup $CMD -p $APP_HOME -port $PORT >> $nohuplog 2>&1 &
`

const goStop = `#!/bin/bash
EXEC=%s
APP_HOME=$(cd "$(dirname $0)/..";pwd)
CMD="$APP_HOME/bin/$EXEC"
$CMD -p $APP_HOME -s stop
`
