#!/bin/sh
echo "starting cloudwatch test"

set -x
#set -xeuf -o pipefail

timeout="${TEST_TIMEOUT:-60}"
retries="${TEST_RETRIES:-3}"

echo "accountname: $ACCOUNT_NAME"
echo "clustername: $CLUSTER_NAME"
echo "timeout: $timeout"
echo "retries: $retries"
echo "testlogsince: $TEST_LOGS_SINCE"


i=0
while [ $i -le $retries ]
do
  ((i++))
  echo "attempt: $i"
  LASTSEENLOG=$(aws logs describe-log-streams --log-group-name $CLUSTER_NAME.$ACCOUNT_NAME.govsvc.uk --log-stream-name-prefix "kubernetes.var.log.containers.sonobuoy_heptio-sonobuoy_istio-proxy" | jq ".logStreams[].lastEventTimestamp" | grep -v "null" | sort -urn | head -n1)

  echo "lastseen: $LASTSEENLOG"

  if (( ${LASTSEENLOG} > ${TEST_LOGS_SINCE} )); then
    echo "PASS: Logs have been reached cloudwatch\nAfter: $TEST_LOGS_SINCE at $LASTSEENLOG in $CLUSTER_NAME.$ACCOUNT_NAME.govsvc.uk/kubernetes.var.log.containers.sonobuoy_heptio-sonobuoy_istio-proxy" 2>&1 | tee /tmp/results
    exit 0
  fi

  sleep ${timeout}
done

echo "FAIL: No logs have been detected reaching cloudwatch since $TEST_LOGS_SINCE" 2>&1 | tee /tmp/results
exit 1
