# How to Organize a Short on your Campus 
You might be thinking: "This is nice and all, but I don't want to spend a week piecing everything together to organize this on my campus.."

_And you would be right!_ There is a lot to organize, a lot to think about. Lucky for you, we've been through this so you don't have to! There are multiple things that go into the organization, so hang tight while I take you through it.

## Prerequisites
There are some things you will need to set up in order for the GitHub interaction to work.
### Server
`shortinette` is listening for activity on the participants' repositories using web hooks. For GitHub to be able to send `shortinette` payloads, you will need a **public IP address**.

For a common mortal, this would require paying for a server - however, as a 42 student, you have access to the [GitHub student pack](https://education.github.com/pack#offers). If you did not unlock it yet, [here](https://github-portal.42.fr) is the link - just log in with your 42 account and you're good to go!

The reason I'm telling you this is that the student pack offers 200$ of DigitalOcean credit for one year. DigitalOcean is a hosting service, which you will be able to use for hosting `shortinette`!

Create a DigitalOcean account (you can just sign up with GitHub). You will have to enter payment information, but you will not be charged _(as long as you follow my instructions)_.
#### Creating the Server
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
#### Installing Packages
You will need three packages:
* **Docker**, to deploy `shortinette`.
* **tmux**, to share terminal sessions with co-organizers (believe me, you do not want to be the only one with access to the `shortinette` process).
* **sqlite**, in case you ever need to manually access the database.



## Announcement
We recommend posting the announcement _at least_ one month before you plan to start the Short. Participation requires a big time investment, and you want people to plan ahead.

As an example, [[this]] is the announcement we made for our very first Rust Short. Feel free to use it and tweak it as you need. 

