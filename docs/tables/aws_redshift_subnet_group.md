# Table: aws_redshift_subnet_group

A cluster subnet group is a collection of subnets (typically private) that are created for a VPC and then designated for redshift clusters.

## Examples

### DB subnet group basic info

```sql
select
  cluster_subnet_group_name,
  description,
  subnet_group_status,
  vpc_id
from
  aws_redshift_subnet_group;
```


### Subnets info of each subnet in subnet group

```sql
select
  cluster_subnet_group_name,
  subnet -> 'SubnetAvailabilityZone' ->> 'Name' as subnet_availability_zone,
  subnet -> 'SubnetAvailabilityZone' ->> 'SupportedPlatforms' as supported_platforms,
  subnet ->> 'SubnetIdentifier' as subnet_identifier,
  subnet ->> 'SubnetStatus' as subnet_status
from
  aws_redshift_subnet_group
  cross join jsonb_array_elements(subnets) as subnet;
```


### List of subnet group without application tag key

```sql
select
  name,
  tags
from
  aws_redshift_subnet_group
where
  not tags :: JSONB ? 'application';
```