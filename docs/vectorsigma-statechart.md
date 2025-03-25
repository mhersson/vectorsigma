# The state chart of VectorSigma

This file contains the actual state chart of VectorSigma, and is used to
generate the current version of the application itself.

```plantuml
@startuml

skin rose

title VectorSigma

[*] --> Initializing

Initializing: do / Initialize
Initializing -[dotted]-> [*]: IsError
Initializing -[bold]-> LoadingInput

LoadingInput: do / LoadInput
LoadingInput -[dotted]-> [*]: IsError
LoadingInput --> ExtractingUML: IsMarkdown
LoadingInput -[bold]-> ParsingUML

ExtractingUML: do / ExtractUML
ExtractingUML -[dotted]-> [*]: IsError
ExtractingUML -[bold]-> ParsingUML

ParsingUML: do / ParseUML
ParsingUML -[dotted]-> [*]: IsError
ParsingUML -[bold]-> GeneratingStateMachine

GeneratingStateMachine: do / GenerateStateMachine
GeneratingStateMachine -[dotted]-> [*]: IsError
GeneratingStateMachine --> GeneratingModuleFiles: IsInitializingModule
GeneratingStateMachine -[bold]-> CreatingOutputFolder

GeneratingModuleFiles: do / GenerateModuleFiles
GeneratingModuleFiles -[dotted]-> [*]: IsError
GeneratingModuleFiles -[bold]-> CreatingInternalOutputFolder

CreatingInternalOutputFolder: do / CreateOutputFolder(internal)
CreatingInternalOutputFolder -[dotted]-> [*]: IsError
CreatingInternalOutputFolder -[bold]-> WritingGeneratedFiles

CreatingOutputFolder: do / CreateOutputFolder
CreatingOutputFolder -[dotted]-> [*]: IsError
CreatingOutputFolder --> MakingIncrementalUpdates: PackageExists
CreatingOutputFolder -[bold]-> WritingGeneratedFiles

MakingIncrementalUpdates: do / MakeIncrementalUpdates
MakingIncrementalUpdates -[dotted]-> [*]: IsError
MakingIncrementalUpdates -[bold]-> FilteringGeneratedFiles
note left of MakingIncrementalUpdates
  Compare the functions in
  actions and guards with
  the new generated code
end note

FilteringGeneratedFiles: do / FilterGeneratedFiles
FilteringGeneratedFiles -[dotted]-> [*]: IsError
FilteringGeneratedFiles -[bold]-> WritingGeneratedFiles
note left of FilteringGeneratedFiles
  If extendedstate.go exists,
  or if the actions and guards
  haven't changed filter them out
end note


WritingGeneratedFiles: do / WriteGeneratedFiles
WritingGeneratedFiles -[dotted]-> [*]: IsError
WritingGeneratedFiles -[bold]-> FormattingCode

FormattingCode: do / FormatCode
FormattingCode -[dotted]-> [*]: IsError
FormattingCode -[bold]-> [*]

@enduml

```
