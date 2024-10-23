# Setup
There are some things you will need to set up in order for the GitHub interaction to work.
## Server
`shortinette` is listening for activity on the participants' repositories using web hooks. For GitHub to be able to send `shortinette` payloads, you will need a **public IP address**.

For a common mortal, this would require paying for a server - however, as a 42 student, you have access to the [GitHub student pack](https://education.github.com/pack#offers). If you did not unlock it yet, [here](https://github-portal.42.fr) is the link - just log in with your 42 account and you're good to go!

The reason I'm telling you this is that the student pack offers 200$ of DigitalOcean credit for one year. DigitalOcean is a hosting service, which you will be able to use for hosting `shortinette`!

Create a DigitalOcean account (you can just sign up with GitHub). You will have to enter payment information, but you will not be charged _(as long as you follow my instructions)_.
### Creating the Server
Cloud servers are called **Droplets** on DigitalOcean. In case you are not familiar with cloud servers, they are just virtual machines hosted by someone else.

To host `shortinette`, the required specs depend heavily on how many participants your Short will have. To be safe, we recommend at least **2GB RAM** and **2 CPUs**. You can create your droplet [here](https://cloud.digitalocean.com/droplets/new?i=3bf27c&region=fra1&size=s-2vcpu-4gb-120gb-intel). 
* **Region**: Choose the one **closest to you**! The further your Droplet is, the more latency you will experience.
* **OS**: `Debian`, version `12 x64`. 
* **Droplet Type**: `Basic`.
* **CPU options**: `Regular`, then choose the **$18/mo** plan with `2GB RAM` and `2 CPUs`. 
* **SSH Key**: Add your SSH key. This will be needed to authenticate when connecting to the server.

You can now click on `Create Droplet`. You will be redirected to the list of your machines, and after about one minute, your Droplet will be ready.

Once you see the little green dot under your Droplet's icon, you are ready to keep going. Click on your droplet's name. This will open its dashboard, and right under the name, you will see its **ipv4 address**.

You can now connect to it via SSH. Open a shell, and run this command (replace `<ipv4>` with the address of your droplet):
```sh
$ ssh root@<ipv4>
```
That was easy, right?
### Installing Packages
You will need three packages:
* **Docker**, to deploy `shortinette`.
* **Tmux**, to share terminal sessions with co-organizers (believe me, you do not want to be the only one with access to the `shortinette` process).
* **SQLite3**, in case you ever need to manually access the database.

I made you an [installation script](scripts/server-setup.sh). Feel free to paste it into your Droplet's command line, it will install all required packages. It will also create a user named `Short`, and print its password to your console. Save this password! If you lose it, you will need your Born2beroot knowledge, _and you probably just followed a tutorial, so we definitely do not want that_.

You can now switch to the `Short` user:
```sh
su Short
```
From now on, you can just log in as `Short` instead of `root` (believe me, it's for your own good).

## GitHub Organisation
`shortinette` uses GitHub as an infrastructure. In order to group all of your participant's repositories, you will need to set up an organisation. You can create your organisation [here](https://github.com/organizations/plan). Choose the free plan, give it a name, and email, and you're good to go. This will make creating and managing repositories easier.

## Secrets
`shortinette` needs a peek into your secrets in order to authenticate with GitHub and manage repositories on your behalf.

* `ORGA_GITHUB`: The name of the GitHub organisation you created, so `shortinette` knows where to create the repositories.
* `TOKEN_GITHUB`: A personal access token with admin access to the `ORGA_GITHUB`. You can create it [here](https://github.com/organizations/Short-Test-Orga/settings/personal-access-tokens).
* `HOST_IP`: `http://<your-public-ip>`, where `<your-public-ip>` is the public ipv4 address of your droplet. Obviously, if you already have a fancy server with secure connection, use `https` - I guess you can figure that out yourself if you managed to set up SSL certificates.
* `WEBHOOK_PORT`: The port you want GitHub to send the web hook payloads to. Choose a port that is not already allocated! If you just created your Droplet following this documentation, just use `8080`.