import requests
import json
import os
from dotenv import load_dotenv
import requests
import logging


logging.basicConfig(level=logging.INFO)
logger = logging.getLogger(__name__)

load_dotenv()

def delete_repo(repo_id: str) -> None:
    url = f"https://api.github.com/repos/{os.getenv('GITHUB_ORGANISATION')}/{repo_id}"
    headers = {
        "Authorization": f"Bearer {os.getenv('GITHUB_TOKEN')}",
        "Accept": "application/vnd.github.v3+json",
    }

    response = requests.delete(url, headers=headers)

    if response.status_code == 204:
        logger.info(f"successfully deleted repo {repo_id}")
        return True
    else:
        logger.error(f"failed to delete repo {repo_id}: {response.status_code} {response.text}")
        return False

with open("shortconfig/shortconfig.json") as f:
    short_config = json.load(f)
    for participant in short_config["participants"]:
        delete_repo(f"{participant['intra_login']}-00")
