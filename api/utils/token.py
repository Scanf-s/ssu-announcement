import hashlib
import hmac
import os


def generate_unsubscribe_token(email: str, category: str) -> str:
    return hmac.new(
        key=os.getenv("SECRET_KEY").encode('utf-8'),
        msg=f"{email}:{category}".encode('utf-8'),
        digestmod=hashlib.sha256
    ).hexdigest()
