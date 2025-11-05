import sys
from pathlib import Path

sys.path.insert(0, str(Path(__file__).parent.parent.parent))

import os
from typing import Any

import boto3
from dotenv import load_dotenv

from api.repository.subscription_repository import repo

load_dotenv("../.env")


class RepositoryTests:

    def __init__(self):
        self.table: Any = boto3.resource(
            "dynamodb", region_name="ap-northeast-2"
        ).Table(os.getenv("SUBSCRIPTION_TABLE"))
        self.test_email: str = os.getenv("TEST_EMAIL")

    def test_get_subscriptions_by_email(self):
        response = repo.get_subscriptions_by_email(
            table=self.table, email=self.test_email
        )
        assert response is not None
        assert response.get("Email") is not None
        assert response.get("Category") is not None
        assert response.get("Email") == self.test_email


test = RepositoryTests()
test.test_get_subscriptions_by_email()
