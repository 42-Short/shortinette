pub struct ApiClient {
    client: reqwest::blocking::Client,
    base_url: reqwest::Url,
}

impl ApiClient {
    pub fn new(base_url: reqwest::Url) -> Self {
        let client = reqwest::blocking::Client::new();

        Self { client, base_url }
    }

    // pub async fn get(&self, endpoint: &str) -> Result<String, Error> {
    //     let url = self.base_url.join(endpoint)?;
    //     let response = self.client.get(url).send().await?;
    //     let body = response.text().await?;
    //     Ok(body)
    // }
}
