#!/bin/bash

_wtx_completions() {
    local cur prev opts
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"
    
    # Main commands
    local commands="list add open rm prune status config version help"
    
    # Flags
    local flags="--help --version --editor --json"
    
    # If completing first argument
    if [ $COMP_CWORD -eq 1 ]; then
        COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
        return 0
    fi
    
    # Command-specific completions
    case "${COMP_WORDS[1]}" in
        open|rm|status)
            # Complete with worktree names
            local worktrees=$(wtx list 2>/dev/null | awk '{print $1}')
            COMPREPLY=( $(compgen -W "${worktrees}" -- ${cur}) )
            return 0
            ;;
        add)
            # No special completion for add
            return 0
            ;;
        list)
            COMPREPLY=( $(compgen -W "--json" -- ${cur}) )
            return 0
            ;;
        *)
            COMPREPLY=( $(compgen -W "${flags}" -- ${cur}) )
            return 0
            ;;
    esac
}

complete -F _wtx_completions wtx
