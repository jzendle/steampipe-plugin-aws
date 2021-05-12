select arn, ec2_instance_id, expire_passwords, status
from aws.aws_ecs_container_instance where arn = "{{ output.arn.value }}"

