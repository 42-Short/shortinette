use crate::client;

#[derive(clap::Parser, Debug)]
/// Easy to use CLI to interact with shortinette
pub struct Cli {
    /// Specifies the base URL which is used to send requests to Shortinette.
    #[arg(long, default_value = "http://127.0.0.1:3000")]
    pub base_url: reqwest::Url,

    #[command(subcommand)]
    pub command: Commands,
}

#[derive(clap::Subcommand, Debug)]
pub enum Commands {
    /// Sets the waiting time of a student to `0`.
    ResetWaittime,
}

impl Commands {
    pub fn exec(self, _client: &client::ApiClient) {
        match self {
            Self::ResetWaittime => {
                println!("Resetting Waittime...");
            }
        }
    }
}
