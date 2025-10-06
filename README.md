# THOR structure definitions

## Introduction

This library provides definitions of structures used in the output of the THOR APT Forensic Scanner. These structures can be used for different use cases:
- generate a schema for THOR JSON logs
- convert JSON logs into text logs
- parse JSON logs

## Versions

There are three versions of the THOR log format:

 - v1: The original THOR log format, used up to and including THOR version 10.7. This is equivalent to the THOR text format, simply serialized as JSON.
 - v2: The format used in THOR version 10.7 with the `--jsonv2` flag. This format introduced a more structured approach to logging,
   with subobjects for reasons, files, and other entities. It is largely open-ended and allows for custom fields.
 - v3: The format used in THOR 11 and later. This format is more strict and versioned, with a defined schema. It introduces the concept of _reportable objects_.

## Parsing Events

There is a parser in the `thorlog/parser` package which can be used to parse an event.
This parser is version aware and can handle all versions of the THOR log format.
The result of the parsing is a `common.Event` object, which is a version-agnostic representation of a THOR event.
It can be cast to the version-specific implementation of this interface, e.g. `thorlog.Finding` for a finding in version 3.

## Textlog Conversion

The `jsonlog.TextlogFormatter` type provides a way to convert an object to a text log format.

This formatter can be used to convert findings and messages to a human-readable format.
However, the text log format is not as rich as the JSON format and may not contain all fields.
When in doubt, use the JSON format for analysis.

## Objects in JSON Log Version 3

Each object in the THOR log contains a `type` field that indicates the object type.
This type determines how the object should be interpreted and what fields it contains.

### Event Types

The object types contained in a THOR log are `THOR finding` and `THOR message`:
 - Findings are the results of THOR's analysis, such as detected threats or anomalies.
 - Messages are informational or status updates from THOR, such as progress updates.

Both findings and messages are together called _events_.

### Reportable Objects

Findings may contain more objects, e.g. as a subject that they report. 
Object types that can appear as subjects are called _reportable objects_.
The most common reportable objects are:
- `file`
- `process`

Reportable objects should contain only fields that relate directly to the object itself.
E.g. when extracting a file from an archive, the file object should contain only fields 
that relate to the file itself, not to the archive.
The archive data will instead appear in the _context_ of the finding.

## Schema

A schema for the version 3 format is attached to each release.
It can also be generated using the `thorlog/jsonschema` package.
