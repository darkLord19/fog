#compdef wtx

_wtx() {
    local -a commands
    commands=(
        'list:List all worktrees'
        'add:Create a new worktree'
        'open:Open a worktree in editor'
        'rm:Remove a worktree'
        'prune:Remove stale worktrees'
        'status:Show worktree status'
        'config:Manage configuration'
        'version:Show version'
        'help:Show help'
    )

    local -a worktrees
    worktrees=(${(f)"$(wtx list 2>/dev/null | awk '{print $1":"$2}')"})

    _arguments -C \
        '1: :->command' \
        '*:: :->args' \
        '--help[Show help]' \
        '--version[Show version]' \
        '--editor[Specify editor]:editor:(vscode cursor neovim vim)' \
        '--json[Output as JSON]'

    case $state in
        command)
            _describe 'command' commands
            ;;
        args)
            case $words[1] in
                open|rm|status)
                    _describe 'worktree' worktrees
                    ;;
                add)
                    _message 'worktree name'
                    ;;
            esac
            ;;
    esac
}

_wtx "$@"
