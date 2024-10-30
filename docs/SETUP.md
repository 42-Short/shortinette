# Setup
To enable GitHub interaction with `shortinette`, you'll need to set up a few things. This guide will walk you through the process step-by-step.

## Server
`shortinette` uses web hooks to listen for activity on participants' repositories. For this to work, you'll need a **public IP address**.

As a 42 student, you have access to the [GitHub student pack](https://education.github.com/pack#offers), which includes $200 of DigitalOcean credit for one year. If you haven't unlocked it yet, visit [this link](https://github-portal.42.fr) and log in with your 42 account.

### Creating the Server
On DigitalOcean, cloud servers are called **Droplets**. Here's how to set one up:

1. Create a DigitalOcean account (GitHub sign-up is available).
2. Go to [this link](https://cloud.digitalocean.com/droplets/new?i=3bf27c&region=fra1&size=s-2vcpu-4gb-120gb-intel) to create your droplet.
3. Choose the settings:
   * **Region**: Select the closest to you for minimal latency.
   * **OS**: `Debian`, version `12 x64`.
   * **Droplet Type**: `Basic`.
   * **CPU options**: `Regular`, then the **$18/mo** plan (2GB RAM, 2 CPUs).
   * **SSH Key**: Add your SSH key for authentication.

4. Click `Create Droplet`. After about a minute, your Droplet will be ready (indicated by a green dot).
5. Click on your Droplet's name to view its dashboard and find its **IPv4 address**.
6. Connect via SSH:
   ```sh
   $ ssh root@<ipv4>
   ```

### Installing Packages
You'll need Docker, Tmux, and SQLite3. Use this [installation script](scripts/server-setup.sh) to set everything up. It will also create a user named `Short` - make sure to save the password it generates! If you don't, you will need to pull up the Born2beroot tutorial you blindly followed to change the password ^^

Switch to the `Short` user:
```sh
su Short
```

Pro tip: Always log in as `Short` instead of `root` for better security practices.

## GitHub Organisation
Create a GitHub organisation [here](https://github.com/organizations/plan) to group your participants' repositories. Choose the free plan and set it up with a name and email.

## Secrets
`shortinette` requires certain secrets for GitHub authentication and repository management:

* `ORGA_GITHUB`: Your newly created GitHub organisation name.
* `TOKEN_GITHUB`: A personal access token with admin rights to `ORGA_GITHUB`. Create it [here](https://github.com/organizations/Short-Test-Orga/settings/personal-access-tokens).
* `HOST_IP`: `http://<your-public-ip>` (use your Droplet's IPv4 address).
* `WEBHOOK_PORT`: The port for GitHub web hook payloads. If you're using a fresh Droplet, `8080` should work fine.
Note: If you have a server with SSL certificates, feel free to use `https` instead of `http` for `HOST_IP`.

## Configuring Participants
```
TODO: Finish when actual configuration logic is ready.
```
The Short takes a configuration in `json` format:
```json
{
  "participants": [
    {
      "github_username": "",
      "intra_login": ""
    }
  ]
}
```
The `intra_login` variable will be used as a UID to build the names of the participant's repos.

