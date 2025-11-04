from typing import Optional, Any, Dict
from boto3.dynamodb.conditions import Key, Attr
import logging

class SubscriptionRepository:

    def __init__(self):
        self.logger = logging.getLogger("SubscriptionRepository")

    def get_subscriptions_by_email(self, table, email: Optional[str]) -> Optional[Dict]:
        response: Dict = table.query(
            KeyConditionExpression=Key('email').eq(email)
        )
        self.logger.debug(f"Query result: {response}")
        return response.get("Items")[0] if response.get("Items") else None

repo: SubscriptionRepository = SubscriptionRepository()