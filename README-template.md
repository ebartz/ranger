# Ranger

*This file is auto-generated from README-template.md, please make any changes there.*

[![Build Status](https://drone-publish.ranger.io/api/badges/ranger/ranger/status.svg?branch=release/v2.6)](https://drone-publish.ranger.io/ranger/ranger)
[![Docker Pulls](https://img.shields.io/docker/pulls/ranger/ranger.svg)](https://store.docker.com/community/images/ranger/ranger)
[![Go Report Card](https://goreportcard.com/badge/github.com/ranger/ranger)](https://goreportcard.com/report/github.com/ranger/ranger)

Ranger is an open source container management platform built for organizations that deploy containers in production. Ranger makes it easy to run Kubernetes everywhere, meet IT requirements, and empower DevOps teams.

> Looking for Ranger 1.6.x info? [Click here](https://github.com/ranger/ranger/blob/master/README_1_6.md)

## Latest Release

To get automated notifications of our latest release, you can watch the announcements category in our [forums](http://forums.ranger.com/c/announcements), or subscribe to the RSS feed `https://forums.ranger.com/c/announcements.rss`.

## Quick Start

    sudo docker run -d --restart=unless-stopped -p 80:80 -p 443:443 --privileged ranger/ranger

Open your browser to https://localhost

## Installation

See [Installing/Upgrading Ranger](https://ranger.com/docs/ranger/v2.6/en/installation/) for all installation options.

### Minimum Requirements

* Operating Systems
  * Please see [Support Matrix](https://ranger.com/support-matrix/) for specific OS versions for each Ranger version. Note that the link will default to the support matrix for the latest version of Ranger. Use the left navigation menu to select a different Ranger version. 
* Hardware & Software
  * Please see [Installation Requirements](https://ranger.com/docs/ranger/v2.6/en/installation/requirements/) for hardware and software requirements.

### Using Ranger

To learn more about using Ranger, please refer to our [Ranger Documentation](https://ranger.com/docs/ranger/v2.6/en/).

## Source Code

This repo is a meta-repo used for packaging and contains the majority of Ranger codebase. For other Ranger projects and modules, [see go.mod](https://github.com/ranger/ranger/blob/release/v2.6/go.mod) for the full list.

Ranger also includes other open source libraries and projects, [see go.mod](https://github.com/ranger/ranger/blob/release/v2.6/go.mod) for the full list.

## Support, Discussion, and Community
If you need any help with Ranger, please join us at either our [Ranger forums](http://forums.ranger.com/) or [Slack](https://slack.ranger.io/) where most of our team hangs out at.

Please submit any Ranger bugs, issues, and feature requests to [ranger/ranger](https://github.com/ranger/ranger/issues).

For security issues, please first check our [security policy](SECURITY.md) and email security-ranger@suse.com instead of posting a public issue in GitHub.  You may (but are not required to) use the GPG key located on [Keybase](https://keybase.io/ranger).

# License

Copyright (c) 2014-CURRENTYEAR [Ranger Labs, Inc.](http://ranger.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
