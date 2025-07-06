# Contributing to s3www

Thank you for your interest in contributing to s3www! This project is a lightweight tool for serving static files from S3-compatible storage, and we welcome contributions from the community to improve its features, performance, and usability. Whether you're fixing bugs, adding features, or improving documentation, your help is appreciated.

## How to Contribute

### Reporting Issues

If you encounter bugs, have feature requests, or find areas for improvement:

1. Check the [GitHub Issues](https://github.com/harshavardhana/s3www/issues) page to see if the issue already exists.
2. If not, open a new issue with a clear title and description, including:
   - Steps to reproduce the issue (if applicable).
   - Expected and actual behavior.
   - Relevant details like your operating system, s3www version, and S3 provider.

### Submitting Changes

To contribute code or documentation:

1. **Fork the Repository**:
   - Fork the [s3www repository](https://github.com/harshavardhana/s3www) to your GitHub account.
   - Clone your fork:
     ```bash
     git clone https://github.com/<your-username>/s3www.git
     cd s3www
     ```

2. **Create a Branch**:
   - Create a new branch for your changes:
     ```bash
     git checkout -b feature/<feature-name>  # For new features
     # or
     git checkout -b fix/<bug-name>         # For bug fixes
     ```

3. **Make Changes**:
   - Follow the project’s coding style (e.g., adhere to Go conventions for code).
   - Keep changes focused and commit messages clear (e.g., `Add support for custom CORS headers`).
   - Update tests if applicable (located in the `*_test.go` files).
   - If modifying documentation, ensure clarity and consistency with the `README.md`.

4. **Test Your Changes**:
   - Build and test locally:
     ```bash
     go build
     go test ./...
     ```
   - Verify your changes work with an S3-compatible bucket.

5. **Commit and Push**:
   - Commit your changes with a descriptive message:
     ```bash
     git commit -m "Add feature X to improve Y"
     ```
   - Push your branch to your fork:
     ```bash
     git push origin feature/<feature-name>
     ```

6. **Open a Pull Request**:
   - Go to the [s3www repository](https://github.com/harshavardhana/s3www) and open a pull request from your branch.
   - Provide a clear description of your changes, referencing any related issues (e.g., `Fixes #123`).
   - Ensure your PR passes any automated checks (e.g., CI tests).

## Code of Conduct

- Be respectful and inclusive in all interactions.
- Follow the [Contributor Covenant Code of Conduct](https://www.contributor-covenant.org/version/2/0/code_of_conduct.html).
- Avoid personal attacks, harassment, or discriminatory language.

## Development Guidelines

- **Code Style**: Follow Go’s standard formatting (`go fmt`) and conventions. Use meaningful variable names and include comments for complex logic.
- **Testing**: Add or update tests for new features or bug fixes. Ensure tests pass before submitting a PR.
- **Documentation**: Update `README.md` or other docs if your changes affect usage or configuration.
- **Dependencies**: Avoid adding new dependencies unless necessary, and justify their inclusion in your PR.

## Areas for Contribution

- **Bug Fixes**: Address issues listed in the [GitHub Issues](https://github.com/harshavardhana/s3www/issues) page.
- **Features**: Add support for new S3 providers, enhance SPA routing, or improve performance.
- **Documentation**: Clarify setup instructions, add examples, or create tutorials.
- **Testing**: Improve test coverage or add integration tests for different S3 providers.

## Getting Help

If you have questions or need assistance:
- Open an issue with the `question` label.
- Reach out to the maintainers via the [GitHub Discussions](https://github.com/harshavardhana/s3www/discussions) page (if available).
- Check the [README.md](README.md) for setup and usage details.

## License

By contributing, you agree that your contributions will be licensed under the [Apache License, Version 2.0](LICENSE).

Thank you for helping make s3www better!