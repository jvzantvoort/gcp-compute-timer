# gcp-compute-timer

Monitor the uptime of gcp computer instances and act accordingly.

# Case

At the moment I have a number of images running in GCP. However
my budget is limited and I sometimes forget to turn of unused
images (scatterbrain). On the otherhand there are images that run
for extended amounts of time (hosting a site or something).

So what I need is a way to:

* get state and uptime from images
* indicate the maximum age of an image
* react to the image exceeding its limit

Google has an excelent API but I don't want to schlep around the SDK
all the time therefor I'm writing this in golang (although I am barely
past the "hello world" phase).

# Configuration

Create a configuration file called
``~/.config/gcp-compute-timer.yml`` with the following content:

```
gcp:
  project: project_name
  zone: europe-west4-a
  bucket: bucket_name
```

# GCP

## API instructions

### BEFORE RUNNING:

1. If not already done, enable the Compute Engine API
   and check the quota for your project at
   https://console.developers.google.com/apis/api/compute
2. This sample uses Application Default Credentials for authentication.
   If not already done, install the gcloud CLI from
   https://cloud.google.com/sdk/ and run
   `gcloud beta auth application-default login`.
   For more information, see
   https://developers.google.com/identity/protocols/application-default-credentials
3. Install and update the Go dependencies by running `go get -u` in the
   project directory.
