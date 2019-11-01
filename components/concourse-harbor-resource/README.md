# Harbor - Docker Image Resource

Includes Notary and Harbor specific setup.

**Important:** This is an extension of the [docker-image-resource]. The upstream documentation is still relevant to users of this resource.

## Source Configuration

In addition to the upstream docs, we're introducing the following:

* `harbor`: *Required.* An object containing further Harbor configuration.

  * `url`: *Required.* An URL pointing to the Harbor core UI. For instance: `https://core.harbor.com`

  * `public`: *Optional.* Allow the image to be publically accessible. Default: `true`

  * `enable_content_trust`: *Optional.* Allow content trust (Notary) signatures for that image. Default: `false`

  * `prevent_vul`: *Optional.* Prevent vulnerable images from being pushed. Default: `true`

  * `severity`: *Optional.* Set severity of vulnerability to be focusing on (`low`, `medium`, `high`). Default: `medium`

  * `auto_scan`: *Optional.* Scan the image for vulnerability. Default: `true`

* `notary`: *Required.* An object containing further Notary configuration.

  * `url`: *Required.* An URL pointing to the Notary API. For instance: `https://notary.harbor.com`

  * `passphrase`: *Required.* An object containing passphrases for different keys.

    * `root`: *Required.* A passphrase string for root key to be used when signing or acting on images.

    * `targets`: *Required.* A passphrase string for targets key to be used when signing or acting on images.
    
    * `snapshot`: *Required.* A passphrase string for snapshot key to be used when signing or acting on images.
    
    * `delegate`: *Required.* A passphrase string for delegate key to be used when signing or acting on images.

  * `keys`: *Required.* A PAM contents in a form of an escaped string combining number of keys for Notary.

  * `delegate_cert`: *Required.* A delegate certificate in a form of an escaped string of a key that will be used for signing images.

**Note:** The `repository` field from the upstream is now required to be formatted as such: `<registry_domain>/<project_name>/<image_name>`. For instance: `docker.io/library/nginx` or `core.harbor.com/team-awesome/app`.

**Note:** `username` and `password` fields are required for Harbor and Notary authentication.


[docker-image-resource]: https://github.com/concourse/docker-image-resource
