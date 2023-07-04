# Ranger 1.6

[![Build Status](https://drone8.ranger.io/api/badges/ranger/ranger/status.svg?branch=master)](https://drone8.ranger.io/ranger/ranger)
[![Docker Pulls](https://img.shields.io/docker/pulls/ranger/server.svg)](https://store.docker.com/community/images/ranger/server)
[![Go Report Card](https://goreportcard.com/badge/github.com/ranger/ranger)](https://goreportcard.com/report/github.com/ranger/ranger)

Ranger is an open source project that provides a complete platform for operating Docker in production. It provides infrastructure services such as multi-host networking, global and local load balancing, and volume snapshots. It integrates native Docker management capabilities such as Docker Machine and Docker Swarm. It offers a rich user experience that enables devops admins to operate Docker in production at large scale.

## Latest Release

* Beta - v1.6.26 - `ranger/server:latest` - Read the full release [notes](https://github.com/ranger/ranger/releases/tag/v1.6.26).

* Stable - v1.6.26 - `ranger/server:stable` - Read the full release [notes](https://github.com/ranger/ranger/releases/tag/v1.6.26).

To get automated notifications of our latest release, you can watch the announcements category in our [forums](http://forums.ranger.com/c/announcements), or subscribe to the RSS feed `https://forums.ranger.com/c/announcements.rss`.

## Installation

Ranger is deployed as a set of Docker containers.  Running Ranger is as simple as launching two containers.  One container as the management server and another container on a node as an agent.  You can install the containers in following approaches.

* [Manually](#launching-management-server)
* [Terraform](https://github.com/ranger/terraform-modules)
* [Puppet](https://github.com/nickschuch/puppet-ranger) (Thanks @nickschuch)
* [Ansible](https://github.com/joshuacox/ansibleplaybook-ranger)
* [SaltStack](https://github.com/komljen/ranger-salt)

### Requirements

* [Supported Docker version](http://ranger.com/docs/ranger/v1.6/en/hosts/#supported-docker-versions)
* Any modern Linux distribution with a [supported Docker version](http://ranger.com/docs/ranger/v1.6/en/hosts/#supported-docker-versions). (Ubuntu, RHEL/CentOS 7 are more heavily tested.) Ranger also works with [RangerOS](https://github.com/ranger/os).
* RAM: 1GB+

### Launching Management Server

    docker run -d --restart=unless-stopped -p 8080:8080 ranger/server

The UI and API are available on the exposed port `8080`.

### Using Ranger

To learn more about using Ranger, please refer to our [Ranger Documentation](http://docs.ranger.com/).

## Source Code

This repo is a meta-repo used for packaging.  The source code for Ranger is in other repos in the ranger organization.  The majority of the code is in https://github.com/ranger/cattle.

## Support, Discussion, and Community
If you need any help with Ranger or RangerOS, please join us at either our [Ranger forums](http://forums.ranger.com/), [#ranger IRC channel](http://webchat.freenode.net/?channels=ranger) or [Slack](https://slack.ranger.io/) where most of our team hangs out at.

Please submit any **Ranger** bugs, issues, and feature requests to [ranger/ranger](//github.com/ranger/ranger/issues). 

Please submit any **RangerOS** bugs, issues, and feature requests to [ranger/os](//github.com/ranger/os/issues).

For security issues, please email security@ranger.com instead of posting a public issue in GitHub.  You may (but are not required to) use the GPG key located on [Keybase](https://keybase.io/ranger).

# License

Copyright (c) 2014-2018 [Ranger Labs, Inc.](http://ranger.com)

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

[http://www.apache.org/licenses/LICENSE-2.0](http://www.apache.org/licenses/LICENSE-2.0)

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
