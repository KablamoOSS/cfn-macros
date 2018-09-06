# cfn-brainfuck

A brainfuck interpreter for CloudFormation.

## Usage

This transform looks for a "Brainfuck" key in the provided JSON document. The
value of this key should be a string containing a brainfuck program.

The brainfuck will be executed, and any data output (i.e. `.`) from the program
will be collected to produce a new CloudFormation template. The program should
be capable of producing a complete CloudFormation template in JSON format.

Note: any keys other than "Brainfuck" in the input template will be discarded.

