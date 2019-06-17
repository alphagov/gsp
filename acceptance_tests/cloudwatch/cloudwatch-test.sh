#!/usr/bin/env bash

set -x

echo "starting cloudwatch test"

#set -xeuf -o pipefail

# AWS_DEFAULT_REGION=eu-west-2
# AWS_REGION=eu-west-2

TIMEOUT="${TEST_TIMEOUT:-30}"
RETRIES="${TEST_RETRIES:-3}"
FARBACK="${TEST_FARBACK:-300}"
#ACCOUNT_NAME="portfolio"
#CLUSTER_NAME="portfolio"
TEST_LOGS_SINCE=$(date '+%s')
echo "Current time: $TEST_LOGS_SINCE"
TEST_LOGS_SINCE=$(($TEST_LOGS_SINCE - $FARBACK))
echo "$FARBACK Secs ago: $TEST_LOGS_SINCE"

echo "accountname: $ACCOUNT_NAME"
echo "clustername: $CLUSTER_NAME"
echo "clusterdomain: $CLUSTER_DOMAIN"
echo "timeout: $TIMEOUT"
echo "retries: $RETRIES"
echo "testlogsince: $TEST_LOGS_SINCE"


i=0
while [ $i -lt $RETRIES ]
do
  ((i++))
  echo "Attempt: $i"
  # LASTSEENLOG=$(AWS_DEFAULT_REGION=eu-west-2 AWS_REGION=eu-west-2 aws logs describe-log-streams --log-group-name $CLUSTER_NAME.$ACCOUNT_NAME.govsvc.uk --log-stream-name-prefix "kubernetes.var.log.containers.sonobuoy_heptio-sonobuoy_istio-proxy" | jq ".logStreams[].lastEventTimestamp" | grep -v "null" | sort -urn | head -n1)

  LASTSEENLOG=$(AWS_DEFAULT_REGION=$AWS_DEFAULT_REGION AWS_REGION=$AWS_REGION aws logs describe-log-streams --log-group-name $CLUSTER_DOMAIN --log-stream-name-prefix "kubernetes.var.log.containers.sonobuoy_heptio-sonobuoy_istio-proxy" | jq ".logStreams[].lastEventTimestamp" | grep -v "null" | sort -urn | head -n1)

  echo "lastseen: $LASTSEENLOG"

  if (( ${LASTSEENLOG} > ${TEST_LOGS_SINCE} )); then
    echo "PASS: Logs have been reached cloudwatch\nAfter: $TEST_LOGS_SINCE at $LASTSEENLOG in $CLUSTER_DOMAIN/kubernetes.var.log.containers.sonobuoy_heptio-sonobuoy_istio-proxy" 2>&1 | tee /tmp/results
    exit 0
  fi

  sleep ${TIMEOUT}
done

echo "FAIL: No logs have been detected reaching cloudwatch since $TEST_LOGS_SINCE" 2>&1 | tee /tmp/results
exit 1
