use clap::Parser;

mod cli;
mod client;

fn main() {
    let cmd = cli::Cli::parse();
    let client = client::ApiClient::new(cmd.base_url);

    cmd.command.exec(&client);
}
