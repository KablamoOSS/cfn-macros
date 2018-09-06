import traceback
import json

tag_types = [
    "AWS::AutoScaling::AutoScalingGroup", # PropagateAtLaunch
    "AWS::CertificateManager::Certificate",
    "AWS::CloudFormation::Stack",
    "AWS::CloudFront::Distribution",
    "AWS::CloudFront::StreamingDistribution",
    "AWS::CloudTrail::Trail",
    "AWS::CodeBuild::Project",
    # "AWS::Cognito::UserPool" -> UserPoolTags?
    # "AWS::DataPipeline::Pipeline" -> PipelineTags
    # "AWS::DAX::Cluster" Check
    "AWS::DMS::Endpoint",
    "AWS::DMS::EventSubscription",
    "AWS::DMS::ReplicationInstance",
    "AWS::DMS::ReplicationSubnetGroup",
    "AWS::DMS::ReplicationTask",
    "AWS::DynamoDB::Table",
    "AWS::EC2::CustomerGateway",
    "AWS::EC2::DHCPOptions",
    "AWS::EC2::Instance",
    "AWS::EC2::InternetGateway",
    "AWS::EC2::NatGateway",
    "AWS::EC2::NetworkAcl",
    "AWS::EC2::NetworkInterface",
    "AWS::EC2::RouteTable",
    "AWS::EC2::SecurityGroup",
    "AWS::EC2::Subnet",
    "AWS::EC2::Volume",
    "AWS::EC2::VPCPeeringConnection",
    "AWS::EC2::VPNConnection",
    "AWS::EC2::VPNGateway",
    # "AWS::EFS::FileSystem" -> FileSystemTags
    "AWS::ElastiCache::CacheCluster",
    "AWS::ElastiCache::ReplicationGroup",
    "AWS::ElasticBeanstalk::Environment",
    "AWS::ElasticLoadBalancing::LoadBalancer",
    "AWS::ElasticLoadBalancingV2::LoadBalancer",
    "AWS::ElasticLoadBalancingV2::TargetGroup",
    "AWS::Elasticsearch::Domain",
    "AWS::EMR::Cluster",
    # "AWS::Inspector::AssessmentTemplate" -> "UserAttributesForFindings"
    # "AWS::Inspector::ResourceGroup" -> "ResourceGroupTags"
    "AWS::Kinesis::Stream",
    "AWS::KMS::Key",
    "AWS::Lambda::Function",
    "AWS::Neptune::DBCluster",
    "AWS::Neptune::DBClusterParameterGroup",
    "AWS::Neptune::DBInstance",
    "AWS::Neptune::DBParameterGroup",
    "AWS::Neptune::DBSubnetGroup",
    "AWS::OpsWorks::Layer",
    "AWS::OpsWorks::Stack",
    "AWS::RDS::DBCluster",
    "AWS::RDS::DBClusterParameterGroup",
    "AWS::RDS::DBInstance",
    "AWS::RDS::DBParameterGroup",
    "AWS::RDS::DBSecurityGroup",
    "AWS::RDS::DBSubnetGroup",
    "AWS::RDS::OptionGroup",
    "AWS::Redshift::Cluster",
    "AWS::Redshift::ClusterParameterGroup",
    "AWS::Redshift::ClusterSecurityGroup",
    "AWS::Redshift::ClusterSubnetGroup",
    # "AWS::Route53::HealthCheck" -> "HealthCheckTags"
    # "AWS::Route53::HostedZone" -> "HostedZoneTags"
    "AWS::S3::Bucket",
    "AWS::SageMaker::Endpoint",
    "AWS::SageMaker::EndpointConfig",
    "AWS::SageMaker::Model",
    "AWS::SageMaker::NotebookInstance",
    "AWS::ServiceCatalog::CloudFormationProduct",
    "AWS::ServiceCatalog::CloudFormationProvisionedProduct",
    "AWS::ServiceCatalog::Portfolio",
    "AWS::SQS::Queue",
    "AWS::SSM::Document"
]

def handler(event, context):

    print(json.dumps(event))

    macro_response = {
        "requestId": event["requestId"],
        "status": "success"
    }
    try:
        params = {
            "params": event["templateParameterValues"],
            "template": event["fragment"],
            "account_id": event["accountId"],
            "region": event["region"]
        }
        response = event["fragment"]

        tags = []
        # TODO: Check existence of below
        tags = response["Metadata"]["AddTags"]

        for k in list(response["Resources"].keys()):
            for tag_type in tag_types:
                if response["Resources"][k]["Type"] == tag_type:
                    if "Properties" not in response["Resources"][k]:
                        response["Resources"][k]["Properties"] = {}
                    if "Tags" not in response["Resources"][k]["Properties"]:
                        response["Resources"][k]["Properties"]["Tags"] = []
                    for tag in tags:
                        res_tag = {
                            'Key': tag["Key"],
                            'Value': ""
                        }
                        if "Value" in tag:
                            res_tag["Value"] = tag["Value"]
                        if response["Resources"][k]["Type"] == "AWS::AutoScaling::AutoScalingGroup" and "PropagateAtLaunch" in tag:
                            res_tag["PropagateAtLaunch"] = tag["PropagateAtLaunch"]
                        response["Resources"][k]["Properties"]["Tags"].append(res_tag)
                                
        macro_response["fragment"] = response
    except Exception as e:
        traceback.print_exc()
        macro_response["status"] = "failure"
        macro_response["errorMessage"] = str(e)

    print(json.dumps(macro_response))
    return macro_response