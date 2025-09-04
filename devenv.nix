{ pkgs, lib, config, inputs, ... }:

{
  # https://devenv.sh/basics/
  env.GREET = "devenv";

  # https://devenv.sh/packages/
  packages = [ 
    pkgs.git
    pkgs.gopls        # Go language server
    pkgs.go-tools     # Additional Go tools (goimports, etc.)
    pkgs.delve        # Go debugger
    pkgs.llvm_17      # LLVM for future backend work
    pkgs.go
  ];

  # https://devenv.sh/languages/
  languages.go.enable = true;

  enterShell = ''
    echo
    echo "🚀 Language Development Environment"
    echo "=================================="
    echo "Go version: $(go version)"
    echo "LLVM version: $(llvm-config --version)"
    echo ""
    echo "Available scripts:"
    echo "  devenv test  - Run all tests"
    echo "  devenv build - Build compiler binary"
    echo "  devenv run   - Run compiler with args"
    echo "  devenv fmt   - Format Go code"
    echo "  devenv lint  - Run Go linter"
    echo ""
  '';

  # https://devenv.sh/pre-commit-hooks/
  git-hooks.hooks = {
    gofmt.enable = true;
    govet.enable = true;
  };

  # See full reference at https://devenv.sh/reference/options/
}