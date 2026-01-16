# Security Policy

## Supported Versions

| Version | Supported          |
|---------|--------------------|
| 1.x.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability in envx, please report it to us privately before disclosing it publicly.

### How to Report

Please send an email to: **security@envx.dev**

Include the following information:
- Type of vulnerability
- Steps to reproduce
- Potential impact
- Any proposed mitigations (if known)

### Response Time

- **Critical**: 24-48 hours
- **High**: 72 hours
- **Medium**: 1 week
- **Low**: 2 weeks

### Security Updates

When a security vulnerability is fixed:
1. We will create a security advisory
2. Fix the vulnerability in a new release
3. Publish the fix and security advisory together

### Public Disclosure

We will disclose vulnerabilities after:
- A fix is available
- Users have reasonable time to update
- Coordination with maintainers if the vulnerability affects other projects

### Security Best Practices

For users of envx:
1. Keep your dependencies updated
2. Use tools like `govulncheck` to scan for known vulnerabilities
3. Review security advisories for your version
4. Follow principle of least privilege when setting environment variables

### Security Features

envx includes several security features:
- Type-safe environment variable parsing
- Validation of required fields
- Support for sensitive data masking in logs
- No execution of external commands