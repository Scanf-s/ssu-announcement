import os
from typing import Any, Dict, Optional, List

from repository.subscription_repository import repo
from utils.error import SubscriptionNotFound
from utils.token import generate_unsubscribe_token
from utils.response import html_response
from templates.template import get_unsubscribe_error_page, get_unsubscribe_success_page

import logging
logger = logging.getLogger("[SSU_ANNOUNCEMENT_API_CONTROLLER]")
logger.setLevel(logging.INFO)

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

def subscribe(event: Dict[str, Any], db_session: Any) -> None:
    # Query parameter에서 Email, Category 추출
    email: str = _get_email_from_event(event)
    category: str = _get_category_from_event(event)

    # Set table
    table: Any = db_session.Table(os.getenv("SUBSCRIPTION_TABLE"))

    # Generate subscription unique token for unsubscribe
    token: str = generate_unsubscribe_token(email=email, category=category)

    # Add subscription with unique token
    repo.add_subscription(table=table, email=email, category=category, token=token)

def unsubscribe(event: Dict[str, Any], db_session: Any) -> Dict[str, Any]:
    try:
        # Get query parameters
        token: str = _get_token_from_event(event)

        # Validate parameters
        if not token:
            return html_response(get_unsubscribe_error_page("잘못된 요청입니다. 이메일에서 제공된 링크를 사용해주세요."), 400)

        # Set table
        table: Any = db_session.Table(os.getenv("SUBSCRIPTION_TABLE"))

        # Verify token
        # bool, dict[str, Any]
        success, data = repo.verify_unsubscribe_token(table, token)
        logger.info(f"Verify token: {success}, {data}")
        if not success:
            return html_response(get_unsubscribe_error_page("유효하지 않은 요청입니다. 토큰이 만료되었거나 올바르지 않습니다."), 400)

        # Delete subscription
        email: str = data.get("Email")
        category: str = data.get("Category")
        table: Any = db_session.Table(os.getenv("SUBSCRIPTION_TABLE"))
        repo.delete_subscription(table=table, email=email, category=category)

        # Return success HTML
        return html_response(get_unsubscribe_success_page(email, category), 200)

    except Exception as e:
        return html_response(get_unsubscribe_error_page(f"오류가 발생했습니다"), 500)

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

def _get_token_from_event(event: Dict[str, Any]) -> str:
    token: Optional[str] = event.get("queryStringParameters", {}).get("token", "")
    if not token:
        raise ValueError("Token parameter is required")
    return token