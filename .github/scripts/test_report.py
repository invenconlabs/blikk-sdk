#!/usr/bin/env python3
import os
import re
import sys
import xml.etree.ElementTree as ET

MARKER = "<!-- blikk-test-report -->"
DEFAULT_PACKAGE_DIR = "blikk/"
MAX_REASON_LENGTH = 200
TEST_FILE_REFERENCE = re.compile(r"([\w.\-]+_test\.go):(\d+):")


def non_empty_lines(text):
    return [line.strip() for line in (text or "").splitlines() if line.strip()]


def failure_reason(failure):
    message_lines = non_empty_lines(failure.get("message"))
    text_lines = non_empty_lines(failure.text)
    for line in message_lines + text_lines:
        if TEST_FILE_REFERENCE.search(line):
            return line
    return (message_lines + text_lines + [""])[0]


def package_dir(suite_name):
    if suite_name and suite_name.endswith("/blikk"):
        return "blikk/"
    return DEFAULT_PACKAGE_DIR


def failed_cases(root):
    cases = []
    for suite in root.iter("testsuite"):
        for case in suite.iter("testcase"):
            failure = case.find("failure")
            if failure is None:
                failure = case.find("error")
            if failure is not None:
                cases.append((case.get("name", ""), failure_reason(failure), package_dir(suite.get("name"))))
    return cases


def count_cases(root):
    cases = list(root.iter("testcase"))
    skipped = sum(1 for case in cases if case.find("skipped") is not None)
    failed = len(failed_cases(root))
    return len(cases) - failed - skipped, failed, skipped


def total_seconds(root):
    return sum(float(suite.get("time") or 0) for suite in root.iter("testsuite"))


def clip(text):
    return text[:MAX_REASON_LENGTH]


def annotation(name, reason, directory):
    reference = TEST_FILE_REFERENCE.search(reason)
    if reference is None:
        return None
    file_name, line = reference.groups()
    escaped_reason = reason.replace("%", "%25")
    return f"::error file={directory}{file_name},line={line},title={name}::{escaped_reason}"


def render_failure_bullet(name, reason):
    safe_reason = clip(reason).replace("`", "'")
    return f"- `{name}` — `{safe_reason}`"


def render_section(title, path):
    if not os.path.exists(path):
        return [f"### {title}", "⏭️ job did not run"]
    root = ET.parse(path).getroot()
    passed, failed, skipped = count_cases(root)
    icon = "❌" if failed else "✅"
    lines = [
        f"### {title}",
        f"{icon} {passed} passed · {failed} failed · {skipped} skipped · {total_seconds(root):.1f}s",
    ]
    for name, reason, directory in failed_cases(root):
        lines.append(render_failure_bullet(name, reason))
        line = annotation(name, reason, directory)
        if line:
            print(line, file=sys.stderr)
    return lines


def main(arguments):
    lines = ["## Test results"]
    for argument in arguments:
        title, path = argument.split("=", 1)
        lines += render_section(title, path)
    lines += ["", MARKER]
    print("\n".join(lines))


if __name__ == "__main__":
    main(sys.argv[1:])
