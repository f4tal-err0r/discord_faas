## Discord FAAS

This is a project that allows small arbitrary functions to run and output as a Discord bot using commands. This is to make small, quick functions for ease of extensibility.

### Architecture (WIP-ish)

- Serverless Event Engine
    - Knative: Easy, can probably make an cloudevent knative template for controller interactions 
    - OpenFAAS:  Exploring this since it allows for more working with the higher level templates
- SCM:
    - Github repo + build integration
