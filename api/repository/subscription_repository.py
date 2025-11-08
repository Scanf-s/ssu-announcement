import logging
from typing import Any, Dict, Optional, List, Tuple

from boto3.dynamodb.conditions import Attr, Key


class SubscriptionRepository:

    def __init__(self):
        self.logger = logging.getLogger("SubscriptionRepository")
        self.logger.setLevel(logging.INFO)

    def get_subscriptions_by_email(self, table, email: Optional[str]) -> Optional[List[Dict]]:
        response: Dict = table.query(KeyConditionExpression=Key("Email").eq(email))
        self.logger.info(f"Query result: {response}")
        return response.get("Items") if response.get("Items") else None

    def add_subscription(self, table: Any, email: Optional[str], category: Optional[str], token: Optional[str]):
        if not email or not category or not token:
            raise ValueError("Email and Category, Token are required")
        self.logger.info(f"Add subscription {token[:6]}... // {email} with {category}")
        table.put_item(
            Item={"Email": email, "Category": category, "UnsubscribeToken": token}
        )

    def delete_subscription(self, table: Any, email: Optional[str], category: Optional[str]):
        if not email or not category:
            raise ValueError("Email and Category are required")
        self.logger.info(f"Delete subscription: {email} with {category}")

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

    def verify_unsubscribe_token(self, table: Any, token: str) -> Tuple[bool, Optional[dict]]:
        try:
            self.logger.info(f"Try to fetch user with token: {token}")
            response: Dict = table.query(
                IndexName='UnsubscribeTokenIndex',
                KeyConditionExpression=Key("UnsubscribeToken").eq(token)
            )
            self.logger.info(f"Query response: {response}")
            return (True, response.get("Items")[0]) if response.get("Items") else (False, None)
        except Exception as e:
            self.logger.error(f"Error occurred: {str(e)}")
            return False, None

repo: SubscriptionRepository = SubscriptionRepository()
