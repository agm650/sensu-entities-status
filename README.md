# entities status

## Table of Contents

- [entities status](#entities-status)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Files](#files)
  - [Usage examples](#usage-examples)
  - [Configuration](#configuration)
    - [Asset registration](#asset-registration)
    - [Resource definition](#resource-definition)
  - [Installation from source](#installation-from-source)
  - [Additional notes](#additional-notes)
  - [Contributing](#contributing)

## Overview

The entities status is a [Sensu CLI][6] tool (`sensuctl`) that ...

## Files

## Usage examples

## Configuration

### Asset registration

[Sensu Assets][10] are the best way to make use of this plugin. If you're not using an asset, please
consider doing so! If you're using sensuctl 5.13 with Sensu Backend 5.13 or later, you can use the
following command to add the asset:

```sh
sensuctl asset add agm650/entities-status
```

If you're using an earlier version of sensuctl, you can find the asset on the [Bonsai Asset Index][https://bonsai.sensu.io/assets/agm650/entities-status].

### Resource definition

## Installation from source

The preferred way of installing and deploying this plugin is to use it as an Asset. If you would
like to compile and install the plugin from source or contribute to it, download the latest version
or create an executable script from this source.

From the local path of the entities-status repository:

```sh
go build
```

## Additional notes

## Contributing

For more information about contributing to this plugin, see [Contributing][1].

[1]: https://github.com/sensu/sensu-go/blob/master/CONTRIBUTING.md
[2]: https://github.com/sensu/sensu-plugin-sdk
[3]: https://github.com/sensu-plugins/community/blob/master/PLUGIN_STYLEGUIDE.md
[4]: https://github.com/sensu/sensuctl-plugin-template/blob/master/.github/workflows/release.yml
[6]: https://docs.sensu.io/sensu-go/latest/sensuctl/reference/
[7]: https://github.com/sensu/sensuctl-plugin-template/blob/master/main.go
[8]: https://bonsai.sensu.io/
[9]: https://github.com/sensu/sensu-plugin-tool
[10]: https://docs.sensu.io/sensu-go/latest/reference/assets/
