import logging
from typing import Any, Dict, Optional, List

from boto3.dynamodb.conditions import Attr, Key


class SubscriptionRepository:

    def __init__(self):
        self.logger = logging.getLogger("SubscriptionRepository")

    def get_subscriptions_by_email(self, table, email: Optional[str]) -> Optional[Dict]:
        response: Dict = table.query(KeyConditionExpression=Key("Email").eq(email))
        self.logger.debug(f"Query result: {response}")
        return response.get("Items")[0] if response.get("Items") else None

    def add_subscription(self, table: Any, email: Optional[str], category: Optional[str]):
        if not email or not category:
            raise ValueError("Email and Category are required")
        self.logger.debug(f"Add subscription: {email} with {category}")
        table.put_item(
            Item={"Email": email, "Category": category}
        )

    def delete_subscription(self, table, email: Optional[str], category: Optional[str]):
        if not email or not category:
            raise ValueError("Email and Category are required")
        self.logger.debug(f"Delete subscription: {email} with {category}")

        if category == "all":
            response: Dict = table.query(
                KeyConditionExpression=Key("Email").eq(email)
            )
            items: List = response.get("Items", [])
            for item in items:
                table.delete_item(
                    Key={"Email": item.get("Email"), "Category": item.get("Category")}
                )
        else:
            table.delete_item(
                Key={"Email": email, "Category": category}
            )


repo: SubscriptionRepository = SubscriptionRepository()
