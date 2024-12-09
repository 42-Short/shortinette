# How to Organize a Short on your Campus 
You might be thinking: "This is nice and all, but I don't want to spend a week
piecing everything together to organize this on my campus.."

_And you would be right!_ There is a lot to organize, a lot to think about.
Lucky for you, we've been through this so you don't have to! There are multiple
things that go into the organization, so hang tight while I take you through it.

## Table of Contents
1. [Server Set Up](#server-setup)
2. [Announce It](#announce-it)
3. [Discord Server Set Up](#discord-server-set-up)

## Server Setup
To enable GitHub interaction with `shortinette`, you'll need to set up a few
things. This guide will walk you through the process step-by-step.

### Server
`shortinette` uses web hooks to listen for activity on participants'
repositories. For this to work, you'll need a **public IP address**.

As a 42 student, you have access to the [GitHub student pack], which includes
$200 of DigitalOcean credit for one year. If you haven't unlocked it yet, visit
the [42 GitHub Portal] and log in with your 42 account.

#### Creating the Server
If you already have a cloud server, you can skip this section.

On DigitalOcean, cloud servers are called **Droplets**. Here's how to set one
up:

1. Create a DigitalOcean account (GitHub sign-up is available).
2. Go to [this link](https://cloud.digitalocean.com/droplets/new?i=3bf27c&region=fra1&size=s-2vcpu-4gb-120gb-intel) to create your droplet.
3. Choose the settings:
   * **Region**: Select the closest to you for minimal latency.
   * **OS**: `Debian`, version `12 x64`.
   * **Droplet Type**: `Basic`.
   * **CPU options**: `Regular`, then the **$18/mo** plan (2GB RAM, 2 CPUs).
   * **SSH Key**: Add your SSH key for authentication.

4. Click `Create Droplet`. After about a minute, your Droplet will be ready
   (indicated by a green dot).
5. Click on your Droplet's name to view its dashboard and find its
   **IPv4 address**.
6. Connect via SSH:
   ```sh
   $ ssh root@<ipv4>
   ```

#### Installing Packages
You'll need Docker, Tmux, and SQLite3. Use this
[installation script](scripts/server-setup.sh) to set everything up. It will
also create a user named `Short` - make sure to save the password it generates!
If you don't, you will need to pull up the Born2beRoot tutorial you blindly
followed to change the password ^^

Switch to the `Short` user:
```sh
su Short
```

Pro tip: Always log in as `Short` instead of `root` for better security
practices.

### GitHub Organisation
Create a [GitHub Organisation] to group your participants' repositories.
Choose the free plan and set it up with a name and email.

### Secrets
`shortinette` requires certain secrets for GitHub authentication and repository
management:

* `ORGA_GITHUB`: Your newly created GitHub organisation name.
* `TOKEN_GITHUB`: A personal access token with admin rights to `ORGA_GITHUB`.
  Create it [here](https://github.com/organizations/Short-Test-Orga/settings/personal-access-tokens).
* `HOST_IP`: `http://<your-public-ip>` (use your Droplet's IPv4 address).
* `WEBHOOK_PORT`: The port for GitHub web hook payloads. If you're using a fres
  Droplet, `8080` should work fine.
> [!NOTE]  
> If you have a server with SSL certificates, feel free to use `https` instead
> of `http` for `HOST_IP`.

### Configuring Participants
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
The `intra_login` variable will be used as a UID to build the names of the
participant's repos.

## Announce It
We recommend posting the announcement _at least_ one month before you plan to
start the Short. Participation requires a big time investment, and you want
people to plan ahead.


<details>

<summary>
Click here to see the announcement we made for our very first Rust Short.
Feel free to use it and tweak it as you need!
</summary>

```md
# It's finally time to get rusty @Students!

We are happy to announce that the Rust Short will take place from `<start date>` to `<end date>`.

##  🦀 What is Rust?

Did you know that about 70% of severe security bugs are caused by memory corruption? While this has a lot to do with skill-issues, (which we obviously do not have here at 42 hehe), languages that let you play loose with memory are definitely part of the problem. 

Rust is a modern systems programming language focused on safety and performance. It's designed to be a safer alternative to languages like C/C++. Rust's main goal is to prevent memory-related vulnerabilities while maintaining high performance. It introduces a way of thinking about memory management that prevents many common issues.
### Why Rust is Awesome
1. **No segfaults:** Rust's ownership system prevents common pointer issues.
2. **No data races:** Easier, safer concurrency.
3. **Safety != slow:** Performance on par with C/C++.
4. **Goodbye Makefiles:** Cargo simplifies project management.
5. **Growing Demand:** Adopted by major tech companies.
##  🚀 Recommended Prerequisite

We recommend having completed the `minishell` and `philosophers` circle before participating in the Rust Piscine. These projects provide insights into C concurrency & memory management, and will help you better appreciate Rust's take on solving these challenges. If you haven't done so yet, don't worry—you can still participate! Just be prepared for a steeper learning curve 💪

##  🏊‍♂️ Details
- **Duration:** 7 days (`<start date>` - `<end date>`)
- **Daily Structure:**
    - Each day, you will receive a new module containing 8 exercises
- **Time Commitment:** The first 5 exercises (50%) take about 6 hours/day, plan at least 10 for passing with 100%
- **On-campus participation required**
##  📝 How do I sign up?

Simply add your GitHub username to the thread below. Mark your calendars and let's get rusty! !🦀💻
```
</details>

We posted the announcement on the Discord server of our campus, and opened a
thread below it, where students could just type their GitHub usernames to sign
up (see [the doc](#setup-configuring-participants) on configuring participants).

## Discord Server Set Up
To streamline communication and avoid overwhelming the main student server,
we've created a dedicated Discord server for the Short. This server is designed
to improve the overall experience for both participants and the Short team.

### Features
1. Dual-Purpose Communication:
    * Facilitates interaction between Short participants
    * Allows the Short team to communicate with participants and monitor the
      event
2. Efficient Information Retrieval:
    * Searchable by exercise
    * Serves as a reliable source of truth
    * Encourages self-help among participants
3. Organized Structure:
    * Forum channels for specific topics:
        * One channel for technical issues
        * Separate channels for each module

### Server Template
We've created a [template](https://discord.new/YuCWVzpYbZns) based on our server
setup. However, please note:
* The template does not automatically include the community feature, you will
  need to manually activate it
* Forum channels must be created individually
### Setup Instructions
1. Use the provided server template as a starting point
2. Activate the community feature in your new server
3. Create forum channels:
    * Technical Issues
    * One for each module (e.g., Module 00: First Steps, ..., Module XX: XX)
4. Set appropriate permissions for the forum channels - the channels dedicated
   for the Short team, like `private-general` will already be set up
5. Invite participants and team members
### Bot
TODO: @ifaoji

### Tips
* Encourage participants to use the search function in the forum channels before
  asking questions
* Regularly pin important announcements or frequently asked questions
* Assign moderators to help manage the server and ensure a positive community
  experience

[GitHub student pack]: https://education.github.com/pack#offers
[42 GitHub Portal]: https://github-portal.42.fr
[GitHub Organisation]: https://github.com/organizations/plan
