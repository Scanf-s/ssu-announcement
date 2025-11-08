import sys
from pathlib import Path
sys.path.insert(0, str(Path(__file__).parent.parent.parent))

from typing import Dict, Any
import os
from dotenv import load_dotenv
from api.app import lambda_handler

load_dotenv(".env")


class ControllerTests:

    def __init__(self):
        self.mock_get_event: Dict[str, Any] = {
            "httpMethod": "GET",
            "path": "/subscribe",
            "queryStringParameters": {
                "email": os.getenv("TEST_EMAIL")
            }
        }
        self.mock_post_event: Dict[str, Any] = {
            "httpMethod": "POST",
            "path": "/subscribe",
            "queryStringParameters": {
                "email": os.getenv("TEST_EMAIL"),
                "category": "test"
            }
        }
        self.mock_delete_all_event: Dict[str, Any] = {
            "httpMethod": "DELETE",
            "path": "/subscribe",
            "queryStringParameters": {
                "email": os.getenv("TEST_EMAIL"),
                "category": "all"
            }
        }
        self.mock_delete_event: Dict[str, Any] = {
            "httpMethod": "DELETE",
            "path": "/subscribe",
            "queryStringParameters": {
                "email": os.getenv("TEST_EMAIL"),
                "category": "test"
            }
        }

    def test_get_subscriptions(self):
        result: Any = lambda_handler(event=self.mock_get_event, context=None)
        assert result is not None
        assert result.get("statusCode") == 200
        assert result.get("headers").get("Content-Type") == "application/json"
        assert result.get("body") is not None
        assert "Email" in result.get("body")
        assert "Category" in result.get("body")
        assert os.getenv("TEST_EMAIL") in result.get("body")

    def test_delete_subscription(self):
        result: Any = lambda_handler(event=self.mock_delete_event, context=None)
        assert result is not None
        assert result.get("statusCode") == 200

    def test_add_subscription(self):
        result: Any = lambda_handler(event=self.mock_post_event, context=None)
        assert result is not None
        assert result.get("statusCode") == 201
        assert result.get("headers").get("Content-Type") == "application/json"
        assert result.get("body") is not None

        # Rollback
        lambda_handler(event=self.mock_delete_event, context=None)

if __name__ == "__main__":
    tests = ControllerTests()
    try:
        tests.test_get_subscriptions()
        tests.test_delete_subscription()
        tests.test_add_subscription()
        print("passed")
    except AssertionError as e:
        print(f"test failed: {e}")
    except Exception as e:
        print(f"test error: {e}")