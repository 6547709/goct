# Task Commands

CloudTower uses asynchronous tasks for most mutation operations (VM creation, power operations, migrations, etc.). Tasks track operation progress and completion status.

## Listing Tasks
```bash
goct task.ls                           # List recent tasks
goct task.ls --id-only                # Output only task IDs
goct --format json task.ls            # JSON output

# Filter tasks by status
goct task.ls | grep RUNNING           # Show running tasks
goct task.ls | grep FAILED            # Show failed tasks
goct task.ls | grep SUCCESS           # Show completed tasks
```

## Task Info
```bash
goct task.info <task-id>              # Show task details and progress
```

## Waiting for Tasks
```bash
goct task.wait <task-id>              # Wait for task to complete
goct task.wait <task-id> --timeout 300  # Wait with timeout (seconds)
```

## Canceling Tasks
```bash
goct task.cancel <task-id>            # Cancel a running task
```

## Understanding Task States

| State | Meaning |
|-------|---------|
| `RUNNING` | Operation in progress |
| `SUCCESS` | Operation completed successfully |
| `FAILED` | Operation failed (check error message) |
| `CANCELLED` | Operation was cancelled |

## Notes

- Most goct mutation commands automatically wait for task completion
- Use `goct task.wait` when working with task IDs directly
- Failed tasks show error reason in task info output