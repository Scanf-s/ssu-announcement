import logging
from typing import Any, Dict, Optional

from boto3.dynamodb.conditions import Attr, Key


class SubscriptionRepository:

    def __init__(self):
        self.logger = logging.getLogger("SubscriptionRepository")

    def get_subscriptions_by_email(self, table, email: Optional[str]) -> Optional[Dict]:
        response: Dict = table.query(KeyConditionExpression=Key("Email").eq(email))
        self.logger.debug(f"Query result: {response}")
        return response.get("Items")[0] if response.get("Items") else None


repo: SubscriptionRepository = SubscriptionRepository()
