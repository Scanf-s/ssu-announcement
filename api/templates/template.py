import os
from string import Template

TEMPLATE_DIR: str = os.path.dirname(os.path.abspath(__file__))

def get_unsubscribe_error_page(message: str) -> str:
    template_path = os.path.join(TEMPLATE_DIR, "unsubscribe_error.html")
    with open(template_path, "r", encoding="utf-8") as f:
        html_template: Template = Template(f.read())

    html_content: str = html_template.substitute(message=message)
    return html_content

def get_unsubscribe_success_page(email: str, category: str) -> str:
    template_path = os.path.join(TEMPLATE_DIR, "unsubscribe_success.html")
    with open(template_path, "r", encoding="utf-8") as f:
        html_template: Template = Template(f.read())

    html_content: str = html_template.substitute(
        email=email,
        category=category
    )
    return html_content