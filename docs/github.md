
# Github Action + Github Page

You can use GitHub Actions to periodically generate an updated page, similar to a cron job, and host it using GitHub Pages. The advantage of this approach is that it’s free (as in free beer) and serverless. Additionally, you can directly update your feed list from the GitHub interface.

To use this method you will need to create a github repository with:

- GitHub Pages enabled
- A [repository variable](https://docs.github.com/en/actions/writing-workflows/choosing-what-your-workflow-does/store-information-in-variables) named FEEDS that contains the list of feeds you want to aggregate
- The following GitHub Actions workflow file located at `.github/workflows/daily.yml`:

```yaml
name: update-demo-daily
on:
  schedule:
    - cron: "0 6 * * *" # every day at 6am
  push:
    tags:
      - "refresh-github-page" # to manually trigger a refresh

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - name: prepare-page
        env:
          FEEDS: ${{vars.FEEDS}}
        run: |
          wget -q -O tinyfeed https://github.com/TheBigRoomXXL/tinyfeed/releases/latest/download/tinyfeed_linux_arm64
          chmod +x tinyfeed
          mkdir www
          echo $FEEDS | ./tinyfeed > ./www/index.html
          cp www/index.html www/404.html 
        # 404.html allows every path to be served by the index.html

      - name: upload-page
        uses: actions/upload-pages-artifact # add @ to specify a version
        with:
          path: www/

  deploy:
    needs: build
    permissions:
      pages: write # to deploy to Pages
      id-token: write # to verify the deployment originates from an appropriate source
    environment:
      name: github-pages
      url: ${{ steps.deployment.outputs.page_url }}
    runs-on: ubuntu-latest
    steps:
      - name: deploy-to-github-pages
        id: deployment
        uses: actions/deploy-pages # add @ to specify a version
```
