[![Build Status](https://travis-ci.org/glesica/farnsworth.svg?branch=master)](https://travis-ci.org/glesica/farnsworth)

# Farnsworth

Farnsworth is a tool to assist in the creation and evaluation of programming
assignments.

To create an assignment, the instructor first implements the project and then
annotates the source code to mark sections that should be hidden from students.

When an archive of the project is created, these sections will be removed. This
might include implementation or selected tests.

Then, once students have completed the assignment, some of those sections (for
example, extra tests) are automatically merged back into each project for
evaluation.

This is a work-in-progress. Right now it supports Java and Go projects.
I intend to add Python and possibly C as well. Adding a project type is pretty
easy to do, particularly right now since Farnsworth still doesn't do a whole
lot.

