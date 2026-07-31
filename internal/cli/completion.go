package cli

// GetCompletionScript returns a basic bash/zsh completion script for stackgenome.
func GetCompletionScript() string {
	return `#!/usr/bin/env bash
# StackGenome completion script for bash/zsh

_stackgenome_completions()
{
    local cur prev words cword
    _init_completion || return

    local commands="analyze completion version help"
    local analyze_flags="--env --recommend --remote --json -h --help"

    if [[ ${cword} -eq 1 ]]; then
        COMPREPLY=( $(compgen -W "${commands}" -- ${cur}) )
        return 0
    fi

    case ${prev} in
        analyze)
            COMPREPLY=( $(compgen -W "${analyze_flags}" -- ${cur}) )
            return 0
            ;;
    esac

    # Handle completion for flags within the analyze subcommand
    if [[ "${words[1]}" == "analyze" && "${cur}" == -* ]]; then
        COMPREPLY=( $(compgen -W "${analyze_flags}" -- ${cur}) )
        return 0
    fi

    # Fallback to directory completion for the target dir argument
    if [[ "${words[1]}" == "analyze" ]]; then
        COMPREPLY=( $(compgen -d -- ${cur}) )
        return 0
    fi
}

complete -F _stackgenome_completions stackgenome
`
}
