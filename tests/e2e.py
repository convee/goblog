#!/usr/bin/env python3
import argparse
import time
from urllib.parse import urljoin, urlparse

from playwright.sync_api import expect, sync_playwright


def url(base_url: str, path: str) -> str:
    return urljoin(base_url.rstrip("/") + "/", path.lstrip("/"))


def set_markdown(page, value: str) -> None:
    page.wait_for_function("window.goblogEditor && window.goblogEditor.state.loaded")
    page.evaluate("value => window.goblogEditor.setMarkdown(value)", value)


def wait_path(page, expected_path: str) -> None:
    page.wait_for_url(
        lambda current: urlparse(current).path.rstrip("/")
        == expected_path.rstrip("/")
    )


def submit_delete(page, row, expected_path: str) -> None:
    page.once("dialog", lambda dialog: dialog.accept())
    row.get_by_role("button", name="删除").click()
    wait_path(page, expected_path)


def run(base_url: str, email: str, password: str) -> None:
    suffix = str(int(time.time() * 1000))[-8:]
    category = f"E2E分类{suffix}"
    page_ident = f"e2e-{suffix}"
    page_title = f"E2E页面{suffix}"
    post_title = f"E2E文章{suffix}"
    tag = f"e2e-tag-{suffix}"
    browser_errors: list[str] = []

    with sync_playwright() as playwright:
        browser = playwright.chromium.launch(headless=True)
        context = browser.new_context()
        page = context.new_page()
        page.on(
            "console",
            lambda message: browser_errors.append(f"console: {message.text}")
            if message.type == "error"
            else None,
        )
        page.on("pageerror", lambda error: browser_errors.append(f"page: {error}"))

        for endpoint in ("/healthz", "/readyz"):
            response = page.request.get(url(base_url, endpoint))
            assert response.status == 204, (endpoint, response.status, response.text())

        response = page.request.get(url(base_url, "/"))
        assert response.status == 200
        assert "script-src 'self'" in response.headers["content-security-policy"]
        assert "unsafe-inline" not in response.headers["content-security-policy"].split("script-src", 1)[1].split(";", 1)[0]
        assert response.headers["x-content-type-options"] == "nosniff"

        page.goto(url(base_url, "/admin/"))
        wait_path(page, "/login")
        csrf_token = page.locator('input[name="_csrf"]').input_value()
        blocked_email = f"blocked-{suffix}@example.com"
        for _ in range(5):
            failed_login = page.request.post(
                url(base_url, "/signin"),
                form={
                    "_csrf": csrf_token,
                    "email": blocked_email,
                    "password": "wrong-password",
                },
            )
            assert failed_login.status == 401
        rate_limited = page.request.post(
            url(base_url, "/signin"),
            form={
                "_csrf": csrf_token,
                "email": blocked_email,
                "password": "wrong-password",
            },
        )
        assert rate_limited.status == 429
        assert rate_limited.headers["retry-after"] == "900"

        page.get_by_label("邮箱").fill(email)
        page.get_by_label("密码").fill(password)
        page.get_by_role("button", name="登录").click()
        wait_path(page, "/admin")

        cookies = {cookie["name"]: cookie for cookie in context.cookies()}
        session = cookies.get("goblog_session")
        assert session and session["httpOnly"] and session["sameSite"] == "Lax"
        assert "email" not in cookies

        rejected = page.request.post(
            url(base_url, "/admin/category/save"), form={"name": "missing-csrf"}
        )
        assert rejected.status == 403

        page.goto(url(base_url, "/admin/category/add"))
        page.get_by_label("分类名称").fill(category)
        page.get_by_role("button", name="保存").click()
        wait_path(page, "/admin/category")
        expect(page.get_by_text(category, exact=True)).to_be_visible()

        page.goto(url(base_url, "/admin/page/add"))
        page.get_by_label("标题").fill(page_title)
        ident_input = page.get_by_label("页面标识")
        ident_input.fill(page_ident)
        assert ident_input.evaluate("element => element.checkValidity()")
        set_markdown(page, "# 初始页面\n\n<script>window.__goblog_xss = true</script>")
        page.get_by_role("button", name="保存").click()
        wait_path(page, "/admin/page")

        page.goto(url(base_url, f"/page/{page_ident}"))
        expect(page.get_by_role("heading", name="初始页面")).to_be_visible()
        assert not page.evaluate("Boolean(window.__goblog_xss)")
        expect(page.locator("article script")).to_have_count(0)

        page.goto(url(base_url, "/admin/page"))
        page_row = page.locator("tr", has_text=page_title)
        page_row.get_by_role("link", name="编辑").click()
        page.get_by_label("标题").fill(page_title + "-已编辑")
        set_markdown(page, "# 已编辑页面")
        page.get_by_role("button", name="保存").click()
        wait_path(page, "/admin/page")
        page.goto(url(base_url, f"/page/{page_ident}"))
        expect(page.get_by_role("heading", name="已编辑页面")).to_be_visible()

        page.goto(url(base_url, "/admin/post/add"))
        page.get_by_label("标题").fill(post_title)
        page.get_by_label("分类").select_option(label=category)
        page.get_by_label("标签").fill(tag)
        page.get_by_label("摘要").fill("端到端搜索摘要")
        set_markdown(
            page,
            "# 安全正文\n\n<script>window.__goblog_post_xss = true</script>\n\n[危险链接](javascript:alert(1))",
        )
        page.get_by_role("button", name="保存").click()
        wait_path(page, "/admin")
        expect(page.get_by_text(post_title, exact=True)).to_be_visible()

        page.goto(url(base_url, "/?per_page=10"))
        expect(page.get_by_role("search")).to_have_count(0)
        page.get_by_label("搜索关键词").fill("端到端搜索摘要")
        page.get_by_label("文章分类").select_option(label=category)
        page.get_by_role("button", name="筛选").click()
        filtered_url = urlparse(page.url)
        assert filtered_url.query.count("per_page=10") == 1
        assert "keyword" in filtered_url.query
        assert "category_id" in filtered_url.query
        expect(page.get_by_text(post_title, exact=True)).to_be_visible()
        expect(page.get_by_text("本页显示 1 篇", exact=True)).to_be_visible()
        expect(page.get_by_text(f"分类：{category}", exact=True)).to_be_visible()

        page.goto(url(base_url, "/"))
        page.get_by_label("搜索关键词").fill("不存在的端到端内容")
        page.get_by_role("button", name="筛选").click()
        filtered_empty_url = urlparse(page.url)
        assert "category_id=" not in filtered_empty_url.query
        assert "per_page=20" not in filtered_empty_url.query
        expect(page.get_by_text("没有找到符合条件的文章")).to_be_visible()
        expect(page.get_by_role("link", name="重置全部")).to_have_attribute(
            "href", "/"
        )

        page.goto(url(base_url, f"/?keyword=端到端搜索摘要&page=9&per_page=10"))
        previous = page.get_by_role("link", name="上一页")
        expect(previous).to_be_visible()
        previous_href = previous.get_attribute("href")
        assert previous_href
        assert "keyword" in previous_href
        assert "page=8" in previous_href
        assert "per_page=10" in previous_href

        page.goto(url(base_url, "/"))
        page.get_by_role("link", name=post_title).click()
        page.wait_for_url("**/post/*")
        post_url = page.url
        expect(page.get_by_role("search")).to_be_visible()
        expect(page.get_by_role("heading", name="安全正文")).to_be_visible()
        assert not page.evaluate("Boolean(window.__goblog_post_xss)")
        expect(page.locator('a[href^="javascript:"]')).to_have_count(0)

        page.goto(url(base_url, "/?keyword=端到端搜索摘要"))
        expect(page.get_by_text(post_title, exact=True)).to_be_visible()
        page.goto(url(base_url, "/"))
        page.get_by_role("link", name=category, exact=True).first.click()
        expect(page.get_by_text(post_title, exact=True)).to_be_visible()
        page.goto(url(base_url, "/tag"))
        tag_link = page.locator("a", has_text=tag)
        expect(tag_link).to_be_visible()
        tag_link.click()
        expect(page.get_by_text(post_title, exact=True)).to_be_visible()

        page.goto(url(base_url, "/admin/"))
        post_row = page.locator("tr", has_text=post_title)
        post_row.get_by_role("link", name="编辑").click()
        page.get_by_label("摘要").fill("已编辑摘要")
        set_markdown(page, "# 已编辑正文")
        page.get_by_role("button", name="保存").click()
        wait_path(page, "/admin")
        page.goto(post_url)
        expect(page.get_by_role("heading", name="已编辑正文")).to_be_visible()

        page.goto(url(base_url, "/admin/"))
        submit_delete(page, page.locator("tr", has_text=post_title), "/admin/")
        deleted_post = page.request.get(post_url)
        assert deleted_post.status == 404
        assert "文章不存在" in deleted_post.text()

        page.goto(url(base_url, "/tag"))
        expect(page.get_by_text(tag, exact=False)).to_be_visible()
        assert "0" in page.locator("a", has_text=tag).inner_text()

        page.goto(url(base_url, "/admin/page"))
        submit_delete(page, page.locator("tr", has_text=page_title), "/admin/page")
        deleted_page = page.request.get(url(base_url, f"/page/{page_ident}"))
        assert deleted_page.status == 404
        assert "页面不存在" in deleted_page.text()

        page.goto(url(base_url, "/admin/category"))
        submit_delete(page, page.locator("tr", has_text=category), "/admin/category")
        expect(page.get_by_text(category, exact=True)).to_have_count(0)

        mobile = context.new_page()
        mobile.set_viewport_size({"width": 390, "height": 844})
        mobile.goto(url(base_url, "/"))
        mobile.wait_for_load_state("networkidle")
        dimensions = mobile.evaluate(
            "() => ({scroll: document.documentElement.scrollWidth, viewport: window.innerWidth})"
        )
        assert dimensions["scroll"] <= dimensions["viewport"], dimensions
        toggle = mobile.get_by_role("button", name="展开导航")
        toggle.click()
        assert toggle.get_attribute("aria-expanded") == "true"
        mobile.close()

        page.goto(url(base_url, "/admin/"))
        page.get_by_role("button", name="退出登录").click()
        wait_path(page, "/login")
        page.goto(url(base_url, "/admin/"))
        wait_path(page, "/login")

        browser.close()

    assert not browser_errors, "\n".join(browser_errors)


def main() -> None:
    parser = argparse.ArgumentParser()
    parser.add_argument("--base-url", default="http://127.0.0.1:9091")
    parser.add_argument("--email", required=True)
    parser.add_argument("--password", required=True)
    args = parser.parse_args()
    run(args.base_url, args.email, args.password)
    print("goblog E2E passed")


if __name__ == "__main__":
    main()
