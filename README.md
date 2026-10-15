# THOR structure definitions

## Introduction

This library defines the structures in the output of the THOR APT Forensic Scanner. Use it to:

- generate schemas for THOR JSON logs and the THOR audit trail
- parse these formats into Go structures
- convert JSON logs into text logs

## Log Formats

THOR has two separate log formats:

- The **THOR log** contains THOR's verdicts in the form of _assessments_.
  Each assessment rates a single _subject_: its severity and THOR's certainty about that (the `score`)
  and why (the `reasons`).
  The subject is the focus of the assessment. All other objects in it are only context for that subject.

  Each assessment contains all information necessary to understand it.
  As a consequence, the same object can appear in many assessments.
- The **audit trail** contains an _audit record_ for each _object_ that THOR observed, suspicious or not.
  Unlike an assessment, an audit record has no overall score.
  Its `reasons` field lists the indicators that matched on the object, if any.
  Related objects each get their own audit record, and the `references` field links these records by ID.

  The records form a graph that you can walk in either direction,
  e.g. from an archive to all files extracted from it, or from a file up to its archive.

The two formats differ in how they handle related objects:
the THOR log repeats a related object in every assessment it is relevant for,
while the audit trail records each object once and links to that record.
E.g., consider a ZIP archive that contains two suspicious files:

- The THOR log contains two assessments, one for each file.
  Each assessment includes the archive in its `ancestors`.
- The audit trail contains three audit records: one for the archive and one for each file.
  Each file's record includes a reference to the archive's record.

### When to use which log format

Use the THOR log to triage what THOR considers suspicious, starting with the highest scores.

Use the audit trail for an in-depth investigation starting from an initial suspicious object.

## Versions

There are three versions of the THOR log format:

- v1: The original THOR log format, used up to and including THOR version 10.7. This is equivalent to the THOR text format, serialized as JSON.
- v2: The format used in THOR version 10.7 with the `--jsonv2` flag. This format introduced a more structured approach to logging,
  with subobjects for reasons, files, and other entities. It is largely open-ended and allows for custom fields.
- v3: The format used in THOR 11 and later. This format is more strict and versioned, with a defined schema.

The audit trail also uses the version 3 format, with its own schema.

## Version 3 Objects

Each object in a version 3 log (THOR log or audit trail) has a `type` field.
The type determines how to interpret the object and which fields it contains.

The objects are defined in the `thorlog/v3` package; the audit trail entry types are defined in its `thorlog/v3/audittrail` subpackage.
Both logs largely share the same objects.

### Reportable Objects

Object types that can appear as the subject of an assessment or as the object of an audit record are called _reportable objects_.
The most common reportable objects are:

- `file`
- `process`

Reportable objects only contain fields that relate directly to the object itself.

E.g. when extracting a file from an archive, the file object contains fields
that relate to the file itself, but not fields relating to the archive.
The archive data will instead appear in the `ancestors` of a THOR log assessment,
or as a related audit record in the audit trail.

## THOR Log Event Types

The object types contained in a THOR log are `THOR assessment` and `THOR message`:

- An assessment is THOR's verdict on a single subject. Subjects are _reportable objects_ (see above), such as files and processes.

  The subject is not necessarily suspicious.
  Its score (0–100) reflects severity and certainty, and its reasons list the indicators that contributed to that score. 
  Investigate assessments with high scores first.
- Messages are informational or status updates from THOR, such as progress updates.

Both assessments and messages are together called _events_.

## Audit Trail Entry Types

The object types contained in the audit trail are `THOR audit record` and `THOR audit message`:

- Audit records document a single object that THOR observed, together with the timestamps known for
  it, its relations to other audit records, and any indicators that matched on it.
- Audit messages are the messages that THOR printed during the scan, in a less verbose form than the
  `THOR message` events.

Both audit records and audit messages are together called _audit entries_.

## Library Usage

### Parsing THOR Log Events

The `thorlog/parser` package parses events from all versions of the THOR log format.
It returns a `common.Event`, a version-agnostic representation of a THOR event.
Type-assert it to the version-specific type, e.g. `*thorlog.Assessment` for an assessment in version 3.

### Converting THOR Log Events to Text

`jsonlog.TextlogFormatter` converts assessments and messages to the human-readable text log format.
However, the text log format is not as rich as the JSON format and may not contain all fields.
When in doubt, use the JSON format for analysis.

### Generating Schemas

There are two schemas for the version 3 format:

- `thor-event.json` describes the events in a THOR log.
- `thor-audit-entry.json` describes the entries in the THOR audit trail.

The schemas are attached to each release.
They can also be generated using the `thorlog/jsonschema` package.
The generator takes the schema to generate as its only argument and must be run from within its directory:

```sh
cd thorlog/jsonschema
go run . log > thor-event.json
go run . audittrail > thor-audit-entry.json
```
