#!/usr/bin/env python3
"""Regenerate Iris models from the pinned Iris checkout's OpenAPI outputs.

Run from this repository's Nix shell:
    python3 scripts/generate-iris-models.py /path/to/iris

Optionally set OAPI_CODEGEN to an existing v2.8.2 binary. Otherwise the pinned
generator runs through Go without changing this module's dependencies.
"""

import hashlib
import json
import os
from pathlib import Path
import subprocess
import sys
import tempfile


VERSION = "v2.8.2"
# Iris commit 23ff21223ab3b7a06c37e5fc1c4e1a3564f6d442. Hash the generated
# inputs too: those outputs are ignored by Iris Git and can differ locally.
SOURCES = {
    "cp-memory/MemoryControlPlaneService.openapi.json":
        "cbcc61b0d078ecd5275eb1cecfca72a4b3f7a33f81f17519229f49a18d606eb2",
    "cp-onprem-memory/MemoryService.openapi.json":
        "411c4d53d09aa757a5094029fe2e531e6f4d9765c8f4cb6a36200c3ae659d0f9",
    "ids-control/IdentityControlService.openapi.json":
        "c1637b7c8181f03acc1c30a9a57a625d99ce5847616c8dd03bcd87a1999f388c",
}


def main():
    if len(sys.argv) != 2:
        raise SystemExit("usage: generate-iris-models.py IRIS_CHECKOUT")
    root = Path(__file__).resolve().parents[1]
    source = Path(sys.argv[1]).resolve() / "generated/smithy/source"
    specs = []
    for name, digest in SOURCES.items():
        data = (source / name).read_bytes()
        if hashlib.sha256(data).hexdigest() != digest:
            raise SystemExit(f"Iris contract changed: {name}; review it before updating SOURCES")
        specs.append(json.loads(data))

    memory, onprem, identity = specs
    # The on-prem service contributes ingestion hooks; common AM owns stores,
    # extractors, memory types and reference data.
    memory["paths"].update({
        path: operations for path, operations in onprem["paths"].items()
        if "ingestion-hooks" in path
    })
    memory["components"]["schemas"].update(onprem["components"]["schemas"])

    binary = os.environ.get("OAPI_CODEGEN")
    command = [binary] if binary else [
        "go", "run", f"github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@{VERSION}"
    ]
    version = subprocess.check_output(command + ["--version"], cwd=root, text=True)
    if VERSION not in version.splitlines():
        raise SystemExit(f"expected oapi-codegen {VERSION}, got {version.strip()}")
    with tempfile.TemporaryDirectory(prefix="iris-models-") as directory:
        for package, spec in [("irisagentmemory", memory), ("irisidentity", identity)]:
            schema = Path(directory) / f"{package}.json"
            schema.write_text(json.dumps(spec))
            subprocess.run(command + [
                "-generate", "types", "-package", package,
                "-o", str(root / "service" / package / "model.gen.go"), str(schema)
            ], cwd=root, check=True)


if __name__ == "__main__":
    main()
