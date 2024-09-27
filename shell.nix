{ pkgs ? import <nixpkgs> {} }:
pkgs.mkShell {
  buildInputs = with pkgs; [
    gnumake
    go
    golangci-lint
  ];
  shellHook = ''
    echo '
┏━╸┏━┓┏┳┓┏━┓┏━┓╻┏━╸┏┓╻┏━╸┏━┓
┃  ┣━┫┃┃┃┣━┛┣━┫┃┃╺┓┃┗┫┣╸ ┣┳┛
┗━╸╹ ╹╹ ╹╹  ╹ ╹╹┗━┛╹ ╹┗━╸╹┗╸
'
  '';
  hardeningDisable = [ "all" ];
}
