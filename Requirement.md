`# User story

The resource has 2 states at the moment:

- Stopped
- Working

In reality, any resource can be in at least 2 states:

- Setup
- Process

The resource can be set up and the desired amount of material can be prepared before it starts.
After the start of processing, the successfully processed material is logged.

## Required changes

Add 2 new statuses for the resource instead of working:

- Setup
- Process

For each new status it must be possible to register a quantity.

The ResourceHub must listen to the new **setup** and **process** endpoints and change the current status of the resource.
`