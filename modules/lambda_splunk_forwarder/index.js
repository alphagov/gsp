/**
 * Stream events from AWS CloudWatch Logs to Splunk
 *
 * This function streams AWS CloudWatch Logs to Splunk using
 * Splunk's HTTP event collector API.
 *
 * Define the following Environment Variables in the console below to configure
 * this function to stream logs to your Splunk host:
 *
 * 1. SPLUNK_HEC_URL: URL address for your Splunk HTTP event collector endpoint.
 * Default port for event collector is 8088. Example: https://host.com:8088/services/collector
 *
 * 2. SPLUNK_HEC_TOKEN: Token for your Splunk HTTP event collector.
 * To create a new token for this Lambda function, refer to Splunk Docs:
 * http://docs.splunk.com/Documentation/Splunk/latest/Data/UsetheHTTPEventCollector#Create_an_Event_Collector_token
 */

'use strict';

const loggerConfig = {
    url: process.env.SPLUNK_HEC_URL || 'https://<HOST>:<PORT>/services/collector',
    token: process.env.SPLUNK_HEC_TOKEN || '<TOKEN>',
};

const SplunkLogger = require('./lib/mysplunklogger');
const zlib = require('zlib');

const logger = new SplunkLogger(loggerConfig);

exports.handler = (event, context, callback) => {
    // CloudWatch Logs data is base64 encoded so decode here
    const payload = Buffer.from(event.awslogs.data, 'base64');
    zlib.gunzip(payload, (err, result) => {
        if (err) {
            callback(err);
        } else {
            const parsed = JSON.parse(result.toString('ascii'));
            console.log('Event Data:', JSON.stringify(parsed, null, 2));
            let count = 0;
            if (parsed.logEvents) {
                parsed.logEvents.forEach((item) => {
                    // Send item JSON object (optional 'context' arg used to add Lambda metadata e.g. awsRequestId, functionName)
                    // Change "item.timestamp" below if time is specified in another field in the event
                    // Change to "logger.log(item.message, context)" if no time field is present in event
                    logger.logWithTime(item.timestamp, item.message, context);
                    count += 1;
                });
            }
            // Send all the events in a single batch to Splunk
            logger.flushAsync((error, response) => {
                if (error) {
                    callback(error);
                } else {
                    console.log(`Response from Splunk:\n${response}`);
                    console.log(`Successfully processed ${count} log event(s).`);
                    callback(null, count); // Return number of log events
                }
            });
        }
    });
};
