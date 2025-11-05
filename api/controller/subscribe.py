"""
Subscribe endpoints handler
"""

import os
from typing import Any, Dict, Optional

from api.repository.subscription_repository import repo
from api.utils.error import SubscriptionNotFound


def get_subscribes(event: Dict[str, Any], db_session: Any) -> Dict[str, Any]:
    # 1. Query parameter에서 Email 추출
    query_params: Optional[dict[str, Any]] = event.get("queryStringParameters")
    email = query_params.get("email", "") if query_params else None
    if not email:
        raise ValueError("Email parameter is required")

    # 2. 테이블 설정
    table: Any = db_session.Table(os.getenv("SUBSCRIPTION_TABLE"))

    # 3. 사용자의 구독 정보 가져오기
    subscribes: Optional[Dict] = repo.get_subscriptions_by_email(table, email)
    if not subscribes:
        raise SubscriptionNotFound(email)

    return {"subscribes": subscribes}
