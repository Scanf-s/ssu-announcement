"""
Subscribe endpoints handler
"""

import os
from typing import Any, Dict, Optional, List

from repository.subscription_repository import repo
from utils.error import SubscriptionNotFound

CATEGORIES: List[str] = [
    "ssu_path",
    "학사",
    "국제교류",
    "봉사",
    "비교과·행사",
    "장학",
    "채용",
    "all",
]

def get_subscribes(event: Dict[str, Any], db_session: Any) -> Dict[str, Any]:
    # 1. Query parameter에서 Email 추출
    email: str = _get_email_from_event(event)

    # 2. 테이블 설정
    table: Any = db_session.Table(os.getenv("SUBSCRIPTION_TABLE"))

    # 3. 사용자의 구독 정보 가져오기
    subscribes: Optional[Dict] = repo.get_subscriptions_by_email(table, email)
    if not subscribes:
        raise SubscriptionNotFound(email)

    return {"subscribes": subscribes}

def add_subscribe(event: Dict[str, Any], db_session: Any) -> None:
    # Query parameter에서 Email, Category 추출
    email: str = _get_email_from_event(event)
    category: str = _get_category_from_event(event)

    # Set table
    table: Any = db_session.Table(os.getenv("SUBSCRIPTION_TABLE"))

    # Add subscription
    repo.add_subscription(table=table, email=email, category=category)

def delete_subscribe(event: Dict[str, Any], db_session: Any) -> None:
    # Query parameter에서 Email, Category 추출
    email: str = _get_email_from_event(event)
    category: str = _get_category_from_event(event)

    # Set table
    table: Any = db_session.Table(os.getenv("SUBSCRIPTION_TABLE"))

    # Delete subscription
    repo.delete_subscription(table=table, email=email, category=category)

def _get_email_from_event(event: Dict[str, Any]) -> str:
    email: Optional[str] = event.get("queryStringParameters", {}).get("email", "")
    if not email:
        raise ValueError("Email parameter is required")
    return email

def _get_category_from_event(event: Dict[str, Any]) -> str:
    category: Optional[str] = event.get("queryStringParameters", {}).get("category", "")
    if not category:
        raise ValueError("Category parameter is required")

    if category not in CATEGORIES:
        raise ValueError("Invalid category")

    return category