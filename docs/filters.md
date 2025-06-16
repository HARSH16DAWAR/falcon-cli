# Falcon CLI Filters

The Falcon CLI provides a powerful filter system that allows you to save and reuse filters across different commands. This document explains how to use the filter system effectively.

## Overview

Filters in Falcon CLI are saved configurations that can be reused across different commands. They are stored in your configuration file and can be used with any command that supports filtering.

## Managing Filters

### Saving a Filter

To save a filter, use the `filter save` command:

```bash
falcon-cli filter save --name "windows-servers" \
                      --type hosts \
                      --description "All Windows servers" \
                      --filter "platform_name:'Windows'"
```

Required flags:
- `--name`: A unique name for your filter
- `--type`: The type of filter (e.g., hosts, detections)
- `--filter`: The filter expression to save

Optional flags:
- `--description`: A description of what the filter does

### Listing Filters

To list all saved filters:

```bash
# List all filters
falcon-cli filter list

# List filters of a specific type
falcon-cli filter list --type hosts
```

### Deleting a Filter

To delete a saved filter:

```bash
falcon-cli filter delete --name "windows-servers" --type hosts
```

Required flags:
- `--name`: The name of the filter to delete
- `--type`: The type of the filter to delete

## Using Saved Filters

### With Hosts Command

You can use saved filters with the hosts command in two ways:

1. Using the filter name:
```bash
falcon-cli hosts --filter-name "windows-servers"
```

2. Using a direct filter:
```bash
falcon-cli hosts --filter "platform_name:'Windows'"
```

Note: You cannot use both `--filter` and `--filter-name` at the same time.

## Filter Examples

Here are some example filters you can save:

### Host Filters

1. Windows Servers:
```bash
falcon-cli filter save --name "windows-servers" \
                      --type hosts \
                      --description "All Windows servers" \
                      --filter "platform_name:'Windows'"
```

2. Online Linux Hosts:
```bash
falcon-cli filter save --name "online-linux" \
                      --type hosts \
                      --description "All online Linux hosts" \
                      --filter "platform_name:'Linux'+status:'online'"
```

3. Critical Production Hosts:
```bash
falcon-cli filter save --name "critical-prod" \
                      --type hosts \
                      --description "Critical production hosts" \
                      --filter "tags:'critical'+tags:'production'"
```

## Filter Syntax

The filter syntax follows the CrowdStrike Falcon API filter format:

- Basic filters: `field:'value'`
- AND operator: `+` (e.g., `field1:'value1'+field2:'value2'`)
- OR operator: `,` (e.g., `field1:'value1',field1:'value2'`)
- NOT operator: `!` (e.g., `!field:'value'`)

## Best Practices

1. Use descriptive names for your filters
2. Always include a description to document the filter's purpose
3. Use consistent naming conventions for filter types
4. Regularly review and clean up unused filters
5. Use the `--type` flag to organize filters by command

## Configuration

Filters are stored in your Falcon CLI configuration file at `~/.falcon-cli/config.yaml`. The structure looks like this:

```yaml
filters:
  - name: "windows-servers"
    type: "hosts"
    description: "All Windows servers"
    filter: "platform_name:'Windows'"
  - name: "online-linux"
    type: "hosts"
    description: "All online Linux hosts"
    filter: "platform_name:'Linux'+status:'online'"
```

## Troubleshooting

1. If a filter is not found:
   - Check the filter name and type
   - Use `filter list` to see all available filters
   - Ensure the filter type matches the command you're using

2. If a filter doesn't work as expected:
   - Verify the filter syntax
   - Test the filter directly using the `--filter` flag
   - Check the CrowdStrike Falcon API documentation for valid filter fields 