import json
import logging
from collections import defaultdict
logger = logging.getLogger()
logger.setLevel(logging.DEBUG)

def _implicit_api_URL(*_):
    return {
        "Description": "API endpoint URL",
        "Export": {
            "Name": {"Fn::Sub": "${AWS::StackName}-Endpoint"}
        },
        "Value":{
            "Fn::Join": [
                ".", [
                {
                    "Fn::Sub": "https://${ServerlessRestApi}"
                },
                "execute-api",
                {
                    "Ref": "AWS::Region"
                },
                {
                    "Fn::Sub": "amazonaws.com/${ServerlessRestApi.Stage}"
                },
                ]
            ]
        }
    }

def _explicit_api_URL(name):
    return {
        "Description": "API endpoint URL",
        "Export": {
            "Name": {"Fn::Sub": "${AWS::StackName}-Endpoint"}
        },
        "Value":{
            "Fn::Join": [
                ".", [
                {
                    "Fn::Sub": f"https://${{{name}}}"
                },
                "execute-api",
                {
                    "Ref": "AWS::Region"
                },
                {
                    "Fn::Sub": f"amazonaws.com/${{{name}.Stage}}"
                },
                ]
            ]
        }
    }

def _generic_output(resource_name, parameter=None):
    return {
        "Export": {
            "Name": {"Fn::Sub": "${AWS::StackName}-" + resource_name}
        },
        "Value": {"Fn:GetAtt": f"{resource_name}.{parameter}"} if parameter else {"Ref": resource_name}
    }


generate_output = defaultdict(
    lambda: _generic_output, {
        'AWS::Serverless::Function': lambda x: _generic_output(x, 'Alias'),
        'AWS::Serverless::Api': _explicit_api_URL,
    }
)


def do_outputs(fragment):
    fragment['Outputs'] = {
        **fragment.get('Outputs', {}),
        **{name: generate_output[specs['Type']](name)
            for name, specs in fragment['Resources'].items()
        }
    }


def main(event, context):
    event_txt = json.dumps(event, indent=4, separators=(',', ': '), sort_keys=True)
    logger.debug(f"event: {event_txt}")

    # Add general Outputs
    do_outputs(event['fragment'])
    # Finally add the URL Output if one exists
    if '"Type": "Api"' in event_txt and not 'RestApiId' in event_txt:
        event['fragment']['Outputs']['ImplicitApiURL'] = _implicit_api_URL()

    fragment_txt = json.dumps(event['fragment'], indent=4, separators=(',', ': '), sort_keys=True)
    logger.debug(f"Returning fragment: {fragment_txt}")
    return {
        "requestId": event["requestId"],
        "status": "success",
        "fragment": event["fragment"],
    }

if __name__ == '__main__':
    ch = logging.StreamHandler()
    ch.setLevel(logging.DEBUG)
    logger.addHandler(ch)

    event_file = '../test_events/test1.json'
    with open(event_file, 'r') as f:
        event = json.load(f)
    main(event, {})
