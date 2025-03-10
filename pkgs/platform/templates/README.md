### Discord Function Template

Template to write a small discord function.

Commands are stored in `./dfaas.yaml`. This file must exist for templates to be deployable to the remote server. 

To test functions, you must have Docker installed locally. You can use `dfaas func build .` inside the current directory and the function will build a template locally. 

You can use `dfaas func invoke <command> --opts option1=value [...]` to test commands locally. They will only display the raw message and raw embed json. You can use a embed generator to preview what your embed will look like, such as https://message.style/. 