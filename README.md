# gitx 🚀

> **Git Identity Switcher** - Manage multiple GitHub identities safely with per-repo binding. Never push to the wrong account again.

[![Go Version](https://img.shields.io/badge/go-1.21+-00ADD8?style=flat-square&logo=go)](https://golang.org/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=flat-square)](https://opensource.org/licenses/MIT)
[![PRs Welcome](https://img.shields.io/badge/PRs-welcome-brightgreen.svg?style=flat-square)](http://makeapullrequest.com)

## ✨ Features

- 🔐 **Per-repo identity binding** - Each repository can be bound to a specific identity
- 🔑 **SSH key management** - Automatic SSH key generation and secure storage
- 🎫 **PAT support** - Personal Access Token support for HTTPS authentication
- 🔒 **Secure storage** - Secrets stored in OS keychain (macOS Keychain, Linux Secret Service)
- 🛡️ **SSH config safety** - Managed SSH config blocks with automatic backups
- 👀 **Dry-run mode** - Preview changes before applying them
- 🎨 **TUI interface** - Interactive text-based UI for identity selection
- 🚦 **Pre-push hooks** - Optional safety hook to prevent unbound pushes
- ⚡ **Fast & lightweight** - Single binary, no dependencies

## 🎯 Why gitx?

Working with multiple GitHub accounts (work, personal, client) is a pain. You've probably experienced:

- ❌ Accidentally pushing to the wrong account
- ❌ Wrong commit author/email
- ❌ SSH config conflicts
- ❌ Manual identity switching

**gitx solves all of this** by binding identities to repositories, so you never have to think about it again.

## 📦 Installation

### From Source

```bash
git clone https://github.com/csawai/gitx.git
cd gitx
go build -o gitx .
sudo mv gitx /usr/local/bin/
```

### Using Go Install

```bash
go install github.com/csawai/git-identity-switcher@latest

# The binary will be named 'git-identity-switcher'
# Create an alias for convenience:
alias gitx='git-identity-switcher'
```

**Upgrading:** Run the same `go install` command again to get the latest version.

### Using Homebrew (macOS)

```bash
brew install csawai/tap/gitx
```

## 🚀 Quick Start

### 1. Add an identity

```bash
gitx add identity
```

Follow the prompts to enter:
- Identity alias (e.g., "work", "personal")
- Name
- Email
- GitHub username
- Auth method (SSH or PAT)

**For SSH:** After key generation, gitx will display your public key. Add it to GitHub at https://github.com/settings/ssh/new. You can also use `gitx show-key <alias>` or `gitx copy-key <alias>` later.

### 2. List identities

```bash
gitx list identities
```

### 3. Bind a repository

```bash
cd /path/to/repo
gitx bind work
```

This automatically:
- Sets `user.name` and `user.email` in the repository's local git config
- Updates the remote URL to use the identity's SSH host alias (for SSH) or HTTPS (for PAT)

### 4. Check status

```bash
gitx status
```

Shows current identity configuration for the repository.

### 5. Unbind a repository

```bash
gitx unbind
```

Reverts gitx changes to the repository.

## 📖 Commands

| Command | Description |
|---------|-------------|
| `gitx version` | Show version |
| `gitx status` | Show current repository identity status |
| `gitx add identity` | Add a new identity |
| `gitx list identities` | List all configured identities |
| `gitx show-key <alias>` | Show SSH public key for an identity |
| `gitx copy-key <alias>` | Copy SSH public key to clipboard |
| `gitx bind <alias>` | Bind repository to an identity |
| `gitx unbind` | Unbind repository from identity |
| `gitx remove identity <alias>` | Remove an identity |
| `gitx tui` | Launch interactive TUI |
| `gitx install-hook` | Install pre-push safety hook |
| `gitx uninstall-hook` | Remove pre-push hook |

## 🔧 How It Works

### SSH Authentication

1. Each identity gets its own SSH key: `~/.ssh/gitx_<alias>`
2. SSH config entries are added using host aliases (e.g., `github.com-work`)
3. Repository remote URLs are updated to use the host alias: `git@github.com-work:org/repo.git`
4. This avoids conflicts with your default GitHub SSH config

### PAT Authentication

1. PAT tokens are stored securely in the OS keychain
2. Remote URLs are converted to HTTPS format
3. Git credential helper is configured to use the keychain

### SSH Config Management

All gitx-managed entries are in a marked block:

```
# BEGIN gitx managed
Host github.com-work
  HostName github.com
  User git
  IdentityFile ~/.ssh/gitx_work
  IdentitiesOnly yes
# END gitx managed
```

- Automatic backups are created before any changes
- Atomic writes ensure config is never corrupted
- Only the managed block is modified; your existing config is untouched

## 🛡️ Safety Features

- **Dry-run mode**: Use `--dry-run` flag to preview changes
- **Automatic backups**: SSH config is backed up before modifications
- **Atomic writes**: Changes are written to temp files, validated, then swapped
- **Pre-push hooks**: Optional hook prevents pushes from unbound repositories

## 📁 Configuration

- **Identities**: `~/.config/gitx/identities.json`
- **SSH keys**: `~/.ssh/gitx_<alias>`
- **Secrets (PATs)**: OS keychain under service name "gitx"

## 💡 Examples

### Workflow Example

```bash
# Add your work identity
gitx add identity
# Alias: work
# Name: John Doe
# Email: john@company.com
# GitHub: johndoe
# Auth: ssh

# Add your personal identity
gitx add identity
# Alias: personal
# Name: John Doe
# Email: john@gmail.com
# GitHub: johndoe-personal
# Auth: pat

# Bind work repo
cd ~/projects/work-project
gitx bind work

# Bind personal repo
cd ~/projects/personal-project
gitx bind personal

# Check status
gitx status
```

### Using TUI

```bash
cd ~/projects/my-repo
gitx tui
# Interactive menu to select and bind identity
```

## 🧪 Testing

```bash
go test ./...
```

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

### Development Setup

```bash
git clone https://github.com/csawai/gitx.git
cd gitx
go mod download
go build -o gitx .
```

## 📋 Requirements

- Go 1.21+
- Git
- SSH (for SSH authentication)
- macOS Keychain or Linux Secret Service (for secure storage)

## 🐛 Troubleshooting

### SSH key not working

Make sure you've added the public key to your GitHub account:

```bash
cat ~/.ssh/gitx_<alias>.pub
# Copy and add to GitHub Settings > SSH and GPG keys
```

### Keychain access issues (macOS)

You may need to grant Terminal/iTerm access to the keychain. Go to System Preferences > Security & Privacy > Privacy > Keychain Access.

### Permission denied errors

Make sure your SSH keys have the correct permissions:

```bash
chmod 600 ~/.ssh/gitx_*
```

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- Built with [Cobra](https://github.com/spf13/cobra) for CLI
- TUI powered by [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- Secure storage via [99designs/keyring](https://github.com/99designs/keyring)

## ⭐ Star History

If you find this project useful, please consider giving it a star ⭐

---

**Made with ❤️ by [csawai](https://github.com/csawai)**
