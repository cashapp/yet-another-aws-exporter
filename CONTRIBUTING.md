# Contributing

Thanks for considering contributing to this project, we're excited to see what changes you have in mind!

## Getting Started

Before submitting a PR, it might be good to do the following:

- Check [Issues](https://github.com/cashapp/yet-another-aws-exporter/issues) and [PRs](https://github.com/cashapp/yet-another-aws-exporter/pulls) to make sure the change you want hasn't already been documented/made
- Open an issue with the changes/new feature you're thinking about
- Take a peak [at the basic Scraper docs](./docs/scrapers.md) to understand the architecture of the exporter
- Fork the this repo into your own account
- Set up [your local development environment](./docs/development-guide.md)
- Submit a PR from your fork


## Releasing Updates

A new release of the exporter will automatically be created once a tag is pushed up to origin. Our GitHub actions Release workflow will do the following:

- Push an image to DockerHub with the tag version
- Compile binaries and add them to the release

THe easiest way to do this is to [use the `gh` CLI tool to push a release](https://cli.github.com/manual/gh_release_create). For example, if you wanted to release version `v1.2.3`, you would run:

```
gh release create v1.2.3
```
