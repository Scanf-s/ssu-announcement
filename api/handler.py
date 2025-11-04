from typing import Dict, Any, Tuple, Callable, List
from api.controller import subscribe, healthcheck
from api.utils.error import SubscriptionNotFound
from api.utils.response import success_response, error_response

import boto3
import logging

logger = logging.getLogger("[SSU_ANNOUNCEMENT_API_HANDLER]")

# Route mapping
ROUTES: Dict[Tuple, Callable] = {
    ("GET", "/subscribe"): subscribe.get_subscribes,
    ("GET", "/health"): healthcheck.health_check,
}

ALLOW_METHODS: List[str] = ["GET", "POST", "DELETE", "OPTIONS", "HEAD"]


def lambda_handler(event: Dict[str, Any], context: Any) -> Dict[str, Any]:
    try:
        # 1. Event에서 httpMethod, 호출 경로 추출
        http_method: str = event.get("httpMethod", "")
        path: str = event.get("path", "")
        if not http_method or not path or http_method not in ALLOW_METHODS:
            return error_response(400, "Invalid request")

        logger.debug(f"Request: {http_method} {path}")

        # 2. RouteKey 매핑
        route_key: Tuple[str, str] = (http_method, path)
        handler: Callable = ROUTES.get(route_key)
        if not handler:
            return error_response(404, "Route not found")

        # 3. DynamoDB 클라이언트 생성
        dynamodb: Any = boto3.resource('dynamodb')

        # 4. Handler 실행 및 응답 반환
        result: Any = handler(event, dynamodb)
        return success_response(result)

    except ValueError as e:
        logger.error(f"Validation error: {str(e)}")
        return error_response(400, str(e))

    except SubscriptionNotFound as e:
        logger.error(f"Subscription not found: {str(e)}")
        return error_response(404, str(e))

    except Exception as e:
        logger.error(f"Internal error: {str(e)}")
        return error_response(500, "Internal server error")
