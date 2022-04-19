# `lessonplan`

This repository contains the sources for my `lessonplan` tool. This tool builds
lesson plans for my [Modelica
Playground](https://playground.modelica.universty).

The `sample` directory includes a sample lesson plan specification. The
structure of the lesson plan is as follows:

## Installing

If you have `go` installed, you can install the `lessonplan` tool with just:

```sh
$ go install github.com/mtiller/lessonplan@latest
```

If you don't have `go` installed, consult the [Releases](#) area for binary
releases for your platform.

## Creating a Lesson Plan

First, create a directory somewhere, _e.g.,_

```sh
$ mkdir mylesson
$ cd mylesson
```

### Master Index

Next, create an `index.json` file in that directory. This file should include a
`"title"` properties as well as a `"contents"` property that contains a list of
the directory names that each lesson can be found in _in the order you want them
to appear in the app_.

For example:

```json
{
  "title": "Sample Lessonplan",
  "contents": ["lesson1", "lesson2"]
}
```

### One Lesson per Directory

Now, in each directory, create another `index.json` file. However, the
`index.json` file for each lesson should contain only a `"title"`, _e.g._,

```json
{
  "title": "Lesson 1"
}
```

### Model

In additional the the `index.json` file, you'll need to create a `model.mo` file
that contains the Modelica code associated with this lesson.

### Explanation

In general, the model you are presenting the user will require some explanation.
This should be placed in a file called `explanation.md` that sits in the same
directory as `model.mo`.

This file is nominally in Markdown syntax. But there are numerous extensions implemented to the Markdown rendering in Modelica Playground. For more information, consult this [How To](https://playground.modelica.university/?toc=howto.json)

### Report

By default, the [Modelica Playground](https://playground.modelica.university)
will show a table of all constants as well as plots of all time varying signals.
However, as explained
[here](https://playground.modelica.university/?toc=howto.json), the simulation
report can be customized in many, many ways. Such customizations should be
placed in the `report.md` file right alongside `model.mo` and `explanation.md`.

### Template Premable

Complex simulation reports may wish to include templates. Because the Modelica
Playground application allows users to modify models and simulation reports
**and** bundles those modifications into special URLs, we want to avoid
including any predefined templates in the URL. This is because if included it
would contribute the overall URL length. But since these templating preambles
are not meant to be modified and are directly associated with specific lessons,
we have carved them out into a separate file to help keep the `report.mod` file
just to what the user might wish to modify.

If you wish to define templating macros for use in your `report.md` file, add
them to a file named `preamble.md` and place them alongside `report.mod` and the
other files.

### Complete Directory Structure

When done, your directory structure should look like this:

```
<dir>/
  index.json
  <each lesson directory>/
    index.json
    model.mo
    explanation.md // optional
    report.md // optional
    preamble.md // optional
```

## Generating a Lesson Plan

To generate the lesson plan, just run the following command:

```sh
$ lessonplan --dir <lessonplan directory> --output mylesson.json
```

## Hosting your Lesson Plan

In this repository, we've deliberately included in the `git` repository the _output file_ for our sample. Normally, one should avoid including compilation artifacts in a version control repository. But by including it here, we can then leverage Github as a hosting platform for our lesson.

To generate the lesson, we simple run:

```sh
$ lessonplan --dir sample -o sample.json
```

From this repository's root directory and then commit the resulting `sample.json` file in Github. The result is that we can now access our lesson plan file from:

```
https://raw.githubusercontent.com/mtiller/lessonplan/master/sample.json
```

This, in turn, allows us to reference our lesson plan when accessing the Modelica Playground. To access a lesson plan from the Modelica Playground, simply add the `?toc=<url>` query string. For example, our sample lesson plan can be referenced directory at:

[https://playground.modelica.university?toc=https://raw.githubusercontent.com/mtiller/lessonplan/master/sample.json](https://playground.modelica.university?toc=https://raw.githubusercontent.com/mtiller/lessonplan/master/sample.json)

## Building from Source

To build from source, you need to have `go` installed. Once that is done, building is as simple as:

```sh
$ go build .
```
