"""Protect tap publishing against corrupt assets and unintended release versions."""

import hashlib
import tempfile
import unittest
from pathlib import Path

from homebrew_formula import homebrew_formula, stable_version, validate_release


class HomebrewFormulaTest(unittest.TestCase):
    def setUp(self):
        temporary = tempfile.TemporaryDirectory()
        self.addCleanup(temporary.cleanup)
        self.directory = Path(temporary.name)
        checksums = []
        for system in ("darwin", "linux"):
            for arch in ("arm64", "amd64"):
                filename = f"rolecule_1.2.3_{system}_{arch}.tar.gz"
                content = filename.encode()
                (self.directory / filename).write_bytes(content)
                checksums.append(f"{hashlib.sha256(content).hexdigest()}  {filename}\n")
        self.manifest = self.directory / "rolecule_1.2.3_SHA256SUMS"
        self.manifest.write_text("".join(checksums))

    def test_all_platforms_use_verified_release_archives(self):
        formula = homebrew_formula("v1.2.3", self.directory)
        for system in ("darwin", "linux"):
            for arch in ("arm64", "amd64"):
                filename = f"rolecule_1.2.3_{system}_{arch}.tar.gz"
                digest = hashlib.sha256((self.directory / filename).read_bytes()).hexdigest()
                self.assertIn(f"/v1.2.3/{filename}", formula)
                self.assertIn(f'sha256 "{digest}"', formula)
        self.assertIn('assert_equal "rolecule 1.2.3"', formula)

    def test_invalid_and_prerelease_tags_are_rejected(self):
        for tag in ("v1.2.3-rc.1", "v1.02.3", "1.2.3", 'v1.2.3"; system("bad")'):
            with self.subTest(tag=tag), self.assertRaises(ValueError):
                homebrew_formula(tag, self.directory)
        self.assertEqual(stable_version("v1.2.3+build.5"), "1.2.3+build.5")

    def test_corrupt_or_missing_archives_are_rejected(self):
        for system in ("darwin", "linux"):
            archive = self.directory / f"rolecule_1.2.3_{system}_arm64.tar.gz"
            content = archive.read_bytes()
            archive.write_bytes(b"corrupt download")
            with self.assertRaisesRegex(ValueError, "mismatched checksum"):
                homebrew_formula("v1.2.3", self.directory)
            archive.unlink()
            with self.assertRaises(FileNotFoundError):
                homebrew_formula("v1.2.3", self.directory)
            archive.write_bytes(content)

    def test_missing_or_duplicate_checksums_are_rejected(self):
        content = self.manifest.read_text()
        self.manifest.write_text("".join(content.splitlines(keepends=True)[1:]))
        with self.assertRaisesRegex(ValueError, "Missing or mismatched checksum"):
            homebrew_formula("v1.2.3", self.directory)
        self.manifest.write_text(content + content)
        with self.assertRaisesRegex(ValueError, "Duplicate checksum"):
            homebrew_formula("v1.2.3", self.directory)

    def test_only_latest_published_stable_release_updates_tap(self):
        release = {"tagName": "v1.2.3", "isDraft": False, "isPrerelease": False}
        self.assertTrue(validate_release("v1.2.3", release, "v1.2.3"))
        with self.assertRaisesRegex(ValueError, "other than"):
            validate_release("v1.2.3", release, "v1.2.4")
        with self.assertRaises(ValueError):
            validate_release("v1.2.4", release, "v1.2.4")
        with self.assertRaises(ValueError):
            validate_release("v1.2.3", dict(release, isDraft=True), "v1.2.3")
        self.assertFalse(validate_release("v1.2.3-rc.1", dict(
            release, tagName="v1.2.3-rc.1", isPrerelease=True,
        ), "v1.2.2"))
        with self.assertRaises(ValueError):
            validate_release("v1.2.3-rc.1", dict(release, tagName="v1.2.3-rc.1"), "v1.2.3-rc.1")


if __name__ == "__main__":
    unittest.main()
