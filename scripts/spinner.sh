#!/bin/bash

function spin {
  spinner="/|\\-/|\\-"
  while :
  do
    for i in `seq 0 7`
    do
      echo -n "${spinner:$i:1}"
      echo -en "\010"
      sleep 0.15
    done
  done
}

function start_spinner {
    spin &
    export _sp_pid=$!
    trap "kill -9 $_sp_pid 2>/dev/null" `seq 0 15`
    disown
}
